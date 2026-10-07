## ADDED Requirements

### Requirement: Managed runtime logs are safe application events

The gateway SHALL provide one managed event representation for startup, runtime, HTTP, recoverable application faults, and shutdown, with opaque instance ID, monotonic string sequence, UTC timestamp, level, source, stable code, safe message, allowlisted scalar fields, and truncation metadata. Sanitization MUST occur before any newly introduced memory, file, or web sink stores an event. The system MUST NOT tee raw process stdout/stderr or treat arbitrary error strings as safe diagnostics. Existing engineering `/debug/logs`, packet history, audit events, and delivery truth SHALL remain distinct.

#### Scenario: Production sources feed the managed view

- **WHEN** known startup/runtime, standard log, slog, Gin access/recovery, and HTTP-server events occur
- **THEN** the managed stream receives classified safe events or explicit fixed suppression events for unclassified raw content
- **AND** required startup/runtime faults include actionable stable codes rather than only a generic success message
- **AND** the implementation records its source-coverage inventory without claiming capture of every OS/runtime crash or third-party byte

#### Scenario: Secret-bearing input never enters a new sink

- **WHEN** a source includes credentials, cookies, authorization headers, DSNs, full configuration, raw SQL/arguments, packet bodies, endpoints, private paths, or nested arbitrary error text
- **THEN** none of those fixture values appear in the managed ring, queues, files, snapshot, SSE, or operator DOM
- **AND** allowlisted safe templates and fields are used before storage, with regex redaction only as defense in depth

#### Scenario: Safe HTTP logs do not inherit full query logging

- **WHEN** a request has a secret query value or raw backend error
- **THEN** the managed HTTP event contains only method, route template, status, duration, safe code, and opaque request ID as applicable
- **AND** it contains neither query values nor raw headers/body/error text
- **AND** existing console formatter compatibility is tested separately from the new safe projection

#### Scenario: Log reads cannot amplify themselves

- **WHEN** clients read snapshots, consume streams, or receive heartbeats
- **THEN** those actions do not generate an ordinary access-log event per read or heartbeat
- **AND** sink failure reporting does not recursively write back into the failing sink

### Requirement: Retention and fan-out have independent hard bounds

The managed log system SHALL enforce the design limits below, with the first reached count or byte limit taking effect. These limits are design requirements, not measured performance claims. Producers MUST NOT wait for disk or network I/O; admission drops, retention eviction, subscriber overflow, and file failures SHALL be separately observable. Every queue, index, counter representation, and browser buffer MUST remain bounded.

- Encoded event: 8 KiB, safe UTF-8 truncation with a flag
- Admission queue: 1,024 records and 2 MiB
- Current-process ring: 2,000 records and 4 MiB, oldest-first eviction
- Concurrent streams: 8; per-stream queue: 256 records and 512 KiB
- Separate per-stream replay copy: 500 records and 1 MiB, with explicit replay-limit gap/tail on overflow
- Stream write deadline: 5 seconds; heartbeat cadence: 15 seconds
- Browser buffer: 1,000 records and 2 MiB

#### Scenario: A burst exceeds capacity

- **WHEN** producers emit more events than admission and retention capacity
- **THEN** each boundary stays within both its count and byte limit
- **AND** drops and eviction are visible without blocking collection or allocating an unbounded overflow list
- **AND** oversized messages are truncated before managed retention without retaining a hidden raw copy

#### Scenario: One reader stops consuming

- **WHEN** one SSE client stalls while another reads normally
- **THEN** the stalled client hits the bounded queue/write policy and is disconnected
- **AND** the other client and runtime producers continue
- **AND** exceeding the concurrent-stream limit returns a safe capacity response without allocating another subscriber

### Requirement: Bounded files preserve early diagnostics without changing data durability

Desktop mode SHALL initialize managed diagnostics before database and runtime initialization and write only safe events to its configured diagnostic path. The default SHALL be `logs/<opaque-db-id>/gateway-runtime.jsonl` under the selected data root, with a distinct namespace per canonical database. Rotation SHALL retain one active file and two backups, each no larger than 5 MiB, with a total limit of 15 MiB; only files owned by this log sink can be rotated. Every configured log file set SHALL acquire exclusive ownership and reject overlapping active/backup sets or aliases to protected database/WAL/SHM/data files before opening or rotating. Identity checks MUST prevent symlink/reparse-target changes from redirecting destructive file operations. The web API SHALL read the current-process ring only, not arbitrary or historical files. Logs MUST NOT substitute for audit or data-delivery persistence.

#### Scenario: Rotate at the size boundary

