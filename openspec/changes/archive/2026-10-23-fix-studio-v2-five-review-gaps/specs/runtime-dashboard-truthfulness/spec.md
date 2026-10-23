## ADDED Requirements

### Requirement: Focused header claims only evidenced stages
The focused header SHALL distinguish available configuration, collector running, stream state and verified SQL delivery. It SHALL NOT claim database persistence from availability or collector running alone.

#### Scenario: No verified database commit
- **WHEN** a device is waiting for its first row, Share-only, stopped, unavailable or failing delivery without current SQL evidence
- **THEN** the header SHALL NOT claim database writes and SHALL retain access to existing recording evidence.
