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

# Git Work Methodology
## General Idea
La siguiente imagen describe la metodología de trabajo que se va a seguir con Git:
![image explaining the git work methodology](https://wac-cdn.atlassian.com/dam/jcr:cc0b526e-adb7-4d45-874e-9bcea9898b4a/04%20Hotfix%20branches.svg?cdnVersion=1090)
En esta página se describe en mayor detalle la imagen [Gitflow Workflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow)

La idea principal es que el `main` solo se deje para versiones estables y que se trabaje en `default`. La rama `default` esta puesta como default en Github, cada vez que se quiera añadir una funcionalidad o arreglar algo lo que se hará es crear una rama nueva a partir de develop (`feature`). Posteriormente se mergeará esa rama con develop de nuevo. Cuando haya suficientes funcionalidades añadidas a develop y queramos sacar una release haríamos lo siguiente:
- A partir de develop (código que queremos utilizar para el release), creamos una rama `release`, en esta rama se testeará la aplicación y se corregirán pequeños cambios para preparar la release.
- Una vez que se haya visto que la herramienta está lista al público:
   - Se añaden esas correcciones a develop mergeando release a develop
   - Se mergea la rama release al main

Adicionalmente puede existir la rama Hotfix, también es una rama auxiliar, lo que permite es hacer correcciones que sean totalmente inmediatas. Por ejemplo, se encuentra un bug en la aplicación, o una librería vulnerable. En ese caso, a partir del main se crea ese hotfix-x, se cambia la librería para actualizarla a una versión no vulnerable, luego:
- Se hace merge con el main para sacar una nueva release que resuelve ese problema, por ejemplo de v1.0.1 a v1.0.2
- Se hace merge con develop para que esa corrección también esté en develop a partir de ahora
## Using Git and Github to implement the methodology
### Creando lo necesario en Github
Imaginemos que queremos añadir una funcionalidad nueva a nuestra herramienta. Lo primero que haríamos sería crear un issue que indicaría lo que queremos mejorar o solucionar. Tras crear el issue, a la derecha debemos de:
- Assignees. Escoger quien va a añadir esa funcionalidad
- Labels. Escoger una etiqueta que clasifique el issue. Por ejemplo "enhancement" para una nueva funcionalidad, o "documentation" para mejoras en la documentación
- Projects. Escogemos un proyecto. Posteriormente cuando el issue se haya creado podremos escoger el estado del issue. Dentro de la tabla hay descripciones de para que sirve cada tablita.
   En el caso en el que vayamos a trabajar en esta mejora ahora mismo, escogeremos el status "In progress". Además, pulsaremos en "Development->Create Branch". Como nombre podemos utilizar algo como feature-n, donde "n" es el número del issue, de esta forma siempre sabemos que la rama feature-n está asociada con el issue n, podemos dirigirnos al issue y ver rápidamente que se pretende hacer o se ha hecho leyendo la descripción.

### Desarrollando la funcionalidad
En este punto el issue está creado, está dentro del proyecto con el estado "In progress", y tenemos una rama feature-n creada en Github. Ahora en local lo que haremos será ponernos a trabajar en esa rama, para ello ejecutamos:
```
git fetch origin
git checkout feature-n
```

Ya estamos listos para hacer las modificaciones en el código. Hacemos varias modificaciones con nuestro editor favorito y nos volvemos a ir a git. Para ver los archivos que se han modificado se puede hacer `git status`.

Ahora podemos ver las modificaciones de cada uno de esos archivos usando `git diff <filename>`.

O ver las modificaciones de todos los archivos si no se añade <filename>.
Lo ideal es ir revisando archivos uno a uno, asegurar que las modificaciones son correctas y que no se nos ha olvidado nada (como comentarios obscenos) y añadir ese archivo con `git add <filename>`.
Cuando añadimos un archivo con `git add`, al hacer `git status`, podemos ver que ahora ese archivo está en verde. Seguimos añadiendo archivos y comprobando que todo está bien o usamos `git add .`  para añadir todos a la vez.
Una vez que todos los archivos están añadidos podemos hacer commit (registrar esta versión del código, lo que nos permite volver si jodemos algo y ver que se ha ido modificando). Para ello hacemos `git commit`, se nos abrirá un editor en la terminal, abajo del todo añadimos una breve descripción de la siguiente forma:
```
# Comentarios indicando que se va a añadir
#
#
#
Esta primera linea es una descripción muy breve

Dejando un salto de linea, aquí podemos poner una descripción más compleja cuando sea necesario.
```

Guardamos y salimos del editor. Aquí ya está guardado esa versión del código. Puedes ver que el commit se ha hecho con `git log`, con este comando puedes ver el historial de todos los commmits anteriores. Ya está todo listo para hacer otras modificaciones, o si eso era todo, es hora de subir el código a github. Para ello, simplemente haz `git push origin feature-n`.

### Haciendo merge de una funcionalidad nueva
Dentro de Github crearemos una nueva "Pull Request". Como base usaremos "develop" y como compare la rama que queremos mergear en "develop". En la descripción simplemente escribe:
"Related with issue #n" donde n es el número del issue. Y añade información adicional a tener en cuenta si es necesario. Una vez creada la Pull Request escoger a un "Reviewer". Una persona que quieres que revise el código antes de añadirlo a develop. Esto es importante porque 4 ojos ven más que 2, y nunca viene mal echar un vistazo por encima para:
- Estar al día de las cosas que se añaden al código. Por esto también es muy bueno añadir comentarios
- Revisar porque siempre se pasan cositas que se pueden mejorar o crear un bug.

Hasta este punto el trabajo de la persona que ha creado la funcionalidad ya está hecho. Ahora otra persona revisará el Pull Request y lo aceptará y/o mergeará en develop. Lo único que queda es, una vez que se haya mergeado feature-n en develop, ir a local y borrar esa rama porque ya no vamos a trabajar en ella. En git en local ejecutamos:
```
git checkout develop
git branch -d feature-n
```

# Building
- Windows event logs
- Other: MUICache, (LastVisitedMRU, LastVisitedPidMRU, OpenSaveMRU), BHO.
