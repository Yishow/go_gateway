## Purpose

Let operators set each device's acquisition interval, optional automatic interval adjustment and protocol-appropriate acquisition optimization from its connection settings, with acquisition starting from those settings when the program starts.

## ADDED Requirements

### Requirement: Each device owns its acquisition interval

Each device SHALL have at most one device-owned polling group. Saving a device acquisition interval SHALL create or update that group and assign all of the device's points to it, without changing other devices' points. Points created afterwards for the device by source rules or runtime device synchronization SHALL use the device-owned group, falling back to the default polling group only when none exists. The interval SHALL be validated between 100 ms and 3600000 ms; invalid values SHALL be rejected with a field error. On upgrade, each existing device SHALL receive a device-owned group. Its interval SHALL be the shortest interval among the enabled polling groups used by the device’s assigned points; those points SHALL be reassigned to the device-owned group. Points in disabled groups and points without a group SHALL retain their prior assignment. A device with no such enabled-group points SHALL use the existing default fallback interval. Re-running the migration SHALL NOT create another owned group. Points previously assigned to a longer enabled interval adopt the shortest interval. Acquisition SHALL start at program start using each device's group without a manual action.

#### Scenario: Interval change affects only one device

- **GIVEN** devices A and B both polling every 1000 ms
- **WHEN** device A saves an interval of 500 ms
- **THEN** all of A's points poll every 500 ms and B's points still poll every 1000 ms

#### Scenario: Upgrade keeps frequency

- **WHEN** the migration runs on a workspace where all points use the 1000 ms default group
- **THEN** every device gets a device-owned group of 1000 ms and its points are reassigned to it

#### Scenario: Upgrade combines mixed intervals at the shortest interval

- **GIVEN** device A has points in enabled groups of 500 ms and 2000 ms, and device B has a point in an enabled 2000 ms group
- **WHEN** the migration runs
- **THEN** A has one device-owned 500 ms group and both of its points use that group
- **AND** B retains a 2000 ms acquisition interval
- **AND** points in disabled groups or without a group retain their prior assignment, and re-running the migration creates no duplicate groups

### Requirement: Automatic interval adjustment

When automatic adjustment is enabled for a device, the effective interval SHALL be the larger of the configured interval and the maximum scan duration of the last 5 scans multiplied by 1.2 and rounded up to a 100 ms multiple. After 5 consecutive on-time scans the effective interval SHALL step back toward the configured interval. Runtime status SHALL report both `configured_interval_ms` and `effective_interval_ms`. When disabled, overruns SHALL be reported without fabricating on-time samples.

#### Scenario: Slow device stretches the interval

- **GIVEN** a configured interval of 1000 ms with automatic adjustment enabled
- **WHEN** the last 5 scans took at most 1250 ms
- **THEN** the effective interval becomes 1500 ms and runtime status shows configured 1000 and effective 1500

### Requirement: Connection settings expose protocol-appropriate acquisition optimization

The connection settings SHALL show the acquisition interval, automatic adjustment, stale threshold, frozen threshold and an acquisition optimization profile, each with contextual help. The profile choice SHALL write only the existing `connection_config.acquisition` declaration and SHALL be validated by the existing device validation. Modbus devices SHALL offer no random profile and SHALL state that contiguous blocks are aggregated. MC 3E devices SHALL offer undeclared, `safe-v1`, `mitsubishi-fx5-cpu-3e-binary-v1` and `mitsubishi-q-iqr-3e-binary-v1`; FATEK devices SHALL offer undeclared, `safe-v1` and `fatek-fbs-ascii-v1`. Saving SHALL take effect without restarting the program.

#### Scenario: Protocol-specific options

| Protocol | Offered profiles |
| -------- | ---------------- |
| Modbus TCP/UDP/RTU | none (contiguous aggregation note only) |
| MC 3E | undeclared, safe-v1, mitsubishi-fx5-cpu-3e-binary-v1, mitsubishi-q-iqr-3e-binary-v1 |
| FATEK FBs | undeclared, safe-v1, fatek-fbs-ascii-v1 |

- **WHEN** the operator opens connection settings for an MC 3E device
- **THEN** the profile select offers exactly the MC 3E row's profiles
