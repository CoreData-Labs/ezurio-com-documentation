// Command ezurio-scraper visits every product support page linked from
// https://www.ezurio.com/support (/support/series/{product-slug}), follows the
// documentation pages (/documentation/{document-slug}) listed on them, finds all
// PDF files and downloads them into the PDFs/ folder.
// It uses only the Go standard library. Logs go to standard error and the
// discovered product support URLs go to standard output.
package main // this file builds an executable program

import ( // standard library packages used by this program
	"bufio"         // wraps the response body so its first bytes can be inspected
	"bytes"         // searches byte slices for the PDF signature
	"crypto/sha256" // creates short hashes that keep file names unique
	"encoding/hex"  // converts hashes to readable text
	"errors"        // creates and compares error values
	"fmt"           // formats text and prints the product URLs
	"html"          // decodes HTML entities such as &amp; inside href values
	"io"            // copies data and provides the discard writer
	"log"           // provides the console loggers
	"net/http"      // performs the HTTP requests
	"net/url"       // parses and resolves URLs
	"os"            // creates folders and files, and provides standard error
	"path"          // extracts the last segment of a URL path
	"path/filepath" // builds file system paths
	"regexp"        // finds links in the HTML text
	"sort"          // sorts URL lists for stable output
	"strings"       // provides string helper functions
	"sync"          // provides mutexes and wait groups for concurrency
	"sync/atomic"   // provides thread-safe counters
	"time"          // provides timeouts and retry delays
) // end of the import list

const ( // fixed configuration values
	supportPageAddress                               = "https://www.ezurio.com/support"                                                                                 // page where scraping starts
	supportSeriesPathPrefix                          = "/support/series/"                                                                                               // path prefix every product support link must have
	documentationPathPrefix                          = "/documentation/"                                                                                                // path prefix of the documentation pages
	canonicalHostname                                = "www.ezurio.com"                                                                                                 // hostname written in every printed URL
	acceptedHostnameWithoutWww                       = "ezurio.com"                                                                                                     // alternative hostname accepted in links
	maximumFetchAttempts                             = 3                                                                                                                // how many times a failing download is tried
	baseDelayBetweenAttempts                         = 2 * time.Second                                                                                                  // retry delay, multiplied by the attempt number
	httpRequestTimeout                               = 30 * time.Second                                                                                                 // maximum time allowed for fetching one HTML page
	pdfDownloadTimeout                               = 5 * time.Minute                                                                                                  // maximum time allowed for downloading one PDF
	concurrentWorkerCount                            = 4                                                                                                                // number of requests running at the same time
	pauseAfterEachItem                               = 250 * time.Millisecond                                                                                           // polite pause a worker takes after each request
	pdfOutputDirectory                               = "PDFs"                                                                                                           // folder that receives the downloaded PDFs
	pdfOutputDirectoryPermissions        os.FileMode = 0o755                                                                                                            // permissions of the PDFs folder
	pdfSignature                                     = "%PDF-"                                                                                                          // text every real PDF file starts with
	pdfSignatureSearchLength                         = 1024                                                                                                             // how many leading bytes are searched for the signature
	sessionCookieEnvironmentVariableName             = "EZURIO_SESSION_COOKIE"                                                                                          // optional environment variable holding the Cookie header of a logged-in browser session
	fileDownloadQueryParameter                       = "file_download"                                                                                                  // query parameter used by Ezurio for gated file downloads (/account?file_download=NAME)
	debugLoggingEnabled                              = true                                                                                                             // set to false to hide debug log lines
	userAgentHeaderValue                             = "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36" // value sent in the User-Agent header
) // end of the constants

// anchorElementPattern captures the attribute text and the inner HTML of every <a ...>...</a> element.
var anchorElementPattern = regexp.MustCompile(`(?is)<a\b([^>]*)>(.*?)</a>`)

// hrefAttributePattern captures the href value from an attribute string. The value may be
// double-quoted, single-quoted or unquoted. Requiring whitespace before "href" prevents
// attributes such as "data-href" from being mistaken for "href".
var hrefAttributePattern = regexp.MustCompile(`(?is)(?:^|\s)href\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+))`)

// quotedAttributeValuePattern captures quoted values of href, src, data and data-* attributes,
// which is where PDF links normally appear (links, embedded viewers, download buttons).
var quotedAttributeValuePattern = regexp.MustCompile(`(?is)\b(?:href|src|data|data-[\w-]+)\s*=\s*(?:"([^"]*)"|'([^']*)')`)

// absolutePDFURLPattern captures absolute URLs ending in .pdf anywhere in the text, including inside scripts.
var absolutePDFURLPattern = regexp.MustCompile(`(?i)https?://[^\s"'<>\\]+?\.pdf(?:\?[^\s"'<>\\]*)?`)

// htmlTagPattern matches any HTML tag so the visible text of an anchor can be isolated.
var htmlTagPattern = regexp.MustCompile(`(?s)<[^>]*>`)

// htmlCommentPattern matches HTML comments so commented-out links are ignored.
var htmlCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)

// sessionCookieHeaderValue holds the optional logged-in session cookie, read once from the environment.
var sessionCookieHeaderValue = strings.TrimSpace(os.Getenv(sessionCookieEnvironmentVariableName)) // empty when not configured

// illegalFileNameCharactersPattern matches characters that operating systems forbid in file names.
var illegalFileNameCharactersPattern = regexp.MustCompile(`[\\/:*?"<>|\x00-\x1f]`)

