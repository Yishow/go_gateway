# hsllogic-performance Spec Delta

## ADDED Requirements

### Requirement: Packet logger object pooling

The `PacketLogger` component SHALL use object pooling (`sync.Pool`) for `PacketLog` entries to minimize memory allocations during high-throughput logging.

#### Scenario: High-frequency packet logging

- GIVEN the system is logging packets at a rate of 1000+ per second
- WHEN a new packet is logged
- THEN the logger acquires a `PacketLog` object from a pool instead of allocating new memory

#### Scenario: Pool object cleanup

- GIVEN a pooled `PacketLog` object is returned to the pool
- WHEN it is retrieved again for a new log entry
- THEN all previous data fields are cleared to prevent data leakage

---

### Requirement: Address parsing cache

The system SHALL cache parsed address structures to avoid redundant parsing of the same address string during repeated polling cycles.

#### Scenario: Cached address retrieval

- GIVEN an address string `D100` has been parsed once
- WHEN the same address is requested again
- THEN the system returns the cached `ParsedAddress` without re-parsing

#### Scenario: Cache invalidation on configuration change

- GIVEN the address cache contains entries
- WHEN the device configuration is reloaded
- THEN the cache is cleared

---

## Cross-References

- Related capability: `protocol-connectors` (uses parsed addresses)
- Related capability: `collection-scheduler` (orchestrates polling)
