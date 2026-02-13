# collection-optimizer Specification

## Purpose
TBD - created by archiving change p2-optimize-packet-aggregation. Update Purpose after archive.
## Requirements
### Requirement: Static request optimization

The system SHALL optimize the polling strategy at startup by aggregating individual tag requests into optimized memory blocks based on protocol constraints.

#### Scenario: Aggregating contiguous registers

- GIVEN tags configured at addresses D100, D101, and D102
- WHEN the request optimizer runs
- THEN it generates a single read request starting at D100 with length 3

#### Scenario: Aggregating with holes

- GIVEN tags at D100 and D105, and `max_hole_size` is 10
- WHEN the optimizer runs
- THEN it generates a single read request for D100 length 6 (discarding D101-D104)

### Requirement: PDU awareness

The optimizer SHALL respect the specific Protocol Data Unit (PDU) limits of each driver to prevent packet fragmentation or rejection.

#### Scenario: Respect Modbus limit

- GIVEN a contiguous block of 200 registers
- AND the Modbus driver specifies a limit of 125 registers per read
- THEN the optimizer splits the block into two requests (e.g., 125 + 75)

