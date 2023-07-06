# gdback-priv

- Hash algorithm = SHA-512
- $unixTime := strconv.FormatInt(time.Now().Unix(), 10) + ".sqlite"
## Syntax Golang Name Packages
### Examples
- gxdb: g->gdback, x->multiOS, db->name package
- gwfiles: g->gdback, w->windows, files->name package
## Functions for data extraction
### RAM Memory Dump File ram_$unixTime.mem (Physical memory) [Belkasoft Live RAM]
### Ram Page Dump File pagefile_$unixTime.sys (Virtual memory) [Belkasoft Live RAM]
### File system_$unixTime.sqlite
- **general** table: Manufacturer, UUID, Processor, RAM, Hostname, Storage, Operating system, MAC address, IP address, System time, System date.
- **devices** table: Device description notified by bus, Type of device, Description of device, Hardware ID.
- **users** table: [PSLoggedon.exe]
- **process** table: [pslist.exe] Name, PID, PRI, THD, HND, PRIV, CPU Time, Elapsed Time.
- **services** table: [PSService64.exe]
### File network_$unixTime.sqlite
- **status** table: [ipconfig /all]
- **established** table: [nbtstat -S, net sessions]
- **active** table: [netstat -an] Protocol, Local address, Foreign address, State.
- **dns** table: [ipconfig /displaydns]
- **arp** table: [arp -a]
### File wireless_$unixTime.sqlite
- **wifi** table: SSID, BSSID, Type, Authentication, Encryption.
- **bluetooth** table:
### File configuration_$unixTime.sqlite
- **firewall** table:
- **security** table: Configuration of Windows Security/Action Center.
### File history_$unixTime.sqlite
- **cmd** table: Command, Time.
- **powershell** table: Command, Time.
- **downloads** table: Filename, Time.
- **clipboard** table: String, Time.
- **internet** tables: Search, Time.
### File files_$unixTime.sqlite
- **general** table: Path, Filename, File Extension, Hash SHA512, Size Bytes, Creation Date, Last Modification Date.
- **shared** table: Shared Folders.
### File hashes_$unixTime.sqlite
- **general** table: Filename, SHA-512.


## Artifacts to extract from Windows

### From Wintriage

So far, it gets the next artifacts, if selected (following forensics best practices based in volatility order):

- Memory dump
- Prefetch
- Commands execution (some of them are native operating system ones and other external tools calls) to get a bunch of live information, network capture, etc…
- Info about encrypted volumes: bitlocker recovery key, truecrypt and veracrypt warnings if mounted
- Alternate Data Streams in every mounted volume
- Windows Registry
- Events: EVT/EVTX and ETL
- SRUM
- Failed spool Jobs
- All users Trash bin in every logical volume
- Active Directory ntds.dit (if executed in a Domain Controller)
- Users artifacts: Registry, shellbags, browsers (IE, Edge, Edge Chromium, Chrome, Firefox, Brave and Opera), Jumplist, Recent, Office Recent, etc…
- Shadow copies. The same artifacts in every shadow copy the system has.
- Live forensic image. It does a live image of C:\ Logical drive in EWF format. This is useful in case the volume is encrypted, so the physical image could be more complicated to be processed later.

### Asking ChatGPT for other artifacts

Remember, the relevance and availability of these artifacts may vary depending on the specific system configuration, the activities being investigated, and the version of Windows used.

- Event Logs: These logs provide information about various system events, including security-related events, application crashes, system startup/shutdown, user logins/logouts, etc. They are stored in files with the extension .evtx.
- Internet History and Cookies: Extracting web browser history and cookies can provide insights into the user's online activities. The locations of these artifacts vary based on the browser used (e.g., Internet Explorer, Edge, Chrome, Firefox).
- User Account Information: Gather information about user accounts, including usernames, passwords (if stored), user profiles, and security identifiers (SIDs).
- System Logs: Extract system logs such as system status logs, application logs, setup logs, and security logs. These logs can provide valuable information about system events and potential security incidents.
- File Metadata: Retrieve file metadata such as file creation/modification timestamps, file sizes, and file permissions. This information can help establish a timeline of events and track file access.
- Printer Spooler Artifacts: Investigate the printer spooler for any evidence of print jobs, including failed or successful print jobs, printer configurations, and related timestamps.
- System Configuration: Extract information about system configuration settings, including installed software, hardware devices, system settings, and driver information.
- USB Device History: Look for traces of USB device connections, including connected USB devices, timestamps, and device metadata.
- System Restore Points: Explore system restore points to gather information from previous system states, including file versions, system settings, and configurations.
- Windows Logs: Extract information from Windows-specific logs, such as the SetupAPI log, which can provide details about the installation and setup of hardware and software components.

### Asking ChatGPT for open-source tools to extract those artifacts

Certainly! If you prefer open-source tools for your forensic investigation, here are some popular ones that can assist in extracting artifacts from a Windows machine:

- Autopsy: Autopsy is an open-source digital forensic platform that provides a range of features for artifact extraction and analysis. It offers support for Windows artifacts, including file system metadata, registry entries, event logs, internet history, and more.
- The Sleuth Kit (TSK): The Sleuth Kit is an open-source forensic toolkit that supports the analysis of various file systems. It includes tools such as "fls" for file system timeline analysis and "ils" for listing allocated and deleted files.
- Volatility: Volatility is an open-source memory forensics framework that allows the extraction and analysis of volatile memory (RAM) artifacts. It provides plugins specifically designed for analyzing Windows memory dumps, including processes, network connections, registry keys, and more.
- Plaso (Log2Timeline): Plaso, also known as Log2Timeline, is an open-source tool for creating timelines and analyzing various log files. It supports parsing Windows event logs (EVT/EVTX), as well as other log formats, to create a cohesive timeline of events.
- RegRipper: RegRipper is an open-source Windows Registry analysis tool that helps in extracting information from Registry hives. It provides a set of plugins to parse and extract relevant artifacts from different Registry keys related to user accounts, installed software, network settings, and more.
- Wireshark: Wireshark is a widely-used open-source network protocol analyzer that allows the capture and analysis of network traffic. It can assist in extracting network artifacts, identifying communication patterns, and detecting potential malicious activities.
- Bulk Extractor: Bulk Extractor is an open-source command-line tool designed for extracting various artifacts from disk images, including email addresses, URLs, credit card numbers, and other potentially sensitive information. It can assist in finding hidden information and patterns within large data sets.


# Building
- Windows event logs
- Other: MUICache, (LastVisitedMRU, LastVisitedPidMRU, OpenSaveMRU), BHO.