var ( // console loggers
	informationLogger = log.New(os.Stderr, "INFO  ", log.LstdFlags) // logs normal progress messages
	errorLogger       = log.New(os.Stderr, "ERROR ", log.LstdFlags) // logs failures and warnings
	debugLogger       = createDebugLogger()                         // logs detailed diagnostics when enabled
) // end of the loggers

// anchorElement describes one hyperlink found in an HTML page.
type anchorElement struct {
	hyperlinkReference string // decoded value of the href attribute
	visibleText        string // text shown to the reader inside the link
}

// concurrentStringSet is a set of strings that many goroutines can fill safely.
type concurrentStringSet struct {
	mutex  sync.Mutex          // protects the map below
	values map[string]struct{} // the stored strings
}

// createDebugLogger returns a logger that prints only when debug logging is enabled.
func createDebugLogger() *log.Logger {
	debugOutputWriter := io.Discard // by default, throw debug output away
	if debugLoggingEnabled {        // when debug logging is switched on
		debugOutputWriter = os.Stderr // send debug output to the console
	}
	return log.New(debugOutputWriter, "DEBUG ", log.LstdFlags|log.Lmicroseconds) // build the logger
}

// main runs the whole workflow: support page, series pages, documentation pages, PDF downloads.
func main() {
	informationLogger.Printf("starting Ezurio support PDF scraper") // announce the start

	pageHTTPClient := &http.Client{Timeout: httpRequestTimeout}     // client for HTML pages
	downloadHTTPClient := &http.Client{Timeout: pdfDownloadTimeout} // client for PDF files, which can be large

	if sessionCookieHeaderValue == "" { // no cookie configured
		informationLogger.Printf("warning: %s is not set, so documents that require signing in cannot be downloaded", sessionCookieEnvironmentVariableName) // explain the limitation
	} else { // cookie configured
		informationLogger.Printf("using the session cookie from %s for requests to Ezurio", sessionCookieEnvironmentVariableName) // never print the cookie itself
	}

	supportPageContent, finalSupportPageURL, fetchError := fetchSupportPage(pageHTTPClient, supportPageAddress) // download the start page
	if fetchError != nil {                                                                                      // if the download failed
		errorLogger.Fatalf("could not fetch the support page: %v", fetchError) // log the error and exit with code 1
	}
	informationLogger.Printf("fetched support page (%d bytes) from %s", len(supportPageContent), finalSupportPageURL) // log the download size

	productSupportURLs := extractSupportProductURLs(supportPageContent, finalSupportPageURL) // find all product support URLs
	if len(productSupportURLs) == 0 {                                                        // if nothing matched
		errorLogger.Fatalf("no support product URLs were found; the links may be rendered by JavaScript or the page structure may have changed") // exit with an explanation
	}
	informationLogger.Printf("found %d unique support product URLs", len(productSupportURLs)) // log the count
	for _, productSupportURL := range productSupportURLs {                                    // go through each URL
		fmt.Println(productSupportURL) // print it to standard output
	}

	documentationPageURLs, directPDFURLs := collectLinksFromSeriesPages(pageHTTPClient, productSupportURLs)                                      // visit every series page
	informationLogger.Printf("series pages gave %d documentation pages and %d direct PDF links", len(documentationPageURLs), len(directPDFURLs)) // log the findings

	pdfURLsFromDocumentationPages := collectPDFLinksFromDocumentationPages(pageHTTPClient, documentationPageURLs) // visit every documentation page
	informationLogger.Printf("documentation pages gave %d PDF links", len(pdfURLsFromDocumentationPages))         // log the findings

	allPDFURLs := mergeUniqueSortedStrings(directPDFURLs, pdfURLsFromDocumentationPages) // combine both sources without duplicates
	if len(allPDFURLs) == 0 {                                                            // if no PDF link was discovered at all
		errorLogger.Fatalf("no PDF links were found; the download links may require signing in or may be rendered by JavaScript") // exit with an explanation
	}
	informationLogger.Printf("found %d unique PDF links in total", len(allPDFURLs)) // log the total

	downloadAllPDFFiles(downloadHTTPClient, allPDFURLs) // download every PDF into the PDFs folder
}

// fetchSupportPage downloads pageAddress and returns its text and the final URL
// after redirects. Network errors, HTTP 429 and HTTP 5xx responses are retried.
func fetchSupportPage(httpClient *http.Client, pageAddress string) (string, *url.URL, error) {
	var lastError error // remembers the most recent failure

	for attemptNumber := 1; attemptNumber <= maximumFetchAttempts; attemptNumber++ { // try up to the maximum
		debugLogger.Printf("fetch attempt %d of %d for %s", attemptNumber, maximumFetchAttempts, pageAddress) // log the attempt

		pageContent, finalPageURL, shouldRetry, fetchError := performSingleFetch(httpClient, pageAddress) // one request
		if fetchError == nil {                                                                            // on success
			return pageContent, finalPageURL, nil // hand the result back
		}

		lastError = fetchError                                                                     // keep the error for the final message
		errorLogger.Printf("attempt %d for %s failed: %v", attemptNumber, pageAddress, fetchError) // log the failure

		if !shouldRetry || attemptNumber == maximumFetchAttempts { // stop when retrying is pointless or attempts are used up
			break // leave the retry loop
		}

		delayBeforeNextAttempt := time.Duration(attemptNumber) * baseDelayBetweenAttempts // wait longer after each failure
		debugLogger.Printf("waiting %s before the next attempt", delayBeforeNextAttempt)  // log the wait
		time.Sleep(delayBeforeNextAttempt)                                                // pause
	}

	return "", nil, fmt.Errorf("giving up on %s: %w", pageAddress, lastError) // report the final failure
}

