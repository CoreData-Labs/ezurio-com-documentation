# 📚 Ezurio Documentation Archive

### A community-friendly library of datasheets, user guides, application notes and more for Ezurio wireless modules, antennas and IoT hardware 📡

👋 **Welcome!** Whether you're a student, hobbyist, repair technician, firmware engineer, researcher or just curious about how wireless devices work, this repository was built for **you**. Everything here is about one thing: **making technical documentation easy to find, easy to keep, and easy to learn from.** 💡

---

## 📑 Table of Contents

1. [What is this repository?](#-what-is-this-repository)
2. [Who is Ezurio?](#-who-is-ezurio)
3. [What's inside the PDFs?](#-whats-inside-the-pdfs)
4. [Product families covered](#-product-families-covered)
5. [Repository layout](#-repository-layout)
6. [How to find the document you need](#-how-to-find-the-document-you-need)
7. [How these documents benefit you](#-how-these-documents-benefit-you)
8. [Education and learning](#-education-and-learning)
9. [Repairs and troubleshooting](#-repairs-and-troubleshooting)
10. [DIY and maker projects](#-diy-and-maker-projects)
11. [Professional engineering](#-professional-engineering)
12. [AI and machine learning uses](#-ai-and-machine-learning-uses)
13. [Benefits for your community](#-benefits-for-your-community)
14. [Reading a datasheet: a quick guide](#-reading-a-datasheet-a-quick-guide)
15. [Glossary](#-glossary)
16. [Keeping the archive fresh](#-keeping-the-archive-fresh)
17. [Contributing](#-contributing)
18. [Legal notice and licensing](#%EF%B8%8F-legal-notice-and-licensing)
19. [FAQ](#-faq)

---

## 🌟 What is this repository?

This repository is an **organized archive of PDF documentation** gathered from the public support and documentation pages at [ezurio.com](https://www.ezurio.com/support). Each product has a support page, and each support page links to documentation pages that offer downloadable PDFs. Those PDFs are collected here into one folder so you can:

- 🔍 **Search** them all at once with your computer's file search or `grep`
- 💾 **Keep an offline copy** for when the internet (or the website) isn't available
- 📖 **Read** them on any device
- 🧠 **Study** them, build on them and learn from them
- 🤝 **Share** the knowledge with your classroom, makerspace or team

Hardware documentation has a bad habit of disappearing: links break, products are retired, websites get redesigned. An archive means the knowledge sticks around. 🧷

---

## 🏢 Who is Ezurio?

Ezurio is a company that designs and supplies **wireless connectivity and IoT hardware**. It is the company previously known as **Laird Connectivity**, so older forum posts, tutorials and product listings may use the Laird name for the same hardware. 🔁

Ezurio's catalog centers on things that let devices talk to each other without wires:

- 📶 **Bluetooth® and Bluetooth Low Energy (BLE)** modules
- 🌐 **Wi-Fi** modules
- 📡 **Cellular and LPWAN** connectivity
- 🔋 **Long-range and low-power radios** for sensors and trackers
- 📻 **Antennas** for many frequency bands
- 🧰 **Development kits, evaluation boards and software tooling**
- 🏭 **Embedded wireless solutions** for industrial, medical, smart-home, logistics and consumer products

Modules like these quietly power things around you: fitness trackers, medical monitors, asset trackers, smart locks, industrial sensors, point-of-sale terminals, wearables and more. 🌍

---

## 📄 What's inside the PDFs?

The documents in this archive come in many flavors. Not every product has every type, but here is what you can typically expect to find:

| 📘 Document type                           | 📝 What it tells you                                                                                                    | 🎯 Great for                       |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------- | ---------------------------------- |
| **Datasheets**                             | Electrical specs, pinouts, power consumption, RF performance, operating temperatures, dimensions, ordering part numbers | Design, selection, troubleshooting |
| **User guides / manuals**                  | Step-by-step setup, features, commands, configuration                                                                   | Getting started, daily use         |
| **Hardware integration guides**            | PCB layout advice, antenna placement, power supply design, grounding, reference schematics                              | Building your own board            |
| **Application notes**                      | Focused how-tos on one topic (low-power tuning, firmware updates, certification tips)                                   | Solving a specific problem         |
| **Quick start guides**                     | Short "unbox and go" walkthroughs                                                                                       | Beginners                          |
| **AT command / software reference guides** | Command lists, parameters, response codes, examples                                                                     | Firmware and software developers   |
| **Development kit documents**              | Board layouts, jumpers, LEDs, connectors, schematics                                                                    | Prototyping                        |
| **Regulatory & certification documents**   | Declarations of conformity, radio approvals, compliance notes                                                           | Getting a product to market        |
| **Antenna documents**                      | Radiation patterns, gain, VSWR, mounting and cable details                                                              | Picking and placing antennas       |
| **Release notes / change notices**         | What changed between firmware or hardware revisions                                                                     | Staying up to date                 |
| **Mechanical drawings**                    | Dimensions, footprints, land patterns                                                                                   | Enclosure and PCB design           |

> 💬 **Tip:** Datasheets are the "source of truth" for a part. When two sources disagree, trust the latest official datasheet.

---

## 🧩 Product families covered

Ezurio organizes products into **series**, each with its own support page. The archive follows those series pages, so you'll find documents grouped by the families Ezurio lists, which broadly fall into categories like:

- 🔵 **Bluetooth / BLE modules** (small, low-power radios for wearables, sensors, beacons and peripherals)
- 🟢 **Wi-Fi modules** (for IoT gateways, appliances and industrial equipment)
- 🟣 **Wi-Fi + Bluetooth combo modules**
- 🟠 **Cellular / LTE-M / NB-IoT connectivity products**
- 🔴 **Long-range radios** (for remote sensing and asset tracking)
- 🟡 **Antennas and RF accessories**
- ⚪ **Development kits, adapters and tools**
- 🛠️ **Legacy product lines** from the Laird Connectivity era

📌 _Product names and families change over time. For the authoritative current list, see the [Ezurio support page](https://www.ezurio.com/support)._

---

## 🗂 Repository layout

```
📦 ezurio-com-documentation
 ┣ 📁 PDFs/            ← 📚 all the documentation lives here
 ┣ 📄 README.md        ← 👋 you are here
 ┗ ⚙️  (tooling used to refresh the archive)
```

Everything you actually care about is inside the **`PDFs/`** folder. 🎉

File names are cleaned up to be **lowercase, readable and unique**, so they're friendly to search tools, scripts and filesystems on every operating system.

---

## 🔎 How to find the document you need

### 🖱️ On GitHub

1. Open the **`PDFs/`** folder.
2. Press **`t`** on your keyboard to use GitHub's _file finder_ and start typing a part number or keyword.
3. Click a PDF to preview or download it.

### 💻 On your computer

Clone the repository, then search however you like:

```bash
git clone https://github.com/PrajwalKoirala638/ezurio-com-documentation.git
cd ezurio-com-documentation/PDFs
```

**List files containing a keyword in the name:**

```bash
ls | grep -i "bluetooth"
```

**Search inside the PDFs (text search):**

```bash
# Linux / macOS (install poppler-utils for pdftotext)
for f in *.pdf; do pdftotext "$f" - 2>/dev/null | grep -il "low power" >/dev/null && echo "$f"; done
```

**Convert everything to plain text for easy searching, notes or AI work:**

```bash
mkdir -p ../text
for f in *.pdf; do pdftotext -layout "$f" "../text/${f%.pdf}.txt"; done
```

### 📱 On your phone or tablet

Download individual PDFs from GitHub and open them in any PDF reader. Many readers support in-document search so you can jump straight to "pinout" or "absolute maximum ratings." 🎯

### 🔤 Search keywords that work well

- Part numbers (e.g., the module name printed on the chip shield)
- `datasheet`, `user guide`, `application note`, `integration`
- `pinout`, `schematic`, `footprint`, `antenna`, `certification`
- `AT commands`, `firmware`, `bootloader`, `dev kit`

---

## 🚀 How these documents benefit you

Technical documents are _free education_ when you know how to use them. Here's the big picture: 🖼️

| 🙋 You are...             | 🎁 What you gain                                                                  |
| ------------------------- | --------------------------------------------------------------------------------- |
| A **student**             | Real-world engineering documents to learn from, far beyond textbook examples      |
| A **hobbyist / maker**    | Confidence to design, wire and program wireless projects correctly the first time |
| A **repair technician**   | Pinouts, voltage levels and block diagrams to diagnose faulty hardware            |
| A **firmware developer**  | Command references and integration details to write correct code                  |
| An **engineer**           | Fast access to specs for part selection and design reviews                        |
| A **teacher**             | Ready-made, professional material for lessons and labs                            |
| A **researcher**          | Primary-source data for comparisons, surveys and reports                          |
| A **community organizer** | Free resources for workshops, repair cafés and maker events                       |

---

## 🎓 Education and learning

Reading real datasheets is one of the **best ways to level up** in electronics and embedded systems. 🧠

### 📚 Things you can learn from these PDFs

- 📶 **How wireless actually works:** radio frequency basics, channels, modulation, range, throughput and interference
- 🔋 **Power management:** sleep modes, current draw, battery life estimation
- 🧮 **Interfaces:** UART, SPI, I²C, USB, GPIO, ADC and how modules communicate with microcontrollers
- 🧱 **PCB design:** land patterns, keep-out zones, ground planes, trace impedance
- 📡 **Antenna fundamentals:** gain, radiation patterns, matching, placement
- 🛡️ **Regulatory basics:** why radios need certification (FCC, CE, ISED, and others) and how it's documented
- 🔐 **Security features:** secure boot, encryption, pairing and authentication concepts
- 🧪 **Testing and validation:** what engineers measure and why

### 🧑‍🏫 Ideas for teachers and self-learners

- ✅ **Datasheet scavenger hunt:** find the supply voltage range, sleep current and operating temperature of a module
- ✅ **Compare two modules:** build a table comparing range, power and price tier to justify a design choice
- ✅ **Design review exercise:** read an integration guide, then critique a sample PCB layout
- ✅ **Vocabulary building:** use the [glossary](#-glossary) and extend it with terms you meet
- ✅ **Estimate battery life:** use current-consumption tables to calculate how long a coin cell would last
- ✅ **Reverse-engineer a block diagram:** redraw a product's architecture in your own words
- ✅ **Write your own quick-start:** turning a dense manual into a one-page guide is a superb way to learn

---

## 🔧 Repairs and troubleshooting

Broken device? Documentation can be the difference between _"throw it away"_ and _"fixed it!"_ ♻️

### 🩺 How the PDFs help diagnose problems

- 🔌 **Check supply voltage requirements** before blaming the module
- 📍 **Use pinout tables** to find the right pad when probing with a multimeter or scope
- 🧭 **Compare expected vs. actual behavior** using documented startup sequences and timings
- 🌡️ **Check operating limits** for temperature and voltage
- 📡 **Inspect antenna details**, since poor range is often an antenna or placement issue, not a dead radio
- 🔁 **Review reset, boot and programming procedures** for devices that won't start
- 🧾 **Identify parts** from ordering information and marking diagrams

### 🛠️ Common issues documentation can help with

| 😕 Symptom                   | 🔎 Where to look                                                    |
| ---------------------------- | ------------------------------------------------------------------- |
| Module won't power up        | Datasheet → supply voltage, power-up sequence                       |
| No communication over serial | User guide → UART settings, baud rate, flow control                 |
| Weak signal / short range    | Integration guide → antenna placement, ground plane, keep-out areas |
| Random resets                | Datasheet → current peaks, decoupling recommendations               |
| Can't pair or connect        | Software guide → pairing, advertising, connection parameters        |
| Firmware update failed       | Application note → bootloader and update procedure                  |
| Overheating                  | Datasheet → absolute maximum ratings, thermal notes                 |

> ⚠️ **Safety first:** always disconnect power and batteries before working on hardware, respect battery safety, and never exceed documented voltage or current limits.

### 🌍 Why this matters for sustainability

The **right to repair** and **extending device life** reduce electronic waste. Accessible documentation helps technicians and owners fix hardware instead of discarding it. 🌱

---

## 🧰 DIY and maker projects

Want to build something cool? These PDFs are your blueprint. 🏗️

### 💡 Project ideas inspired by wireless modules

- 🌡️ **Wireless temperature and humidity sensor** that reports to your phone over BLE
- 🚪 **Door/window sensor** with months of battery life
- 🐾 **Pet or asset tracker** using low-power radios
- 🌱 **Smart garden monitor** for soil moisture and watering reminders
- 🏠 **Home automation nodes** that join your Wi-Fi network
- 🎛️ **Bluetooth remote control** for lights, music or cameras
- 📟 **BLE beacon** for indoor navigation or museum guides
- 🏃 **Wearable prototype** for fitness or health sensing
- 🔔 **Wireless doorbell or alert button**
- 📊 **Data logger** that syncs to a phone or gateway

### 🧭 A simple DIY workflow

1. 🎯 **Define** what you want to build and your range, battery and size needs.
2. 📖 **Pick a module** by comparing datasheets.
3. 🔌 **Read the hardware integration guide** before drawing any schematic.
4. 🧪 **Prototype** on a development kit.
5. 🧱 **Design your PCB** using the documented footprint and layout rules.
6. 💻 **Program** using the software/command reference.
7. 📡 **Test range and battery life** and tune.
8. 📝 **Document your project** and share it with the community. 🎉

### 🧠 Maker pro tips

- Always start with the **dev kit documentation**; it's the gentlest on-ramp.
- Follow **antenna keep-out guidance** exactly. It's the #1 cause of poor wireless performance.
- Add **decoupling capacitors** where the datasheet recommends.
- Keep a **project notebook** and save the PDF version you used so your build stays reproducible.

---

## 👷 Professional engineering

For working engineers, this archive speeds up the daily grind. ⏱️

- ⚖️ **Part selection:** compare specs across product families in minutes
- 🧾 **Design reviews:** have the exact reference document on hand
- 🧱 **Hardware integration:** consult layout and RF guidance during PCB design
- ✅ **Certification planning:** review regulatory documents early to avoid surprises
- 📦 **Bill of materials checks:** confirm ordering codes and variants
- 🔄 **Migration planning:** compare older and newer revisions using release notes
- 🧑‍🤝‍🧑 **Onboarding:** hand new teammates one folder instead of a dozen bookmarks

---

## 🤖 AI and machine learning uses

Technical PDFs are a rich source of **structured, factual, domain-specific text**. People use collections like this for many productive purposes. 🧠✨

### 🧪 Ideas

- 💬 **Question-answering assistants** that cite specific datasheets (retrieval-augmented generation, RAG)
- 🔍 **Semantic search** across thousands of pages of technical content
- 🏷️ **Information extraction:** pulling pinouts, voltages, temperatures and part numbers into structured tables or JSON
- 🗃️ **Knowledge graphs** linking products, interfaces, features and document types
- 📊 **Comparison tools** that rank modules by range, power use or interface support
- 🧑‍🏫 **Tutoring and study aids:** quizzes, flashcards and explanations generated from real specs
- 🌐 **Translation and simplification** of dense technical language
- 🔧 **Repair-assistant prototypes** that help diagnose symptoms using documented behavior
- 🧪 **Benchmarks and evaluations** for table, figure and layout understanding in documents

### 🛠️ A simple pipeline

```
📥 PDFs  →  🔤 text extraction / OCR  →  ✂️ chunking by section
        →  🧮 embeddings  →  🗄️ vector index  →  💬 question answering with citations
```

### ✅ Best practices

- 📎 **Always cite the source document** and its version/date in answers
- 🧾 **Keep the original PDF accessible** so humans can verify AI output
- 🔬 **Validate extracted numbers** (voltages, currents, pin numbers). Errors here can damage hardware
- 🧱 **Preserve tables and figures** where possible; many key specs live in tables
- 🕒 **Track versions**, since specs change between revisions
- ⚖️ **Respect licensing and terms of use** (see the [legal notice](#%EF%B8%8F-legal-notice-and-licensing)). Check whether your intended use, especially model training or redistribution, is permitted before you proceed

> 🛑 **Important:** An AI summary is a helpful _starting point_, never a replacement for the official datasheet when designing or repairing real hardware.

---

## 🏘 Benefits for your community

Knowledge grows when it's shared. 🌳

- 🔧 **Repair cafés & fix-it events:** volunteers can look up specs on the spot
- 🏫 **Schools & after-school clubs:** real-world material for STEM lessons
- 🧑‍💻 **Makerspaces & hackerspaces:** a shared reference library for members
- 🎓 **University labs:** supplementary reading for embedded systems and wireless courses
- 🚀 **Startups & small businesses:** lower the barrier to prototyping
- 🌍 **Regions with limited internet:** offline copies make the information accessible anywhere
- 🧓 **Lifelong learners:** a free, self-paced path into electronics
- 🤝 **Open-source projects:** a handy reference for drivers, firmware and hardware designs

### 💬 Ways to spread the benefit

- Host a **"read a datasheet" workshop** 🧑‍🏫
- Build a **local mirror** for a classroom or lab 🖥️
- Create **beginner guides** from the dense originals ✍️
- Share **your own project write-ups** that cite these documents 📝
- Help others **find the right document** on forums and chat servers 💬

---

## 📖 Reading a datasheet: a quick guide

Datasheets can be intimidating. Here's a friendly roadmap. 🗺️

### 1️⃣ Start with the overview

Look at the **features list** and **block diagram** to understand what the part does.

### 2️⃣ Check the absolute maximum ratings ⚠️

These are limits you must **never exceed**. They are not normal operating conditions.

### 3️⃣ Read the recommended operating conditions ✅

This is where the supply voltage and temperature range you should actually use are listed.

### 4️⃣ Study the pinout 📍

Match **pin names and numbers** to functions such as power, ground, serial lines and GPIO.

### 5️⃣ Look at electrical characteristics ⚡

Find current consumption in each mode (active, idle, sleep) to estimate battery life.

### 6️⃣ Review RF performance 📡

Look for transmit power, receiver sensitivity, supported frequencies and expected range.

### 7️⃣ Follow layout and antenna guidance 🧱

Respect **keep-out areas** and **ground plane** advice.

### 8️⃣ Don't skip ordering information 🧾

Variants differ (antenna type, memory, temperature grade), so confirm the right part number.

### 9️⃣ Check the revision history 🕒

Know which version of the document you're reading.

---

## 🔤 Glossary

| Term                         | Meaning                                                                   |
| ---------------------------- | ------------------------------------------------------------------------- |
| **BLE**                      | Bluetooth Low Energy, a power-sipping flavor of Bluetooth                 |
| **RF**                       | Radio Frequency                                                           |
| **dBm**                      | A unit of radio power. Higher means stronger signal                       |
| **RSSI**                     | Received Signal Strength Indicator, a measure of how strong a signal is   |
| **Antenna gain**             | How well an antenna focuses energy in a direction                         |
| **VSWR**                     | Voltage Standing Wave Ratio, a measure of antenna matching quality        |
| **GPIO**                     | General-Purpose Input/Output pin                                          |
| **UART**                     | Serial interface common for module communication                          |
| **SPI / I²C**                | Common chip-to-chip communication buses                                   |
| **AT commands**              | Text commands used to control a module over serial                        |
| **Bootloader**               | Small program that starts the device and can update firmware              |
| **OTA**                      | Over-the-Air (firmware update performed wirelessly)                       |
| **Keep-out area**            | Region of a PCB that must be free of copper or components near an antenna |
| **Ground plane**             | Large copper area that acts as a reference for circuits and antennas      |
| **Land pattern / footprint** | The copper pads layout on a PCB where a part is soldered                  |
| **Decoupling capacitor**     | Small capacitor that stabilizes the power supply near a chip              |
| **Certification**            | Official approval that a radio device meets regulations                   |
| **FCC / CE / ISED**          | Regulatory bodies/marks for the US, Europe and Canada                     |
| **SoC**                      | System on Chip, a single chip containing processor and radio              |
| **IoT**                      | Internet of Things, everyday devices connected to networks                |
| **LPWAN**                    | Low-Power Wide-Area Network, long-range, low-data-rate wireless           |
| **PCB**                      | Printed Circuit Board                                                     |
| **BOM**                      | Bill of Materials, the parts list for a product                           |

---

## 🔄 Keeping the archive fresh

Documentation gets updated. 🆕 This archive is refreshed by re-collecting the PDFs from Ezurio's public pages, so newer revisions and newly added products can be picked up over time.

**What you can do:**

- 📅 Check the **last updated** date (maintainers: add it near the top!)
- 🔗 For critical designs, **verify against the official [Ezurio support site](https://www.ezurio.com/support)**
- 🐛 Spot an outdated or missing document? **Open an issue.**

---

## 🤝 Contributing

Contributions are warmly welcome! 💖

- 🐞 **Report problems:** broken or unreadable PDFs, misnamed files, missing products
- 🗂️ **Improve organization:** suggest better folder structures or indexes
- 📝 **Write guides:** beginner-friendly summaries or tutorials that point back to the official documents
- 🔤 **Extend the glossary** and FAQ
- 🌐 **Translations** of guides and explanations
- 🧠 **Share learning materials:** lesson plans, workshop outlines, project write-ups

**How to contribute:**

1. 🍴 Fork the repository
2. 🌿 Create a branch (`git checkout -b my-improvement`)
3. ✍️ Make your changes
4. 📬 Open a pull request with a clear description

Please be kind, respectful and constructive. We're all here to learn. 🌈

---

## ⚖️ Legal notice and licensing

> 📢 **Please read this section carefully.**

- 🏷️ **All documents in the `PDFs/` folder are the property of their respective copyright holders (Ezurio, formerly Laird Connectivity, and any third parties named inside the documents).** All trademarks, logos and product names belong to their owners. The Bluetooth® word mark is a registered trademark of Bluetooth SIG, Inc.
- 🚫 This repository is **not affiliated with, endorsed by, or sponsored by Ezurio.**
- 📜 The PDFs are provided here for **convenience, reference and educational purposes.** Their use and redistribution remain subject to the terms set by the copyright holders and the Ezurio website.
- 🔐 **Some Ezurio documents may be intended to be accessed only through their website or after signing in.** If you are a rights holder and want a file removed, or you believe something here should not be shared, please open an issue or contact the maintainer, and it will be taken down promptly.
- 🤖 **Using documents for AI/ML training, commercial products or redistribution may require permission** from the copyright holder. Please review the applicable terms and obtain permission where needed.
- ⚠️ **No warranty.** Information may be outdated, incomplete or altered. **Always verify against official sources before designing, manufacturing, certifying or repairing any product.** The maintainers are not responsible for damage, injury or loss resulting from the use of this information.
- 🛡️ Follow **local laws and radio regulations** in your country when building or modifying wireless devices.

---

## ❓ FAQ

**🙋 Do I need to be an engineer to use this?**
Not at all! Beginners can start with quick-start guides and user guides, then work up to datasheets.

**🙋 Are these documents free?**
They are gathered from publicly reachable pages of the Ezurio website, but copyright still belongs to their owners. See the [legal notice](#%EF%B8%8F-legal-notice-and-licensing).

**🙋 Can I use these for my school project or thesis?**
Absolutely. Cite the document title, publisher and revision/date, and link to the official source where possible. 🎓

**🙋 Can I use them to train an AI model?**
Possibly, but check the copyright holder's terms first. For personal learning, search and retrieval experiments, cite your sources and keep the originals accessible. For redistribution or commercial training, get permission.

**🙋 Is the information current?**
It's current as of the last refresh. For any real design, confirm against the official website. 🔎

**🙋 What if a PDF won't open?**
Try another PDF reader or re-download it. If it's still broken, open an issue. 🛠️

**🙋 Why does the Laird name appear in some documents?**
Ezurio is the successor brand of Laird Connectivity, so older documents may carry the earlier name.

**🙋 Can I mirror this archive for my classroom or makerspace?**
For personal and educational use, cloning is simple. For public mirroring or redistribution, see the legal notice first.

---

## 💖 Final words

Great technology is built on shared knowledge. 🌐 By keeping documentation accessible, we help students learn faster, makers build smarter, technicians repair more devices, and communities grow stronger.

⭐ **If this archive helped you, please star the repo and share it with someone who'd love it!**

**Happy learning, building and fixing!** 🔧📡🚀

---

<sub>📌 Disclaimer: Independent community archive. Not affiliated with Ezurio or Laird Connectivity. All trademarks and copyrights belong to their respective owners.</sub>
