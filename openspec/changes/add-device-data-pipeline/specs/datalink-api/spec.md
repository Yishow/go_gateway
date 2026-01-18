## ADDED Requirements
### Requirement: Management APIs
The system SHALL expose APIs for devices, points, tags, and mappings.

#### Scenario: Create mapping via API
- WHEN a client posts a mapping
- THEN the mapping is saved and returned with an id

### Requirement: Connection and preview APIs
The system SHALL expose APIs to test device connections and preview mapping results.

#### Scenario: Mapping preview API
- WHEN a client requests a preview for a mapping
- THEN the API returns raw and transformed values

### Requirement: Write settings API
The system SHALL expose APIs to read and update write settings, including timestamp precision and partition interval.

#### Scenario: Update write precision
- WHEN a client updates the write precision to milliseconds
- THEN subsequent writes use millisecond precision