// performSingleFetch sends one GET request. The boolean result says whether a retry could help.
func performSingleFetch(httpClient *http.Client, pageAddress string) (string, *url.URL, bool, error) {
	httpRequest, requestCreationError := http.NewRequest(http.MethodGet, pageAddress, nil) // build the request
	if requestCreationError != nil {                                                       // if the address is invalid
		return "", nil, false, fmt.Errorf("building request: %w", requestCreationError) // retrying cannot help
	}
	httpRequest.Header.Set("User-Agent", userAgentHeaderValue)          // identify as a regular browser
	httpRequest.Header.Set("Accept", "text/html,application/xhtml+xml") // ask for HTML content
	httpRequest.Header.Set("Accept-Language", "en-US,en;q=0.9")         // ask for English content
	applySessionCookie(httpRequest)                                     // send the logged-in session cookie to Ezurio hosts, if configured

	httpResponse, requestError := httpClient.Do(httpRequest) // send the request and follow redirects
	if requestError != nil {                                 // on a network failure or timeout
		return "", nil, true, fmt.Errorf("sending request: %w", requestError) // worth retrying
	}
	defer httpResponse.Body.Close() // always release the connection

	debugLogger.Printf("response status=%q content-type=%q content-length=%d", httpResponse.Status, httpResponse.Header.Get("Content-Type"), httpResponse.ContentLength) // log response details

	if httpResponse.StatusCode != http.StatusOK { // any status other than 200 OK is a problem
		isTemporaryFailure := httpResponse.StatusCode == http.StatusTooManyRequests || httpResponse.StatusCode >= http.StatusInternalServerError // rate limit or server error
		return "", nil, isTemporaryFailure, fmt.Errorf("unexpected HTTP status %s", httpResponse.Status)                                         // only temporary failures are retried
	}

	contentType := strings.ToLower(httpResponse.Header.Get("Content-Type")) // read the content type
	if contentType != "" && !strings.Contains(contentType, "html") {        // warn if it does not look like HTML
		errorLogger.Printf("warning: unexpected content type %q for %s, parsing anyway", contentType, pageAddress) // log the warning
	}

	responseBodyBytes, readError := io.ReadAll(httpResponse.Body) // read the whole body
	if readError != nil {                                         // if the transfer broke midway
		return "", nil, true, fmt.Errorf("reading response body: %w", readError) // worth retrying
	}
	if len(responseBodyBytes) == 0 { // an empty page is unusable
		return "", nil, true, errors.New("response body was empty") // worth retrying
	}

	return string(responseBodyBytes), httpResponse.Request.URL, false, nil // success: body text and final URL
}

// extractAnchorElements returns every <a href="..."> element of pageContent with its decoded href and visible text.
func extractAnchorElements(pageContent string) []anchorElement {
	contentWithoutComments := htmlCommentPattern.ReplaceAllString(pageContent, "")             // drop commented-out markup
	allAnchorMatches := anchorElementPattern.FindAllStringSubmatch(contentWithoutComments, -1) // find every anchor element
	anchorElements := make([]anchorElement, 0, len(allAnchorMatches))                          // result slice sized for all anchors

	for _, anchorMatch := range allAnchorMatches { // check each anchor
		hrefMatch := hrefAttributePattern.FindStringSubmatch(anchorMatch[1]) // look for an href attribute
		if hrefMatch == nil {                                                // anchors without href (such as named anchors) are skipped
			continue // go to the next anchor
		}

		rawHyperlinkReference := firstNonEmptyString(hrefMatch[1], hrefMatch[2], hrefMatch[3]) // pick the captured quoting style
		rawHyperlinkReference = strings.TrimSpace(html.UnescapeString(rawHyperlinkReference))  // decode entities and trim spaces
		if rawHyperlinkReference == "" {                                                       // ignore empty hrefs
			continue // go to the next anchor
		}

		visibleText := htmlTagPattern.ReplaceAllString(anchorMatch[2], " ") // remove tags nested inside the anchor
		visibleText = strings.TrimSpace(html.UnescapeString(visibleText))   // decode entities and trim spaces

		anchorElements = append(anchorElements, anchorElement{hyperlinkReference: rawHyperlinkReference, visibleText: visibleText}) // keep this anchor
	}

	debugLogger.Printf("found %d anchor elements with an href attribute", len(anchorElements)) // log how many were found
	return anchorElements                                                                      // hand back the anchors
}

// resolveHyperlinkReference turns a raw href into an absolute URL without a fragment.
func resolveHyperlinkReference(baseURL *url.URL, rawHyperlinkReference string) (*url.URL, error) {
	parsedReference, parseError := url.Parse(rawHyperlinkReference) // parse the href
	if parseError != nil {                                          // if it is malformed
		return nil, parseError // report the problem
	}
	resolvedURL := baseURL.ResolveReference(parsedReference) // turn relative links into absolute ones
	resolvedURL.Fragment = ""                                // ignore the #fragment part
	resolvedURL.RawFragment = ""                             // ignore its raw form too
	return resolvedURL, nil                                  // hand back the absolute URL
}

