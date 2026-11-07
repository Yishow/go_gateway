## Purpose

Provide five interchangeable memory presentations for source setup while preserving authoritative occupancy, address geometry, live values and safe actions.

## ADDED Requirements

### Requirement: Layout changes preserve shared memory context
The memory view SHALL expose grid focus, range navigation, horizontal map, range chapters and grid with detail layouts. Switching SHALL retain selected area, address window, focus and owner selection without restarting the parent live-value hook or mutating configuration.

#### Scenario: Switch all layouts at a located address
- **WHEN** an operator locates a valid address and switches through all five layouts
- **THEN** each layout displays the same area, address window, ownership and current values without configuration writes

### Requirement: Every layout preserves safe cell actions
Every layout SHALL use the same occupancy authority and action callbacks. Empty cells SHALL offer address-prefilled creation; used cells SHALL open their owner; conflicts SHALL require owner selection; unconfirmed cells MUST NOT permit creation. Recovery, local and stale blocking SHALL remain effective.

#### Scenario: Occupancy and blocked actions remain identical
- **WHEN** the operator switches layouts with empty, used, conflicting and unconfirmed cells or a blocked source
- **THEN** actions remain identical to the shared occupancy state and unknown cells cannot invoke creation

### Requirement: Typed values retain physical address geometry
The system SHALL distinguish data types by text and visual treatment without a left color bar. Every physical register SHALL retain its address and focus target. Adjacent registers of a confirmed value SHALL visibly connect with continuation information at row and page boundaries. Unknown widths MUST NOT be presented as confirmed spans.

#### Scenario: Multi-register values wrap without losing identity
- **WHEN** one-register integers, two-register floats and four-register values cross a displayed row or page
- **THEN** each visible register identifies its true address and value start while joining only adjacent confirmed registers of that value

#### Scenario: Uncertain string width remains uncertain
- **WHEN** source evidence reports a string-width issue
- **THEN** the known partial start and unconfirmed remainder remain distinct without inventing a confirmed string span

### Requirement: Layout-specific navigation shares the current data
Range navigation and chapters SHALL derive summaries and jump actions from the same area rules. The horizontal map SHALL represent the current bounded window. Detail SHALL describe the selected cell and expose only its permitted actions.

#### Scenario: Auxiliary navigation changes the shared window
- **WHEN** the operator chooses a saved range or an address in the occupancy map
- **THEN** the common window or selected address changes consistently and no duplicate memory subscription is created

### Requirement: Themes and responsive access preserve functionality
All layouts SHALL consume the existing light and dark theme and preserve keyboard access, selected tab visibility, focus restoration and all information at desktop, tablet and phone widths. Layout changes SHALL refresh rendered-column measurement and visible-value observation.

#### Scenario: Narrow or themed layout keeps interaction valid
- **WHEN** the viewport or theme changes in any layout
- **THEN** controls remain reachable, keyboard movement follows rendered columns, live values remain observed and unknown occupancy remains non-activatable
