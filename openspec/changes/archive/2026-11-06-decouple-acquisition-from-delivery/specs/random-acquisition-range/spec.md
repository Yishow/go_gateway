## Purpose

Let devices that support random reads acquire many scattered addresses of different areas and data types in as few requests as possible, with MC 3E word, double word and bit access, a documented per-profile random limit that operators can lower, and a range form that warns when a selection approaches that limit.

## ADDED Requirements

### Requirement: MC random read covers word, double word and bit devices

The MC 3E random read (command 0403, subcommand 0000) SHALL send separate word and double word access sections. 16-bit items SHALL use word access; 32-bit items SHALL use double word access; bit devices SHALL be read in word units of 16 points aligned to a multiple of 16 and the requested bit SHALL be extracted from the returned word. Long timer and long counter devices (LTS, LTC, LSTS, LSTC, LCS, LCC) SHALL NOT be placed in a random request. LTN devices SHALL be rejected with an explicit unsupported reason under the supported subcommand 0000 boundary; no request for an LTN device SHALL be sent with subcommand 0000 or 0002. Strings SHALL NOT use random access. The word access points plus double word access points of one request SHALL NOT exceed the profile limit.

#### Scenario: Mixed word, double word and bit in one request

- **GIVEN** an MC device declared with a random-capable profile and points D39 int16, D100 int32 and M5 bool
- **WHEN** one scan executes
- **THEN** one 0403 request carries D39 and M0 (16-point unit) in the word section and D100 in the double word section
- **AND** M5 is decoded from bit 5 of the M0 word

#### Scenario: Long timer is excluded

- **WHEN** a scattered range contains an LTN item
- **THEN** planning reports an explicit unsupported reason for that item and sends no request for it

##### Example: long timer beside random items

- **GIVEN** a random-capable MC device with scattered items D10 int16, D20 int16 and LTN0
- **WHEN** the plan is built
- **THEN** planning rejects LTN0 with an explicit unsupported reason and sends no request for it; rejecting the whole input with that reason satisfies this boundary

### Requirement: Random profiles carry documented limits

The acquisition profile catalog SHALL include `mitsubishi-q-iqr-3e-binary-v1` with a random limit of 192 points per request (word plus double word access points, subcommand 0000, per MELSEC Communication Protocol Reference Manual SH-080008), excluding the ZR file register. `mitsubishi-fx5-cpu-3e-binary-v1` SHALL keep 192. `fatek-fbs-ascii-v1` SHALL keep at most 64 components with a weighted cost of at most 64 words, a 32-bit component costing 2. The optional `acquisition.limits.random_points` SHALL lower the per-request random limit and SHALL be rejected when it is not a positive integer or exceeds the profile limit. The connection settings SHALL show the profile limit as the default value and allow lowering it.

#### Scenario: Limit cannot be raised

| Profile | Protocol limit | `random_points` | Result |
| ------- | -------------- | --------------- | ------ |
| mitsubishi-q-iqr-3e-binary-v1 | 192 | 120 | accepted, requests hold at most 120 points |
| mitsubishi-q-iqr-3e-binary-v1 | 192 | 300 | rejected with a field error |
| fatek-fbs-ascii-v1 | 64 weighted words | 48 | accepted |

- **WHEN** an operator saves `random_points` above the profile limit
- **THEN** the save is rejected and the stored declaration is unchanged

### Requirement: Scattered mixed rules are persisted and acquired

Source rules SHALL persist `layout_mode` and `mixed_items` through the source rule repository and the Studio v2 workspace rule API. Each mixed item SHALL carry its own address and data type and items of one rule SHALL be allowed to use different areas. Saving SHALL validate item spans, data types and addresses and SHALL reject overlapping items with a field error naming the conflicting items. Each mixed item SHALL be expanded into one point that keeps its item identity. For a device whose declared profile supports random reads, acquisition planning SHALL accept mixed rules and plan their items through the optimizer; for undeclared, `safe-v1` or Modbus devices, mixed rules SHALL remain excluded with the reason `source-layout-unsupported`. Homogeneous rules SHALL keep their existing behavior.

#### Scenario: Mixed rule round trip

- **GIVEN** an MC device declared with `mitsubishi-q-iqr-3e-binary-v1`
- **WHEN** a rule is saved with items D39 int16, D100 int32 and M5 bool
- **THEN** reading the rule back returns the same three items and three points exist for it

#### Scenario: Overlap is rejected

- **WHEN** a mixed rule contains D3 uint32 and D4 int16
- **THEN** saving is rejected with a field error naming both items

#### Scenario: Safe profile keeps excluding mixed rules

- **WHEN** a device declared with `safe-v1` has a mixed rule
- **THEN** acquisition planning excludes that rule with the reason `source-layout-unsupported`

### Requirement: Scattered address range with usage warning

For a device whose declared profile supports random reads, the read range form SHALL offer a scattered address mode that stores a mixed layout of items, each with its own address and data type, across areas. The form SHALL show the computed random usage against the effective limit: for MC, each 16-bit or 32-bit item counts 1 point and bit items count one point per distinct 16-point unit; for FATEK, 16-bit items count 1 and 32-bit items count 2. At 90 percent or more of the limit the form SHALL show a warning; above the limit it SHALL explain that acquisition will be split into multiple requests without blocking the save. For devices without a random-capable profile the mode SHALL be unavailable with a stated reason.

#### Scenario: Warning at 90 percent

- **GIVEN** an MC device with an effective random limit of 192
- **WHEN** the operator's scattered range reaches 173 points
- **THEN** the form warns that usage is at least 90 percent of 192

#### Scenario: One range replaces several

- **GIVEN** a random-capable MC device
- **WHEN** the operator creates one scattered range with D39, D40, D100, D200, D1000 and D2000 next to an existing D0–D31 range
- **THEN** acquisition uses at most one contiguous request for D0–D31 and one random request for the six scattered items, or a single random request when that is cheaper
