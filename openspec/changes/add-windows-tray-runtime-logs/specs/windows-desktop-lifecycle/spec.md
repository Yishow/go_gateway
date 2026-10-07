## ADDED Requirements

### Requirement: Windows desktop delivery does not create a console

The Windows desktop release SHALL be one GUI-subsystem executable containing its backend, complete embedded frontend asset graph, and tray icon. Normal launch and browser-opening actions MUST NOT create or flash a console window. Generated configuration data and bounded diagnostic files SHALL NOT require an additional installed runtime or companion executable. Development console builds and Linux/macOS CLI builds SHALL remain available without desktop dependencies.

#### Scenario: Launch the standalone desktop artifact

- **WHEN** an operator launches the Windows desktop release from Explorer on a supported interactive Windows system without Node.js or Go installed
- **THEN** the gateway creates its tray entry without a console or a separate application taskbar window
- **AND** the embedded setup and log routes load from the executable's own assets after readiness
- **AND** opening either route does not launch a command-shell window

#### Scenario: Preserve console and headless execution

- **WHEN** a console development build, Linux/macOS CLI build, or Windows release with explicit `--headless` is started
- **THEN** no tray or native interactive dialog is required
- **AND** existing CLI configuration precedence, exit reporting, and redirected output remain usable
- **AND** every build in explicit `--headless` mode suppresses automatic browser opening even when AUTO_OPEN_BROWSER is true
- **AND** legacy console invocation without that explicit flag retains its opt-in AUTO_OPEN_BROWSER behavior

#### Scenario: Capture a Windows GUI artifact in headless mode

- **WHEN** the GUI-subsystem release runs with `--headless` under cmd, PowerShell, and a supervisor harness that supplies redirected standard handles
- **THEN** tests verify actual output capture, process wait behavior, and exit status in each launch method
- **AND** absent console handles do not cause a new console allocation or a false claim of captured diagnostics

#### Scenario: Windows-only release flags

- **WHEN** the documented Windows desktop and console build commands and a non-Windows build command run
- **THEN** only the Windows desktop artifact has the GUI subsystem
- **AND** all embedded release paths fail on missing frontend assets instead of reporting a complete UI delivery
- **AND** a filename ending in `.exe` is not used as evidence of the target operating system

### Requirement: Desktop data selection is explicit and non-destructive

The desktop launcher SHALL establish a canonical data source before migration, listeners, or acquisition. Explicit process configuration SHALL take precedence over executable-directory dotenv and defaults; existing database aliases SHALL retain `GATEWAY_DB_PATH`, then `DB_PATH`, then `SQLITE_PATH` precedence. Desktop-relative settings SHALL use the documented executable/data-root bases, while CLI working-directory semantics SHALL remain unchanged. The system MUST NOT silently migrate, overwrite, select an ambiguous legacy database, or create a replacement database after a path/permission failure.

#### Scenario: An explicit existing database wins

- **WHEN** a desktop launch supplies an absolute existing database through the highest-priority supported variable
- **THEN** that database and its canonical identity determine the owner guard and data root
- **AND** lower-priority aliases and the executable-directory default do not create another database

#### Scenario: A legacy working directory is ambiguous

- **WHEN** the launch working directory differs from the executable directory and contains a competing `.env` or `datalink.db`, without explicit unambiguous configuration
- **THEN** startup stops before migration or acquisition with a native explanation of the competing locations
- **AND** the operator is instructed to configure the intended absolute database path
- **AND** neither database is moved, overwritten, or newly created

#### Scenario: Confirm a genuinely new desktop database

- **WHEN** no explicit database is configured and the resolved default database does not exist
- **THEN** the native launcher asks before creating a new database at the displayed location
- **AND** cancel leaves storage and acquisition unchanged
- **AND** accepting does not imply that every historical installation was discovered or migrated

#### Scenario: An unwritable root fails visibly

- **WHEN** the selected database or initial diagnostic directory cannot be used
- **THEN** desktop startup reports the failure through a native dialog and exits nonzero
- **AND** it does not elevate permissions or fall back to another database location

### Requirement: One Windows owner controls each canonical database

Windows desktop and headless modes SHALL acquire the same OS-backed exclusive ownership guard for a canonical database before initialization. Lock metadata SHALL NOT replace live ownership verification, and stale files or reused PIDs MUST NOT trigger process termination. A secondary desktop instance SHALL only request a fixed open-setup action from an authenticated same-user, same-session owner over access-controlled local IPC; it MUST NOT accept arbitrary commands or URLs.

#### Scenario: Concurrent launch cannot duplicate acquisition

- **WHEN** two Windows processes start against the same database through equivalent path spellings
- **THEN** at most one process opens/migrates the database and starts services
- **AND** the other opens the verified existing setup page or reports the already-running owner and exits
- **AND** no second collector, delivery worker, or Share listener is created

#### Scenario: A stale owner marker does not block recovery

- **WHEN** a previous process has terminated but diagnostic PID metadata remains
- **THEN** OS ownership determines whether a new process can acquire the database
- **AND** no process is killed solely because a recorded PID exists

#### Scenario: Handoff is unavailable or untrusted

- **WHEN** the owner is headless, in another session, or fails the local IPC identity check
- **THEN** the secondary process reports that the gateway is already running without launching another runtime
- **AND** it does not send an arbitrary command, bypass the owner guard, or open an unverified service

