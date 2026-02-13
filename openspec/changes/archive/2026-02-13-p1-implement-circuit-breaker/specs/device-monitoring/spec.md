# device-monitoring Spec Delta

## ADDED Requirements

### Requirement: Device health profiling

The system SHALL track the connectivity health of each configured device using a sliding window of historical request results.

#### Scenario: Track success rate

- GIVEN a device has processed 100 requests
- WHEN 5 requests fail
- THEN the system calculates a 5% error rate for that time window

### Requirement: Circuit breaker state machine

The system SHALL maintain a circuit breaker state (Healthy, Unstable, Dead) for each device to protect system resources.

#### Scenario: Trip circuit breaker

- GIVEN a device is in Healthy state
- WHEN the error rate exceeds the configured threshold (e.g., 20%)
- THEN the system transitions the device to Dead state and suspends regular polling

#### Scenario: Background probing recovery

- GIVEN a device is in Dead state
- WHEN the background prober successfully connects to the device
- THEN the system transitions the device back to Healthy state and resumes polling
