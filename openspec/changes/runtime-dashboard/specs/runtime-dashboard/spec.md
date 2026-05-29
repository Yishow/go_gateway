## ADDED Requirements

### Requirement: Navigation and Layout
The system SHALL provide a dedicated dashboard route at `/studio/dashboard` and a navigation link in the top bar of the main layout, presenting a structured layout containing metrics, device status grid, live values table, output queue status, and real-time logs.

#### Scenario: User navigates to dashboard
- **WHEN** user clicks on the "Dashboard" navigation link in the top bar
- **THEN** the system redirects to `/studio/dashboard` and renders all five main widget panels

### Requirement: Real-time Metrics Processing
The system SHALL poll or stream real-time performance indicators, calculating and showing the polling throughput (points per second), success rate (%), and average polling latency (ms).

#### Scenario: Dashboard metrics update
- **WHEN** new performance metrics are received from the backend
- **THEN** the system updates the display of points per second, success rate, and latency metrics instantly

##### Example: Metrics computation
| Received Points (1s window) | Failed Polls | Latencies | Computed Points/Sec | Computed Success Rate | Avg Latency |
|-----------------------------|--------------|-----------|---------------------|-----------------------|-------------|
| 100                         | 5            | [10, 20]  | 100                 | 95.0%                 | 15.0ms      |
| 0                           | 0            | []        | 0                   | 100.0%                | 0.0ms       |

### Requirement: Device Status Grid Visualization
The system SHALL render a grid of connected PLC devices, indicating each device's name, protocol type, connectivity status (online, offline, testing), and current latency value.

#### Scenario: Device connectivity changes
- **WHEN** a device's connection status changes from online to offline
- **THEN** the status grid item for this device turns red, displaying "Offline" with an updated timestamp

### Requirement: Live Value Streaming
The system SHALL dynamically display a live values table containing point names, Modbus/PLC addresses, raw values, transformed values, and updated times.

#### Scenario: Value transform and display
- **WHEN** the raw register value changes
- **THEN** the system applies the scale and offset multiplier and shows the updated raw and transformed values in the table

##### Example: Linear scale transformation
- **GIVEN** a point with scale_multiplier = 0.1 and scale_offset = 5
- **WHEN** the system reads a raw value of 150
- **THEN** the transformed value is calculated as 20.0 (150 * 0.1 + 5)

### Requirement: Output Queue Backlog Indicator
The system SHALL monitor the database and Local Modbus output queues, rendering the current backlog count and alerting the user if the queue length exceeds the threshold.

#### Scenario: Queue backlog alert
- **WHEN** the output queue backlog count exceeds 500 items
- **THEN** the system displays a warning warning banner and changes the queue indicator color to amber

### Requirement: Real-time Log Stream
The system SHALL stream recent errors, timeouts, or protocol parse warnings from the backend Datalink scheduler, showing them in a scrollable console panel.

#### Scenario: Error log output
- **WHEN** the backend collector encounters a TCP handshake timeout for a device
- **THEN** the dashboard appends a timestamped error message to the log stream console in red text
