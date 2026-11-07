## ADDED Requirements

### Requirement: Selected memory facts remain readable without expanding address rows
The five memory views SHALL retain compact physical-address cells when an address is selected or inspected. Complete selected facts SHALL remain available in a separate visible, keyboard-reachable information region, with stable reserved space so focus changes do not move action controls. Cell accessible names and exact datatype, address, value and evidence state MUST remain complete. Continuation cells and the information region SHALL use the same canonical value identity for visible-value observation, without introducing another live-value owner.

#### Scenario: A long selected description is displayed
- **WHEN** the operator selects a cell whose rule name, full description or value is long
- **THEN** the selected facts remain readable outside the address grid without increasing the other cells in that row
- **AND** focus, blur, pointer inspection and touch do not move the action controls between press and activation.

#### Scenario: A continuation is visible after a window change
- **WHEN** a visible cell continues a confirmed multi-register value whose canonical start lies outside the current window
- **THEN** the visible subscription and selected facts use that canonical start and preserve the current value and freshness semantics
- **AND** unknown, partial, conflict, stale and failed evidence retain existing fail-closed behavior rather than fabricated successful data.