// extractSupportProductURLs finds all valid support product links in pageContent,
// resolves them against baseURL and returns them de-duplicated and sorted.
func extractSupportProductURLs(pageContent string, baseURL *url.URL) []string {
	uniqueURLs := make(map[string]struct{}) // set used to remove duplicates

	for _, linkedAnchor := range extractAnchorElements(pageContent) { // check each anchor
		resolvedURL, resolveError := resolveHyperlinkReference(baseURL, linkedAnchor.hyperlinkReference) // make the link absolute
		if resolveError != nil {                                                                         // if it is malformed
			errorLogger.Printf("skipping unparsable href %q: %v", linkedAnchor.hyperlinkReference, resolveError) // log and skip
			continue                                                                                             // go to the next anchor
		}

		if !isValidSupportProductURL(resolvedURL) { // keep only support product links
			debugLogger.Printf("ignoring non-matching link: %s", linkedAnchor.hyperlinkReference) // log the skipped link
			continue                                                                              // go to the next anchor
		}

		escapedProductPath := strings.TrimSuffix(resolvedURL.EscapedPath(), "/") // keep original encoding, drop trailing slash
		canonicalURL := "https://" + canonicalHostname + escapedProductPath      // always output https on the canonical host

		if _, alreadyCollected := uniqueURLs[canonicalURL]; alreadyCollected { // check for a duplicate
			debugLogger.Printf("duplicate skipped: %s", canonicalURL) // log the duplicate
			continue                                                  // go to the next anchor
		}
		uniqueURLs[canonicalURL] = struct{}{}            // remember this URL
		debugLogger.Printf("accepted: %s", canonicalURL) // log the accepted URL
	}

	return sortedKeys(uniqueURLs) // hand back the sorted result
}

// extractDocumentationPageURLs finds all links to /documentation/{document-slug} pages on the Ezurio website.
func extractDocumentationPageURLs(pageContent string, baseURL *url.URL) []string {
	uniqueURLs := make(map[string]struct{}) // set used to remove duplicates

	for _, linkedAnchor := range extractAnchorElements(pageContent) { // check each anchor
		resolvedURL, resolveError := resolveHyperlinkReference(baseURL, linkedAnchor.hyperlinkReference) // make the link absolute
		if resolveError != nil {                                                                         // skip malformed links silently, they are logged elsewhere
			continue // go to the next anchor
		}

		if !isEzurioHostname(resolvedURL.Hostname()) { // only follow links on the Ezurio website
			continue // go to the next anchor
		}

		documentSlug := strings.Trim(strings.TrimPrefix(resolvedURL.Path, documentationPathPrefix), "/") // isolate what follows /documentation/
		if !strings.HasPrefix(resolvedURL.Path, documentationPathPrefix) || documentSlug == "" {         // must be a real documentation page
			continue // go to the next anchor
		}

		canonicalURL := "https://" + canonicalHostname + strings.TrimSuffix(resolvedURL.EscapedPath(), "/") // always https on the canonical host
		uniqueURLs[canonicalURL] = struct{}{}                                                               // remember this page
	}

	return sortedKeys(uniqueURLs) // hand back the sorted result
}

// extractPDFURLs finds every PDF link in pageContent. It checks link and embed attributes,
// anchors labelled "Download PDF", and any absolute .pdf URL written anywhere in the text.
func extractPDFURLs(pageContent string, baseURL *url.URL) []string {
	uniquePDFURLs := make(map[string]struct{})                                                               // set used to remove duplicates
	normalizedContent := strings.ReplaceAll(htmlCommentPattern.ReplaceAllString(pageContent, ""), `\/`, `/`) // remove comments and unescape JSON slashes

	for _, attributeMatch := range quotedAttributeValuePattern.FindAllStringSubmatch(normalizedContent, -1) { // check every href/src/data attribute
		registerPDFCandidate(uniquePDFURLs, baseURL, firstNonEmptyString(attributeMatch[1], attributeMatch[2]), true) // accept it only if it ends in .pdf
	}

	for _, linkedAnchor := range extractAnchorElements(normalizedContent) { // check every anchor
		lowercaseVisibleText := strings.ToLower(linkedAnchor.visibleText)                                        // compare text case-insensitively
		if strings.Contains(lowercaseVisibleText, "download") && strings.Contains(lowercaseVisibleText, "pdf") { // a "Download PDF" button
			registerPDFCandidate(uniquePDFURLs, baseURL, linkedAnchor.hyperlinkReference, false) // accept it even without a .pdf extension
		}
	}

	for _, absoluteURLText := range absolutePDFURLPattern.FindAllString(normalizedContent, -1) { // check every absolute .pdf URL in the text
		registerPDFCandidate(uniquePDFURLs, baseURL, absoluteURLText, true) // accept it only if it ends in .pdf
	}

	return sortedKeys(uniquePDFURLs) // hand back the sorted result
}

