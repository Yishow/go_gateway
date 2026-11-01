## MODIFIED Requirements

### Requirement: Static request optimization
The system SHALL produce deterministic immutable protocol-aware plans for the active acquisition Points at startup and each relevant configuration revision. It SHALL execute these plans through the existing collector, not a second acquisition engine. Recording membership SHALL NOT change the acquisition set. Plans SHALL group only compatible device, station, connection generation, access function/area, protocol capability and polling cadence. Invalid addresses SHALL produce explicit errors, never silent omission.

#### Scenario: Aggregating contiguous registers
- **GIVEN** active int16 Points D100, D101 and D102 with compatible capability and cadence
- **WHEN** the planner runs
- **THEN** it produces one contiguous read covering all three and emits each Point separately.

#### Scenario: Aggregating with holes
- **GIVEN** Modbus offsets 0, 10 and 20 with max_hole_size 10 but no explicit readable gap coverage
- **WHEN** the planner runs
- **THEN** it produces three reads, not a guessed covering interval.

#### Scenario: Explicit readable coverage
- **GIVEN** the same three Points and versioned explicit readable coverage of 0 through 20 without excluded or forbidden cells
- **WHEN** protocol and configured limits permit the interval
- **THEN** one 21-word request reads the range and emits only the three original Points.

#### Scenario: Source range is not recording selection
- **GIVEN** 21 active int16 Points D0 through D20 and a record selecting only D0, D10 and D20
- **WHEN** the record changes membership
- **THEN** all 21 acquisition Points remain scheduled and no duplicate source reader is created.

### Requirement: PDU awareness
The planner SHALL enforce the strictest protocol, client and configured-device item, byte and address limits, using native address radix and component width. A typed scalar SHALL fit completely within one request; an individually oversized scalar SHALL fail explicitly. It SHALL preserve every original Point mapping when unioning nested or overlapping spans.

Byte budgets SHALL be immutable Protocol, Client, and Device pairs of RequestBytes and ResponseBytes. Each direction SHALL use the strictest positive limit; zero SHALL mean unspecified, not a proven zero cost. Costs SHALL cover the complete supported application frame, including Modbus MBAP or RTU unit and CRC, MC3E binary headers and end code, or FATEK ASCII framing and checksum, and SHALL exclude IP, TCP, and Ethernet overhead. Costs SHALL be derived from the actual command serializers and response fixtures, not a generic payload-size formula. Unknown frame costs SHALL be reported as unknown; a positive applicable budget with unknown cost SHALL reject that candidate with unsupported-frame-cost rather than claim compliance. Native item limits SHALL also be clamped to the actual client implementation, not only configured or theoretical protocol limits. Cost-model and budget changes SHALL participate in plan and capability identity.

#### Scenario: Complete frame budget boundaries
- **GIVEN** a supported command whose request and normal response frame lengths are verified against its actual serializer and fixtures
- **WHEN** Protocol, Client, and Device provide different positive request or response byte limits
- **THEN** each direction uses its smallest positive limit; an exactly fitting frame is permitted and an over-limit frame is split only at complete scalar boundaries or explicitly rejected if indivisible
- **AND** native bit packing, component width, address encoding length, and actual client item caps remain enforced independently of byte budgets

#### Scenario: Unknown frame cost cannot pass a budget
- **GIVEN** a command or frame variant without a verified frame-cost model
- **WHEN** an applicable request or response byte budget is positive
- **THEN** the candidate is rejected with unsupported-frame-cost, not accepted using zero or payload-only bytes
- **AND** an unspecified budget does not turn an unknown cost into a known estimate or enable an otherwise unsupported command

#### Scenario: Respect Modbus limit
- **GIVEN** 200 contiguous int16 registers and a 125-register capability limit
- **WHEN** the plan is built
- **THEN** it produces requests of 125 and 75 without crossing a forbidden boundary.

#### Scenario: Nested overlapping spans
- **GIVEN** D0 int64 and D1 int16 on compatible word memory
- **WHEN** their raw spans are deduplicated
- **THEN** the union remains D0 through D3 and neither Point loses required words.

FATEK command 48 SHALL distinguish native component count from its independent 64-word protocol budget. Its N field SHALL equal the number of requested components; each discrete or 16-bit component SHALL cost one word and each 32-bit component SHALL cost two words. More restrictive configured component limits and byte limits SHALL remain independently enforced. Native address increments SHALL follow the protocol address space: DR/DD components advance two underlying word addresses, W-prefixed bit components advance sixteen bit addresses, and DW-prefixed bit components advance thirty-two bit addresses. Plain word and discrete components advance one address. The supported WX/WY/WM/WS and DWX/DWY/DWM/DWS start addresses SHALL be multiples of eight, not necessarily multiples of their component stride; complete scalar footprints SHALL remain within the declared address bounds. Complete scalar expansion, contiguous compatibility, random deduplication, frame costs and result scatter SHALL use this same address contract. A valid overlapping or odd-address component SHALL NOT be silently replaced by another component or rejected merely to conceal a stride mismatch.

#### Scenario: FATEK component count is not word cost
- **WHEN** a random request contains twenty 32-bit components and twenty-four discrete or 16-bit components
- **THEN** its N is 44 and its word cost is 64; it fits only when all independent configured limits also permit it
- **AND** thirty-three int32 Points at distinct DR wire addresses exceed 64 words and split into at least two requests without splitting a scalar, even when the configured component cap is 64