- **WHEN** appending a safe event would exceed the active file limit
- **THEN** the sink rotates before appending and retains at most the active file and two bounded backups
- **AND** it never deletes database, outbox, receipt, audit, or unrelated logger files

#### Scenario: Diagnostic storage fails after startup

- **WHEN** the diagnostic disk becomes full, blocked, or unwritable during acquisition
- **THEN** runtime producers and the bounded memory view continue independently
- **AND** tray/log status exposes diagnostic degradation and file-loss counters
- **AND** recovery reports the lost interval without claiming complete persisted history

#### Scenario: Separate databases do not share a rotation set

- **WHEN** two permitted instances use different databases in the same parent directory
- **THEN** their default log namespaces and file ownership are distinct
- **AND** explicitly configuring a shared or overlapping active/backup set causes a safe startup failure rather than cross-instance rotation

#### Scenario: A configured log path aliases protected data

- **WHEN** LOG_FILE or a rotation target resolves to the database, its WAL/SHM files, another protected data file, an unowned existing file, or an unverified symlink/reparse target
- **THEN** the configuration is rejected before a log open, truncate, rename, or delete can modify that target
- **AND** changing a link between validation and use cannot bypass the protected-file check

#### Scenario: Initial diagnostics cannot be saved

- **WHEN** desktop startup cannot initialize its diagnostic file
- **THEN** startup fails visibly through native error reporting before acquisition
- **AND** the message explicitly says that no diagnostic file was successfully saved
- **AND** it does not silently select another database or claim a web log page is already available

### Requirement: Log data is restricted to the actual local connection

Both GET `${API_BASE_PATH}/system/logs` and GET `${API_BASE_PATH}/system/logs/stream` SHALL independently check access before reading or subscribing. The actual socket peer MUST be numeric loopback, Host MUST be an allowed loopback authority with the actual listener port, and any supplied Origin MUST exactly match the request origin. Sec-Fetch-Site SHALL be absent or exactly one valid `same-origin` or `none` value; `same-site`, `cross-site`, duplicate, and malformed values SHALL be denied. Forwarded, X-Forwarded-* and X-Real-IP headers SHALL be rejected. The routes MUST override global wildcard CORS and MUST NOT introduce tokens, cookies, or readiness-token authentication. This is a local-machine trust boundary, not user authentication or protection from local processes and deliberately relayed local proxies.

#### Scenario: Normal embedded same-origin reads work

- **WHEN** the embedded page accesses logs over IPv4 or IPv6 loopback with the correct Host/port and no Origin header or an exact same-origin Origin
- **THEN** snapshot and SSE requests can pass the local guard
- **AND** responses use no-store, nosniff, and same-origin resource policy without wildcard or cross-origin credential headers

#### Scenario: Spoofed proxy and LAN requests are denied

- **WHEN** a LAN peer sends a loopback-looking forwarded address, or a request includes forwarding headers
- **THEN** both log endpoints return a fixed safe 403 before reading/subscribing
- **AND** Gin proxy trust or the global CORS configuration does not grant access

#### Scenario: A hostile browser origin or Host is denied

- **WHEN** a loopback connection carries a DNS-rebinding Host, wrong port, foreign Origin, or disallowed Sec-Fetch-Site, including `same-site` with no Origin
- **THEN** both log endpoints deny the request without exposing records or raw diagnostic details
- **AND** the global wildcard CORS middleware cannot override that decision

#### Scenario: Remote-only binding does not expand log exposure

- **WHEN** the gateway only listens on a specific non-loopback address
- **THEN** the tray/page explains that local log viewing is unavailable under that binding
- **AND** the feature does not add a listener, change HOST, or expose unauthenticated remote logs

### Requirement: Snapshot queries are bounded and truthful about their scope

The snapshot API SHALL accept level/source/safe-text search, limit, and a before cursor over the current retained safe records. Level SHALL mean minimum severity in debug/info/warn/error order; source SHALL match a stable source ID exactly; text search SHALL use Unicode case-folded literal substring matching over safe message, code, and allowlisted string fields. The default limit SHALL be 200 and maximum 500; out-of-range limits, search longer than 256 characters, regex operations, and malformed cursors SHALL be rejected with safe 400 responses. Results SHALL be chronological and include instance identity, oldest/latest cursors, capture-drop/eviction counters, and sink health. The UI MUST state that search covers limited current-process history rather than all application history.

#### Scenario: Filter a retained window

- **WHEN** an operator selects a severity/source and searches safe text
- **THEN** the result contains only matching retained events within the requested bounded page
- **AND** an empty successful result is distinct from unavailable or denied access

