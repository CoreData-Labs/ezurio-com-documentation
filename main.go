// Command ezurio-scraper discovers every product support URL of the form
// /support/series/{product-slug} that is linked from https://www.ezurio.com/support.
// It uses only the Go standard library. Logs go to standard error and the
// extracted URLs go to standard output.
package main // this file builds an executable program

import ( // standard library packages used by this program
	"errors"   // creates simple error values
	"fmt"      // formats text and prints the final URLs
	"html"     // decodes HTML entities such as &amp; inside href values
	"io"       // reads response bodies and provides the discard writer
	"log"      // provides the console loggers
	"net/http" // performs the HTTP request
	"net/url"  // parses and resolves URLs
	"os"       // provides standard error output
	"regexp"   // finds anchor elements in the HTML text
	"sort"     // sorts the final URL list
	"strings"  // provides string helper functions
	"time"     // provides timeouts and retry delays
) // end of the import list

const ( // fixed configuration values
	supportPageAddress         = "https://www.ezurio.com/support"                                                                                 // page where scraping starts
	supportSeriesPathPrefix    = "/support/series/"                                                                                               // path prefix every product support link must have
	canonicalHostname          = "www.ezurio.com"                                                                                                 // hostname written in every printed URL
	acceptedHostnameWithoutWww = "ezurio.com"                                                                                                     // alternative hostname accepted in links
	maximumFetchAttempts       = 3                                                                                                                // how many times a failing download is tried
	baseDelayBetweenAttempts   = 2 * time.Second                                                                                                  // retry delay, multiplied by the attempt number
	httpRequestTimeout         = 30 * time.Second                                                                                                 // maximum time allowed for one HTTP request
	debugLoggingEnabled        = true                                                                                                             // set to false to hide debug log lines
	userAgentHeaderValue       = "Mozilla/5.0 (X11; CrOS x86_64 14541.0.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/150.0.0.0 Safari/537.36" // value sent in the User-Agent header
) // end of the constants

// anchorHyperlinkReferencePattern captures the href value of every <a ...> element.
// The value may be double-quoted, single-quoted or unquoted. The (?:[^>]*?\s)?
// prefix prevents attributes such as "data-href" from being treated as "href".
var anchorHyperlinkReferencePattern = regexp.MustCompile(`(?is)<a\b(?:[^>]*?\s)?href\s*=\s*(?:"([^"]*)"|'([^']*)'|([^\s"'>]+))`)

// htmlCommentPattern matches HTML comments so commented-out links are ignored.
var htmlCommentPattern = regexp.MustCompile(`(?s)<!--.*?-->`)

var ( // console loggers
	informationLogger = log.New(os.Stderr, "INFO  ", log.LstdFlags) // logs normal progress messages
	errorLogger       = log.New(os.Stderr, "ERROR ", log.LstdFlags) // logs failures and warnings
	debugLogger       = createDebugLogger()                         // logs detailed diagnostics when enabled
) // end of the loggers

// createDebugLogger returns a logger that prints only when debug logging is enabled.
func createDebugLogger() *log.Logger {
	debugOutputWriter := io.Discard // by default, throw debug output away
	if debugLoggingEnabled {        // when debug logging is switched on
		debugOutputWriter = os.Stderr // send debug output to the console
	}
	return log.New(debugOutputWriter, "DEBUG ", log.LstdFlags|log.Lmicroseconds) // build the logger
}