### Requirement: Readiness and startup failures remain observable

The desktop launcher SHALL initialize safe diagnostics and the tray, reserve the configured HTTP listener, initialize services, and serve the embedded UI before announcing web readiness. Runtime availability SHALL be represented separately from HTTP readiness. Initial tray failure or fatal startup failure MUST produce a native actionable error with a stable safe code and release acquired resources. The diagnostic file path SHALL be shown only when that file was actually written; inability to save diagnostics SHALL be explicit.

#### Scenario: A different application owns the configured port

- **WHEN** listener binding fails because the configured port is occupied
- **THEN** the gateway reports a port-conflict startup failure and exits nonzero after cleanup
- **AND** it does not kill the owner, change ports, start acquisition, or open the other application's page

#### Scenario: HTTP serves while acquisition is unavailable

- **WHEN** the HTTP listener and assets are available but runtime startup fails under the existing degraded-start policy
- **THEN** the tray presents web availability and degraded acquisition separately
- **AND** it never uses a healthy collector or SQL-committed indication solely because HTTP is running

#### Scenario: Tray creation fails before service startup

- **WHEN** the initial Windows tray cannot be created
- **THEN** a native startup error appears and the process cleans up and exits nonzero
- **AND** no invisible background acquisition instance remains

#### Scenario: An early fatal occurs without a console

- **WHEN** configuration, storage, migration, or service initialization fails before the web page is available
- **THEN** the operator receives a safe native error and next step
- **AND** a successfully saved diagnostic is identified, or the dialog explicitly states that saving failed
- **AND** raw credentials, DSNs, configuration dumps, and backend exception details are not displayed

### Requirement: The tray exposes a small truthful local menu

The tray SHALL provide Open setup, View runtime logs, read-only process and acquisition status, Version information, and Exit with confirmation. Setup SHALL target `/studio/v2`; logs SHALL target `/studio/logs`. URLs MUST be constructed from a verified listener and fixed product paths without credentials. Wildcard binding SHALL resolve to a usable loopback URL. A specific non-loopback-only bind SHALL NOT weaken the log access boundary. State text SHALL distinguish starting, serving, degraded, stopping, and stop-timeout. No restart, startup-registration, service-installation, or remote-control menu item SHALL be introduced.

#### Scenario: An operator uses the menu

- **WHEN** the ready operator opens the tray menu with mouse or keyboard
- **THEN** setup and log actions open their fixed routes and status/version actions remain read-only
- **AND** version metadata is actual product/build/commit information or explicitly unknown/dev
- **AND** setup is disabled before readiness and logs explain any unavailable local-only access

#### Scenario: The browser is closed or minimized

- **WHEN** the operator closes or minimizes any browser window opened from the tray
- **THEN** acquisition, HTTP serving, and logging continue without a lifecycle transition
- **AND** the tray can reopen the page

#### Scenario: Explorer rebuilds the notification area

- **WHEN** Windows Explorer recreates the taskbar after successful startup
- **THEN** the application re-registers the tray entry with current truthful state
- **AND** it does not restart the gateway or duplicate workers
- **AND** icon visibility in the Windows overflow area is not misrepresented as a failed process

### Requirement: Exit is confirmed, coordinated, and honest about delivery

User-initiated tray exit SHALL require confirmation that acquisition and serving will stop. Cancel SHALL leave the running state unchanged. Confirmed exit, supported process signals, and OS-session shutdown SHALL share an idempotent coordinator that stops new work, ends log streams, settles in-flight work, stops runtime intake before group delivery, closes Share/connections/storage, then finalizes diagnostics and ownership. The system MUST preserve existing durable accepted-data and unknown-outcome semantics and MUST NOT claim that all backlog was delivered merely because shutdown was requested.

#### Scenario: Cancel then confirm exit

- **WHEN** an operator first cancels Exit and later confirms it
- **THEN** cancellation leaves all services running
- **AND** confirmation transitions once to stopping, disables repeated exit actions, and invokes each cleanup phase at most once
- **AND** successful completion removes the tray and owner guard after resource cleanup

#### Scenario: The shutdown deadline expires

- **WHEN** the shared 15-second shutdown notification deadline expires before cleanup is verified
- **THEN** the tray/native status reports stop-timeout and the pending phase without claiming stopped or delivered
- **AND** it keeps observing the original coordinator without closing storage beneath active workers
- **AND** a desktop operator can continue waiting or separately confirm force termination with an explicit incomplete/unknown-data warning
- **AND** headless reporting remains non-successful and does not create a dialog or falsely report a graceful exit

#### Scenario: Durable backlog survives exit

- **WHEN** accepted group data remains undelivered or an external commit outcome remains unknown at shutdown
- **THEN** shutdown retains the journal/outbox/receipt and existing recovery rules
- **AND** a later start neither silently clears backlog nor treats unknown as committed
- **AND** the UI does not promise recovery of legacy data that was never durably accepted

#### Scenario: Multiple shutdown causes race

- **WHEN** a tray confirmation and a supported termination signal arrive concurrently while log SSE clients are connected
- **THEN** one coordinator cancels streams and services without double-close panics or duplicate cleanup
- **AND** no indefinitely live log subscriber alone prevents shutdown progress