#### Scenario: Native component addresses preserve scalar identity
- **WHEN** DR0 and DR10 int64 Points expand into 32-bit native components
- **THEN** their component addresses are DR0/DR2 and DR10/DR12, not DR0/DR1 and DR10/DR11
- **AND** separate valid int32 Points at DR0 and DR1 do not become a contiguous request whose second returned component is DR2
- **AND** a WY8 int32 scalar expands to WY8/WY24 while retaining the protocol's fixed-width address encoding

#### Scenario: Native random limits
- **WHEN** expanded MC random units exceed 255 or FATEK units exceed 64 under those configured client limits
- **THEN** requests split at complete scalar boundaries with no dropped or duplicated Point result.

## ADDED Requirements

### Requirement: Capability-aware random acquisition
The system SHALL integrate existing MC and FATEK random-read clients for verified compatible capabilities. Unknown or unsupported capabilities SHALL use safe contiguous or scalar reads with an explicit reason. MC word-only random restrictions and FATEK native component widths SHALL be respected.

The only persisted capability declaration SHALL be the optional `acquisition` object in the existing Device connection configuration. Version 1 SHALL accept a known profile, optional stricter positive native-unit and directional application-frame limits, and optional revision-bound explicit readable/forbidden/excluded native ranges. Unknown fields, null declarations, unsupported versions/profiles and invalid limits/ranges SHALL be rejected by the existing device validation boundary. A single pure resolver SHALL serve installation and future preview; registry schema SHALL describe that same declaration. Coverage content as well as its revision SHALL participate in identity. No recording selector, client implementation, connection success or speculative read SHALL establish device support.

#### Scenario: Historical configuration stays safe
- **GIVEN** a saved MC or FATEK connection without an acquisition declaration, including a named FATEK subtype
- **WHEN** scheduled acquisition resolves its capabilities
- **THEN** random support remains unknown and only safe contiguous/scalar requests are eligible
- **AND** an explicit `safe-v1` profile disables random without inventing readable holes.

#### Scenario: Verified family and command-specific limits
- **GIVEN** an explicit compatible `fatek-fbs-ascii-v1` or `mitsubishi-fx5-cpu-3e-binary-v1` declaration
- **WHEN** the latter selects FX5 CPU built-in Ethernet 3E binary word random acquisition
- **THEN** its 0403/0000 protocol limit is 192 access points independently of the client's 255-item serializer limit
- **AND** `RandomProtocolLimit` only adds a stricter random-command cap; it never relaxes existing shared limits or constrains contiguous/bit commands by accident
- **AND** unverified MC/Q/L/R/module families remain unknown, while all native-unit, frame-byte and address limits remain independently enforced.


#### Scenario: Three sparse word Points
- **GIVEN** supported MC or FATEK random capability and three int16 Points D0, D10 and D20
- **WHEN** one error-free scan executes
- **THEN** one protocol read request obtains the three Point results.

#### Scenario: Capability unavailable
- **WHEN** random capability is unknown or the requested area is unsupported
- **THEN** the planner does not probe arbitrary addresses and reports its safe fallback plan.

### Requirement: Identity-preserving raw deduplication
The system SHALL deduplicate only compatible raw units within a scan and scatter complete results using each Point's immutable effective format, datatype, identity and acquisition lease. It SHALL NOT reuse an earlier scan's value as a new observation or use one global byte order for all results.

#### Scenario: Distinct decoders over identical words
- **GIVEN** two Points sharing words 0x0102 and 0x0304 with effective ABCD and CDAB formats
- **WHEN** one raw read succeeds
- **THEN** separate results are 0x01020304 and 0x03040102, each with its own Point identity and fingerprint.

#### Scenario: Different cadences
- **GIVEN** Points with 500ms and 2000ms polling cadences
- **WHEN** the fast group becomes due alone
- **THEN** it does not acquire or emit a new observation for the slow group.

### Requirement: Revision-bound request preview
The system SHALL expose a scoped read-only preview generated by the same planner used for execution, bound to configuration and capability revisions. It SHALL report acquisition Point count, unique wire units, planned protocol requests, limits and fallback reasons without PLC I/O. It SHALL distinguish estimates from actual requests, retries, fallback attempts, transport frames and SQL writes. Relevant changes SHALL invalidate old plans and previews.

#### Scenario: Preview parity and stale revision
- **WHEN** preview and error-free execution use the same immutable inputs
- **THEN** their initial protocol request counts match; changed inputs reject stale preview revision with conflict instead of displaying an obsolete estimate.

### Requirement: Reproducible optimization evidence
The system SHALL validate planner strategies against deterministic scalar baselines using identical acquisition sets, simulator data and transport conditions. Reports SHALL include request counts, bytes, retries, per-device cycle latency and freshness, queue/drop counts, CPU/allocation conditions and cold versus warm connections. It SHALL NOT claim an unmeasured speedup.

#### Scenario: Safe sparse and contiguous baselines
- **WHEN** no-error fixtures execute three sparse word Points and 21 contiguous word Points
- **THEN** verified MC/FATEK sparse requests are 3 to 1, uncovered Modbus sparse remains 3 to 3, and permitted contiguous requests are 21 to 1, with equivalent per-Point values and quality.

#### Scenario: Optimization regression
- **WHEN** a strategy increases controlled-RTT cycle duration relative to the identical scalar baseline or changes values
- **THEN** it is not enabled until the regression is resolved and rerun evidence passes.
