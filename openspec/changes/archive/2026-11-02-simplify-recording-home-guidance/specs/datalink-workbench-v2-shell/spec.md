## MODIFIED Requirements

### Requirement: Design tokens scope
The v2 shell SHALL retain scoped existing typography and brand semantics while the new recording surfaces and affected reused controls use root DESIGN.md light/dark semantic tokens under the shared ThemeProvider. Tokens MUST NOT leak local layout styling to unrelated routes. New surfaces SHALL avoid forced dark color-scheme, decorative gradients, glow and glass blur. The existing data-workbench-v2 root SHALL remain; Inter/system and JetBrains Mono typography SHALL be preserved. The retired /studio route MUST NOT be restored.

#### Scenario: Tokens scoped to v2 subtree
- **WHEN** the operator navigates from /studio/v2 to an existing unrelated route
- **THEN** that route retains its existing layout and shared theme preference without inheriting V2-specific canvas or layout styling.

#### Scenario: Tokens applied on v2
- **WHEN** the new recording surface mounts in either resolved theme
- **THEN** it uses the data-workbench-v2 scope, shared semantic tokens, Inter/system text and JetBrains Mono for code, without a forced dark canvas.