// registerPDFCandidate resolves a raw reference and adds it to the set when it qualifies as a PDF link.
func registerPDFCandidate(uniquePDFURLs map[string]struct{}, baseURL *url.URL, rawReference string, requirePDFExtension bool) {
	trimmedReference := strings.TrimSpace(html.UnescapeString(rawReference)) // decode entities and trim spaces
	if trimmedReference == "" {                                              // ignore empty values
		return // nothing to register
	}

	resolvedURL, resolveError := resolveHyperlinkReference(baseURL, trimmedReference) // make the link absolute
	if resolveError != nil {                                                          // ignore malformed values
		return // nothing to register
	}

	if resolvedURL.RawQuery != "" { // links such as /account?file_download=a b (1).pdf contain spaces
		resolvedURL.RawQuery = resolvedURL.Query().Encode() // re-encode the query so the HTTP request line is valid
	}
	if resolvedURL.Path == "/account" && resolvedURL.Query().Get(fileDownloadQueryParameter) == "" { // plain sign-in links are not downloads
		return // nothing to register
	}
	if resolvedURL.Scheme != "http" && resolvedURL.Scheme != "https" { // reject mailto:, javascript:, data: and similar
		return // not a web link
	}

	if requirePDFExtension && !hasPDFExtension(resolvedURL) { // plain attributes must end in .pdf
		return // not a PDF link
	}

	if _, alreadyRegistered := uniquePDFURLs[resolvedURL.String()]; !alreadyRegistered { // log each new PDF link only once
		debugLogger.Printf("PDF link found: %s", resolvedURL.String()) // log the discovery
	}
	uniquePDFURLs[resolvedURL.String()] = struct{}{} // remember this PDF link
}

// hasPDFExtension reports whether the URL path, or the file_download query value, ends in .pdf (ignoring case).
func hasPDFExtension(candidateURL *url.URL) bool {
	if strings.HasSuffix(strings.ToLower(candidateURL.Path), ".pdf") { // a direct link such as /files/datasheet.pdf
		return true // it is a PDF link
	}
	downloadedFileName := candidateURL.Query().Get(fileDownloadQueryParameter) // a gated link such as /account?file_download=datasheet.pdf
	return strings.HasSuffix(strings.ToLower(downloadedFileName), ".pdf")      // true when the requested file is a PDF
}

// applySessionCookie adds the optional session cookie to requests that go to the Ezurio website only.
func applySessionCookie(httpRequest *http.Request) {
	if sessionCookieHeaderValue == "" { // no cookie configured
		return // nothing to add
	}
	if !isEzurioHostname(httpRequest.URL.Hostname()) { // never send the cookie to other websites
		return // leave the request untouched
	}
	httpRequest.Header.Set("Cookie", sessionCookieHeaderValue) // identify the logged-in session
}

// containsSignInRequirement reports whether a page tells anonymous visitors to sign in.
func containsSignInRequirement(pageContent string) bool {
	lowercaseContent := strings.ToLower(pageContent)                                                                           // compare case-insensitively
	return strings.Contains(lowercaseContent, "sign in required") || strings.Contains(lowercaseContent, "sign in to continue") // the two phrases Ezurio uses
}

// isValidSupportProductURL reports whether candidateURL is an Ezurio link of the
// form /support/series/{product-slug}, where the slug is one non-empty path segment.
func isValidSupportProductURL(candidateURL *url.URL) bool {
	if candidateURL.Scheme != "http" && candidateURL.Scheme != "https" { // reject mailto:, javascript:, and similar
		return false // not a web link
	}

	if !isEzurioHostname(candidateURL.Hostname()) { // reject links to other websites
		return false // wrong website
	}

	if !strings.HasPrefix(candidateURL.Path, supportSeriesPathPrefix) { // path must start with /support/series/
		return false // wrong section of the website
	}

	productSlug := strings.TrimSuffix(strings.TrimPrefix(candidateURL.Path, supportSeriesPathPrefix), "/") // isolate the slug
	if productSlug == "" || strings.Contains(productSlug, "/") {                                           // slug must be a single segment
		return false // empty slug or deeper path
	}

	return true // every check passed
}

// isEzurioHostname reports whether hostname belongs to the Ezurio website.
func isEzurioHostname(hostname string) bool {
	lowercaseHostname := strings.ToLower(hostname)                                                   // compare hostnames case-insensitively
	return lowercaseHostname == canonicalHostname || lowercaseHostname == acceptedHostnameWithoutWww // accept with and without www
}

// collectLinksFromSeriesPages visits every product support page and returns the
// documentation page URLs and the direct PDF URLs found on them.
func collectLinksFromSeriesPages(pageHTTPClient *http.Client, seriesPageURLs []string) ([]string, []string) {
	documentationPageSet := newConcurrentStringSet() // collects documentation page URLs
	directPDFSet := newConcurrentStringSet()         // collects PDF URLs linked directly from series pages
	var visitedPageCount atomic.Int64                // counts finished pages for progress logging

	informationLogger.Printf("visiting %d product support pages with %d workers", len(seriesPageURLs), concurrentWorkerCount) // announce the stage

	processConcurrently(seriesPageURLs, func(seriesPageURL string) { // handle each series page on a worker
		pageContent, finalPageURL, fetchError := fetchSupportPage(pageHTTPClient, seriesPageURL) // download the page
		currentPosition := visitedPageCount.Add(1)                                               // count this page
		if fetchError != nil {                                                                   // if it could not be downloaded
			errorLogger.Printf("[%d/%d] could not visit %s: %v", currentPosition, len(seriesPageURLs), seriesPageURL, fetchError) // log and move on
			return                                                                                                                // skip this page
		}

		foundDocumentationPageURLs := extractDocumentationPageURLs(pageContent, finalPageURL) // find documentation pages
		foundPDFURLs := extractPDFURLs(pageContent, finalPageURL)                             // find direct PDF links
		documentationPageSet.addAll(foundDocumentationPageURLs)                               // remember the documentation pages
		directPDFSet.addAll(foundPDFURLs)                                                     // remember the PDF links

		informationLogger.Printf("[%d/%d] %s: %d documentation pages, %d PDF links", currentPosition, len(seriesPageURLs), seriesPageURL, len(foundDocumentationPageURLs), len(foundPDFURLs)) // log progress
	})

	return documentationPageSet.sortedValues(), directPDFSet.sortedValues() // hand back both sorted lists
}