// main runs the whole scraping workflow.
func main() {
	informationLogger.Printf("starting Ezurio support product URL scraper") // announce the start
	debugLogger.Printf("target support page: %s", supportPageAddress)       // log the target address

	httpClient := &http.Client{Timeout: httpRequestTimeout} // client that gives up after the timeout

	pageContent, finalPageURL, fetchError := fetchSupportPage(httpClient, supportPageAddress) // download the page
	if fetchError != nil {                                                                    // if the download failed
		errorLogger.Fatalf("could not fetch the support page: %v", fetchError) // log the error and exit with code 1
	}
	informationLogger.Printf("fetched support page (%d bytes) from %s", len(pageContent), finalPageURL) // log the download size

	productSupportURLs := extractSupportProductURLs(pageContent, finalPageURL) // find all product support URLs
	if len(productSupportURLs) == 0 {                                          // if nothing matched
		errorLogger.Fatalf("no support product URLs were found; the links may be rendered by JavaScript or the page structure may have changed") // exit with an explanation
	}

	informationLogger.Printf("found %d unique support product URLs", len(productSupportURLs)) // log the final count
	for _, productSupportURL := range productSupportURLs {                                    // go through each URL
		fmt.Println(productSupportURL) // print it to standard output
	}
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

		lastError = fetchError                                                 // keep the error for the final message
		errorLogger.Printf("attempt %d failed: %v", attemptNumber, fetchError) // log the failure

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
		errorLogger.Printf("warning: unexpected content type %q, parsing anyway", contentType) // log the warning
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

// extractSupportProductURLs finds all valid support product links in pageContent,
// resolves them against baseURL and returns them de-duplicated and sorted.
func extractSupportProductURLs(pageContent string, baseURL *url.URL) []string {
	contentWithoutComments := htmlCommentPattern.ReplaceAllString(pageContent, "")                           // drop commented-out markup
	allHyperlinkMatches := anchorHyperlinkReferencePattern.FindAllStringSubmatch(contentWithoutComments, -1) // find every anchor href
	debugLogger.Printf("found %d anchor elements with an href attribute", len(allHyperlinkMatches))          // log how many were found

	uniqueURLs := make(map[string]struct{}) // set used to remove duplicates

	for _, hyperlinkMatch := range allHyperlinkMatches { // check each anchor
		rawHyperlinkReference := firstNonEmptyString(hyperlinkMatch[1], hyperlinkMatch[2], hyperlinkMatch[3]) // pick the captured quoting style
		rawHyperlinkReference = strings.TrimSpace(html.UnescapeString(rawHyperlinkReference))                 // decode entities and trim spaces
		if rawHyperlinkReference == "" {                                                                      // ignore empty hrefs
			continue // go to the next anchor
		}

		parsedReference, parseError := url.Parse(rawHyperlinkReference) // parse the href
		if parseError != nil {                                          // if it is malformed
			errorLogger.Printf("skipping unparsable href %q: %v", rawHyperlinkReference, parseError) // log and skip
			continue                                                                                 // go to the next anchor
		}

		resolvedURL := baseURL.ResolveReference(parsedReference) // turn relative links into absolute ones
		if !isValidSupportProductURL(resolvedURL) {              // keep only support product links
			debugLogger.Printf("ignoring non-matching link: %s", rawHyperlinkReference) // log the skipped link
			continue                                                                    // go to the next anchor
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

	sortedURLs := make([]string, 0, len(uniqueURLs)) // slice sized for all unique URLs
	for uniqueURL := range uniqueURLs {              // copy the set into the slice
		sortedURLs = append(sortedURLs, uniqueURL) // add one URL
	}
	sort.Strings(sortedURLs) // sort for stable output
	return sortedURLs        // hand back the result
}

// isValidSupportProductURL reports whether candidateURL is an Ezurio link of the
// form /support/series/{product-slug}, where the slug is one non-empty path segment.
func isValidSupportProductURL(candidateURL *url.URL) bool {
	if candidateURL.Scheme != "http" && candidateURL.Scheme != "https" { // reject mailto:, javascript:, and similar
		return false // not a web link
	}

	hostname := strings.ToLower(candidateURL.Hostname())                         // compare hostnames case-insensitively
	if hostname != canonicalHostname && hostname != acceptedHostnameWithoutWww { // reject links to other websites
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

// firstNonEmptyString returns the first argument that is not empty.
func firstNonEmptyString(candidateStrings ...string) string {
	for _, candidateString := range candidateStrings { // go through the candidates in order
		if candidateString != "" { // found a usable value
			return candidateString // return it
		}
	}
	return "" // none were usable
}