#### Scenario: Reject an unbounded query

- **WHEN** a caller supplies an excessive page limit, excessive search string, regex operation, or invalid cursor
- **THEN** the API rejects it before scanning/allocating an unbounded result
- **AND** it returns a stable safe error rather than raw parser or backend details

### Requirement: Streaming preserves boundaries and declares loss

The stream SHALL use the same level/source/search filters as its snapshot and support an initial after cursor and reconnect Last-Event-ID using instance-plus-string-sequence IDs. Conflicting cursors SHALL be rejected. Snapshot watermark capture and later replay-cutoff capture/live registration SHALL each be atomic with appends, without holding a broker lock across HTTP requests, serialization, or network I/O. A bounded replay copy SHALL precede live events beyond its cutoff and SHALL NOT be enqueued into the smaller live queue. Event IDs SHALL support deduplication. The server MUST explicitly report expired history, admission loss, replay-limit loss, subscriber overflow, and a new process instance through typed gap/reset metadata; a bounded available tail MUST NOT be represented as complete history. Filtered-out sequence IDs MUST NOT themselves be classified as loss. Unknown future cursors SHALL fail safely. A stalled connection that cannot receive gap metadata SHALL be disconnected within its write bound, with loss exposed through bounded counters and the next handshake instead of an unbounded attempt to send a final event.

#### Scenario: Append between snapshot and subscription

- **WHEN** records arrive while a client transitions from snapshot to streaming
- **THEN** retained eligible records after the snapshot watermark are replayed before live delivery without an unexplained missing interval
- **AND** duplicates can be removed using stable event IDs

#### Scenario: Resume after eviction or process restart

- **WHEN** a reconnect cursor is older than retention or belongs to a prior process instance
- **THEN** the server sends a typed gap or reset and only the bounded available tail
- **AND** the UI identifies the missing history rather than implying continuous capture

#### Scenario: A filter matches no new records

- **WHEN** valid records arrive but none match a stream's filter
- **THEN** bounded handshake/heartbeat metadata advances its processed watermark independently from its latest matching record
- **AND** no filtered-out sequence alone is reported as loss or replayed forever
- **AND** progress never advances past an unsent matching record unless an explicit gap has accounted for it

#### Scenario: Replay is larger than the live queue

- **WHEN** a valid reconnect needs more records than the live queue can hold
- **THEN** replay uses its separate bounded copy without holding the broker lock over network writes
- **AND** a replay exceeding 500 records or 1 MiB produces an explicit replay-limit gap and a bounded tail
- **AND** runtime producers are not blocked by replay network I/O

#### Scenario: Close subscriptions safely

- **WHEN** a browser leaves the page, changes its query, disconnects, or the gateway begins shutdown
- **THEN** its subscription, timers, and writer resources are released
- **AND** late writes and concurrent cancellation do not panic, leak subscribers, or keep shutdown open indefinitely

### Requirement: The embedded viewer exposes bounded operator controls

The existing SPA SHALL expose `/studio/logs` and reachable links from setup/runtime without altering existing route identities. The viewer SHALL support severity and source filters, bounded text search, pause/resume, follow-tail, and browser-only clear. Pausing or clearing a view MUST NOT stop acquisition, delete server records, or change logging levels. All content SHALL be rendered as text with accessible controls and en/zh-TW messages. It MUST distinguish loading, empty, connected, reconnecting, paused, denied, failed, gap, truncation, and degraded diagnostic storage states.

#### Scenario: Pause, clear, and resume

- **WHEN** the operator pauses the view, clears its visible rows, and resumes after more events arrive
- **THEN** the gateway has continued acquisition and logging throughout
- **AND** resume uses the last consumed cursor and shows any retention gap
- **AND** clear changes only browser state and neither files nor the server ring

#### Scenario: Filters change during an outstanding response

- **WHEN** filter A has an outstanding request and the operator selects filter B
- **THEN** the old request/stream is canceled and a new snapshot/stream is established
- **AND** late A data cannot replace B results or create duplicate reconnect timers

#### Scenario: Reconnect truthfully after a dropped connection

- **WHEN** streaming disconnects
- **THEN** the view retains its last data with a disconnected/reconnecting indication until a new server handshake
- **AND** one reconnect owner retries with jittered 1/2/4/8/16/30-second capped backoff
- **AND** a safe 403 stops automatic retries and explains local-only access

#### Scenario: Malicious text is not executable

- **WHEN** a safe retained message contains HTML-looking text or script delimiters
- **THEN** the viewer renders the characters literally without executing them
- **AND** the browser buffer remains independently bounded while following the tail