// collectPDFLinksFromDocumentationPages visits every documentation page and returns the PDF URLs found on them.
func collectPDFLinksFromDocumentationPages(pageHTTPClient *http.Client, documentationPageURLs []string) []string {
	pdfSet := newConcurrentStringSet()       // collects PDF URLs
	var visitedPageCount atomic.Int64        // counts finished pages for progress logging
	var signInRequiredPageCount atomic.Int64 // counts pages that ask anonymous visitors to sign in

	informationLogger.Printf("visiting %d documentation pages with %d workers", len(documentationPageURLs), concurrentWorkerCount) // announce the stage

	processConcurrently(documentationPageURLs, func(documentationPageURL string) { // handle each documentation page on a worker
		pageContent, finalPageURL, fetchError := fetchSupportPage(pageHTTPClient, documentationPageURL) // download the page
		currentPosition := visitedPageCount.Add(1)                                                      // count this page
		if fetchError != nil {                                                                          // if it could not be downloaded
			errorLogger.Printf("[%d/%d] could not visit %s: %v", currentPosition, len(documentationPageURLs), documentationPageURL, fetchError) // log and move on
			return                                                                                                                              // skip this page
		}

		foundPDFURLs := extractPDFURLs(pageContent, finalPageURL) // find PDF links on the page
		if containsSignInRequirement(pageContent) {               // the page hides its content or downloads behind a login
			signInRequiredPageCount.Add(1) // count it
		}
		pdfSet.addAll(foundPDFURLs) // remember them

		debugLogger.Printf("[%d/%d] %s: %d PDF links", currentPosition, len(documentationPageURLs), documentationPageURL, len(foundPDFURLs)) // log progress
		if currentPosition%50 == 0 {                                                                                                         // report progress every 50 pages
			informationLogger.Printf("visited %d of %d documentation pages", currentPosition, len(documentationPageURLs)) // log the milestone
		}
	})

	informationLogger.Printf("%d of %d documentation pages require signing in", signInRequiredPageCount.Load(), len(documentationPageURLs)) // explain missing PDFs
	return pdfSet.sortedValues()                                                                                                            // hand back the sorted list
}

// downloadAllPDFFiles downloads every PDF into the PDFs folder and logs a summary.
func downloadAllPDFFiles(downloadHTTPClient *http.Client, pdfURLs []string) {
	if createError := os.MkdirAll(pdfOutputDirectory, pdfOutputDirectoryPermissions); createError != nil { // make sure the folder exists
		errorLogger.Fatalf("could not create the %s folder: %v", pdfOutputDirectory, createError) // exit when it cannot be created
	}

	destinationFileNames := planPDFFileNames(pdfURLs)                           // decide a unique file name for every PDF URL
	var downloadedCount, skippedCount, failedCount, processedCount atomic.Int64 // thread-safe counters for the summary

	informationLogger.Printf("downloading %d PDF files into %s/ with %d workers", len(pdfURLs), pdfOutputDirectory, concurrentWorkerCount) // announce the stage

	processConcurrently(pdfURLs, func(pdfURL string) { // handle each PDF on a worker
		destinationPath := filepath.Join(pdfOutputDirectory, destinationFileNames[pdfURL])            // full path of the target file
		alreadyExisted, downloadError := downloadPDFFile(downloadHTTPClient, pdfURL, destinationPath) // download (or skip) the file
		currentPosition := processedCount.Add(1)                                                      // count this file

		switch { // report the outcome
		case downloadError != nil: // the download failed
			failedCount.Add(1)                                                                                            // count the failure
			errorLogger.Printf("[%d/%d] failed to download %s: %v", currentPosition, len(pdfURLs), pdfURL, downloadError) // log the failure
		case alreadyExisted: // the file was downloaded in an earlier run
			skippedCount.Add(1)                                                                                              // count the skip
			informationLogger.Printf("[%d/%d] already present, skipped: %s", currentPosition, len(pdfURLs), destinationPath) // log the skip
		default: // the file was downloaded just now
			downloadedCount.Add(1)                                                                            // count the success
			informationLogger.Printf("[%d/%d] downloaded %s", currentPosition, len(pdfURLs), destinationPath) // log the success
		}
	})

	informationLogger.Printf("finished: %d downloaded, %d already present, %d failed", downloadedCount.Load(), skippedCount.Load(), failedCount.Load()) // log the summary
}

