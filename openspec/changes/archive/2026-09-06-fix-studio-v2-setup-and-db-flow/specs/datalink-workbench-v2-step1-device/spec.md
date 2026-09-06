## ADDED Requirements

### Requirement: Protocol switch auto-populates dedicated communication presets

When adding or editing a device in Step 1, selecting any supported communication protocol SHALL auto-populate dedicated default communication parameters for that protocol:
- `mc_3e` (MC Protocol): `port` SHALL default to `6000`, `station` SHALL default to `0`, `network_no` SHALL default to `0`, `pc_no` SHALL default to `255`, `module_io_no` SHALL default to `1023`, `module_station_no` SHALL default to `0`.
- `fatek_fbs`: `port` SHALL default to `500`, `station` SHALL default to `1`.
- `modbus_tcp`: `port` SHALL default to `502`, `slave_id` SHALL default to `1`, `timeout` SHALL default to `5`.
- `modbus_udp`: `port` SHALL default to `502`, `slave_id` SHALL default to `1`, `timeout` SHALL default to `5`.
- `modbus_rtu`: `port` SHALL default to `/dev/ttyUSB0`, `baud` SHALL default to `9600`, `parity` SHALL default to `N`, `slave_id` SHALL default to `1`.
- `mqtt`: `broker` SHALL default to `mqtt://127.0.0.1:1883`, `client_id` SHALL default to `gw-01`.

#### Scenario: Selecting MC Protocol in Step 1
- **WHEN** the operator selects protocol `mc_3e` in the Device creation modal or editor
- **THEN** the port field SHALL default to `6000`
- **AND** the station field SHALL default to `0`

#### Scenario: Selecting Fatek Protocol in Step 1
- **WHEN** the operator selects protocol `fatek_fbs`
- **THEN** the port field SHALL default to `500`
- **AND** the station field SHALL default to `1`

#### Scenario: Selecting Modbus RTU Protocol in Step 1
- **WHEN** the operator selects protocol `modbus_rtu`
- **THEN** the serial port field SHALL default to `/dev/ttyUSB0`
- **AND** the baud rate SHALL default to `9600`
- **AND** parity SHALL default to `N`

### Requirement: Protocol form frontend-backend field alignment

The Device Connection form in Step 1 SHALL accurately render protocol-specific fields and labels, and SHALL serialize parameters according to each backend protocol driver specification without missing or mismatched field names.

#### Scenario: Submitting MC Protocol device
- **GIVEN** an MC Protocol device with port `6000` and station `0`
- **WHEN** the operator saves the device
- **THEN** the payload sent to backend SHALL contain `port: 6000` and station configuration matching the MC Protocol backend driver specification
