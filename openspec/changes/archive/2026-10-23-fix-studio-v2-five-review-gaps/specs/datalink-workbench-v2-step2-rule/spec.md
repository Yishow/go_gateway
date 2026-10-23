## ADDED Requirements

### Requirement: Complete occupied range without changing points
The rule summary SHALL distinguish item start positions from the complete occupied range using the source stride, datatype span and existing protocol address representation, without changing actual points.

#### Scenario: Multiword occupied end
- **WHEN** four float32 items start at D0 or 40001 with contiguous two-word stride
- **THEN** the complete occupied range SHALL end at D7 or 40008 respectively, while still reporting four items and eight words.

#### Scenario: Other widths and address radices
- **WHEN** 64-bit, single-word or bit values are summarized in decimal or protocol-specific radix areas
- **THEN** the summary SHALL include the final occupied unit and preserve the parser radix and actual point addresses.