// planPDFFileNames assigns each PDF URL a file name. When two different URLs would
// get the same name, a short hash of the URL is added so no file overwrites another.
func planPDFFileNames(pdfURLs []string) map[string]string {
	fileNamesByURL := make(map[string]string, len(pdfURLs)) // result: URL to file name
	ownerByLowercaseName := make(map[string]string)         // which URL already owns a name (case-insensitive)

	for _, pdfURL := range pdfURLs { // the input is sorted, so the plan is stable between runs
		fileName := buildFileNameFromURL(pdfURL)                                 // name derived from the URL
		ownerURL, nameIsTaken := ownerByLowercaseName[strings.ToLower(fileName)] // check for a collision
		if nameIsTaken && ownerURL != pdfURL {                                   // another URL already uses this name
			fileExtension := filepath.Ext(fileName)                                            // for example ".pdf"
			fileStem := strings.TrimSuffix(fileName, fileExtension)                            // name without extension
			fileName = fmt.Sprintf("%s-%s%s", fileStem, buildShortHash(pdfURL), fileExtension) // make the name unique
		}
		ownerByLowercaseName[strings.ToLower(fileName)] = pdfURL // reserve the name
		fileNamesByURL[pdfURL] = fileName                        // record the plan
	}

	return fileNamesByURL // hand back the plan
}

// buildFileNameFromURL derives a safe file name ending in .pdf from the last segment of the URL path.
func buildFileNameFromURL(pdfURL string) string {
	parsedURL, parseError := url.Parse(pdfURL) // parse the URL
	if parseError != nil {                     // fall back for malformed URLs
		return "document-" + buildShortHash(pdfURL) + ".pdf" // hash-based name
	}

	lastPathSegment := parsedURL.Query().Get(fileDownloadQueryParameter) // gated downloads carry the file name in the query
	if lastPathSegment == "" {                                           // normal links carry it in the path
		lastPathSegment = path.Base(parsedURL.Path) // use the last path segment
	}
	lastPathSegment = illegalFileNameCharactersPattern.ReplaceAllString(lastPathSegment, "_") // remove characters the file system forbids
	lastPathSegment = strings.Trim(lastPathSegment, " .")                                     // remove leading and trailing spaces and dots
	if lastPathSegment == "" || lastPathSegment == "/" || lastPathSegment == "_" {            // nothing usable in the path
		lastPathSegment = "document-" + buildShortHash(pdfURL) // hash-based name
	}

	if !strings.HasSuffix(strings.ToLower(lastPathSegment), ".pdf") { // make sure the extension is .pdf
		lastPathSegment += ".pdf" // append it
	}

	return lastPathSegment // hand back the file name
}

// buildShortHash returns the first eight hexadecimal characters of the SHA-256 hash of text.
func buildShortHash(text string) string {
	hashBytes := sha256.Sum256([]byte(text))    // hash the text
	return hex.EncodeToString(hashBytes[:])[:8] // shorten it for use in file names
}

// downloadPDFFile saves pdfURL to destinationPath. It returns true when the file already existed
// and was skipped. Temporary failures are retried, and responses that are not real PDFs are rejected.
func downloadPDFFile(downloadHTTPClient *http.Client, pdfURL string, destinationPath string) (bool, error) {
	if existingFileInformation, statError := os.Stat(destinationPath); statError == nil && existingFileInformation.Size() > 0 { // a previous run saved it
		debugLogger.Printf("skipping existing file %s", destinationPath) // log the skip
		return true, nil                                                 // nothing to download
	}

	var lastError error // remembers the most recent failure

	for attemptNumber := 1; attemptNumber <= maximumFetchAttempts; attemptNumber++ { // try up to the maximum
		debugLogger.Printf("download attempt %d of %d for %s", attemptNumber, maximumFetchAttempts, pdfURL) // log the attempt

		shouldRetry, downloadError := attemptPDFDownload(downloadHTTPClient, pdfURL, destinationPath) // one attempt
		if downloadError == nil {                                                                     // on success
			return false, nil // the file was downloaded
		}

		lastError = downloadError                                                                // keep the error for the final message
		debugLogger.Printf("attempt %d for %s failed: %v", attemptNumber, pdfURL, downloadError) // log the failure

		if !shouldRetry || attemptNumber == maximumFetchAttempts { // stop when retrying is pointless or attempts are used up
			break // leave the retry loop
		}
		time.Sleep(time.Duration(attemptNumber) * baseDelayBetweenAttempts) // wait longer after each failure
	}

	return false, lastError // report the final failure
}

