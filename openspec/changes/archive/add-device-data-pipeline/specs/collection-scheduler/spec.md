## ADDED Requirements
### Requirement: Scheduled collection
The system SHALL schedule point reads using fixed intervals per polling group.

#### Scenario: Interval read
- WHEN a polling group is set to 5s
- THEN the system executes reads every 5 seconds

### Requirement: Retry and backoff
The system SHALL retry failed reads with a configurable retry count and backoff.

#### Scenario: Retry after failure
- WHEN a read fails
- THEN the system retries up to the configured count before marking an error

### Requirement: Concurrency control
The system SHALL limit concurrent reads per device to avoid protocol collisions.

#### Scenario: Serial device lock
- WHEN two tasks target the same serial device
- THEN the system serializes the reads

### Requirement: Quality and timestamp
The system SHALL attach a timestamp and quality flag to every collected value.

#### Scenario: Bad quality
- WHEN a read times out
- THEN the collected value is marked with a bad quality flag