// attemptPDFDownload performs one download attempt. The boolean result says whether a retry could help.
func attemptPDFDownload(downloadHTTPClient *http.Client, pdfURL string, destinationPath string) (bool, error) {
	httpRequest, requestCreationError := http.NewRequest(http.MethodGet, pdfURL, nil) // build the request
	if requestCreationError != nil {                                                  // if the address is invalid
		return false, fmt.Errorf("building request: %w", requestCreationError) // retrying cannot help
	}
	httpRequest.Header.Set("User-Agent", userAgentHeaderValue)                                   // identify as a regular browser
	httpRequest.Header.Set("Accept", "application/pdf,application/octet-stream;q=0.9,*/*;q=0.8") // prefer PDF content
	httpRequest.Header.Set("Accept-Language", "en-US,en;q=0.9")                                  // ask for English content
	applySessionCookie(httpRequest)                                                              // send the logged-in session cookie to Ezurio hosts, if configured

	httpResponse, requestError := downloadHTTPClient.Do(httpRequest) // send the request and follow redirects
	if requestError != nil {                                         // on a network failure or timeout
		return true, fmt.Errorf("sending request: %w", requestError) // worth retrying
	}
	defer httpResponse.Body.Close() // always release the connection

	if httpResponse.StatusCode != http.StatusOK { // any status other than 200 OK is a problem
		isTemporaryFailure := httpResponse.StatusCode == http.StatusTooManyRequests || httpResponse.StatusCode >= http.StatusInternalServerError // rate limit or server error
		return isTemporaryFailure, fmt.Errorf("unexpected HTTP status %s", httpResponse.Status)                                                  // only temporary failures are retried
	}

	responseReader := bufio.NewReader(httpResponse.Body)                     // buffered reader that allows looking ahead
	leadingBytes, peekError := responseReader.Peek(pdfSignatureSearchLength) // look at the first bytes without consuming them
	if peekError != nil && !errors.Is(peekError, io.EOF) {                   // a short file ends with EOF, anything else is a real error
		return true, fmt.Errorf("reading response start: %w", peekError) // worth retrying
	}
	if len(leadingBytes) == 0 { // an empty response is unusable
		return true, errors.New("response body was empty") // worth retrying
	}
	if !bytes.Contains(leadingBytes, []byte(pdfSignature)) { // the content is not a PDF, for example a sign-in page
		return false, fmt.Errorf("response is not a PDF (content-type %q); the file may require signing in", httpResponse.Header.Get("Content-Type")) // retrying cannot help
	}

	temporaryFile, temporaryFileError := os.CreateTemp(pdfOutputDirectory, "download-*.part") // temporary file in the target folder
	if temporaryFileError != nil {                                                            // if the file system refuses
		return false, fmt.Errorf("creating temporary file: %w", temporaryFileError) // retrying cannot help
	}
	temporaryFilePath := temporaryFile.Name() // remember the temporary path for cleanup and renaming

	_, copyError := io.Copy(temporaryFile, responseReader) // stream the PDF to disk
	closeError := temporaryFile.Close()                    // flush and close the file
	if copyError != nil {                                  // the transfer broke midway
		os.Remove(temporaryFilePath)                                 // delete the partial file
		return true, fmt.Errorf("saving PDF content: %w", copyError) // worth retrying
	}
	if closeError != nil { // the file could not be finalised
		os.Remove(temporaryFilePath)                                       // delete the partial file
		return false, fmt.Errorf("closing temporary file: %w", closeError) // retrying cannot help
	}

	if renameError := os.Rename(temporaryFilePath, destinationPath); renameError != nil { // move the finished file to its final name
		os.Remove(temporaryFilePath)                                          // delete the temporary file
		return false, fmt.Errorf("renaming downloaded file: %w", renameError) // retrying cannot help
	}

	return false, nil // success
}

// processConcurrently runs processItem for every item using a fixed number of workers.
func processConcurrently(items []string, processItem func(item string)) {
	itemChannel := make(chan string) // hands items to the workers
	var workerGroup sync.WaitGroup   // lets us wait for all workers to finish

	for workerNumber := 0; workerNumber < concurrentWorkerCount; workerNumber++ { // start the workers
		workerGroup.Add(1) // register one more worker
		go func() {        // run the worker in its own goroutine
			defer workerGroup.Done()        // report completion when the channel is drained
			for item := range itemChannel { // take items until the channel is closed
				processItem(item)              // do the actual work
				time.Sleep(pauseAfterEachItem) // be polite to the server
			}
		}()
	}

	for _, item := range items { // feed all items to the workers
		itemChannel <- item // blocks until a worker is free
	}
	close(itemChannel) // tell the workers there is nothing more
	workerGroup.Wait() // wait until every worker has finished
}

// newConcurrentStringSet creates an empty set.
func newConcurrentStringSet() *concurrentStringSet {
	return &concurrentStringSet{values: make(map[string]struct{})} // allocate the map
}

// addAll stores every given string in the set.
func (stringSet *concurrentStringSet) addAll(newValues []string) {
	stringSet.mutex.Lock()               // block other goroutines while writing
	defer stringSet.mutex.Unlock()       // release the lock when done
	for _, newValue := range newValues { // go through the new strings
		stringSet.values[newValue] = struct{}{} // store each one
	}
}

// sortedValues returns the stored strings in alphabetical order.
func (stringSet *concurrentStringSet) sortedValues() []string {
	stringSet.mutex.Lock()              // block other goroutines while reading
	defer stringSet.mutex.Unlock()      // release the lock when done
	return sortedKeys(stringSet.values) // hand back the sorted strings
}

// mergeUniqueSortedStrings combines several lists into one sorted list without duplicates.
func mergeUniqueSortedStrings(stringLists ...[]string) []string {
	uniqueValues := make(map[string]struct{}) // set used to remove duplicates
	for _, stringList := range stringLists {  // go through each list
		for _, value := range stringList { // go through each value
			uniqueValues[value] = struct{}{} // store it once
		}
	}
	return sortedKeys(uniqueValues) // hand back the sorted result
}

// sortedKeys returns the keys of a set in alphabetical order.
func sortedKeys(stringSet map[string]struct{}) []string {
	keys := make([]string, 0, len(stringSet)) // slice sized for all keys
	for key := range stringSet {              // copy the keys
		keys = append(keys, key) // add one key
	}
	sort.Strings(keys) // sort for stable output
	return keys        // hand back the result
}

// firstNonEmptyString returns the first argument that is not empty.
func firstNonEmptyString(candidateStrings ...string) string {
	for _, candidateString := range candidateStrings { // go through the candidates in order
		if candidateString != "" { // found a usable value
			return candidateString // return it
		}
	}
	return "" // none were usable
}
