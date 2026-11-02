## MODIFIED Requirements

### Requirement: UI / API / Diagnostics options
The General settings SHALL retain three SectionCards for UI, API and Diagnostics.
- UI SHALL retain theme dark/light/auto, locale zh-TW/en with zh-TW default and addr_format modbus/hex/raw with modbus default. A genuinely absent theme preference SHALL resolve to system/auto rather than treating an injected dark mapper default as saved intent.
- API SHALL retain api_base text in mono defaulting to http://localhost:8080 with the no-trailing-slash hint, api_version v1/v2 default v1 and timeout_seconds default 30 seconds.
- Diagnostics SHALL retain log_level trace/debug/info/warn/error default info, sse_heartbeat_seconds default 15, enable_debug_panel default false and enable_audit_log default true.
Changes MUST write through state.settings.general. Existing explicit whole-settings Save behavior SHALL remain unchanged. A theme-only interaction SHALL persist only the theme key through the existing updateKey mutation or equivalent theme-limited adapter; it MUST NOT invoke full snapshot saving or flush other unsaved settings. Theme selection SHALL also update the shared global provider without introducing per-workspace preferences; auto SHALL follow live OS preference. Pending or failed persistence MUST NOT be represented as persisted success. Other settings behavior SHALL remain unchanged.

#### Scenario: Change theme updates the actual interface
- **WHEN** the operator selects light from the theme select
- **THEN** state.settings.general.theme becomes light and the shared provider resolves light for the new recording surfaces and affected controls without changing recording state.

#### Scenario: Select automatic theme
- **WHEN** the operator selects auto and the OS preference later changes
- **THEN** the provider follows that change without persisting the resolved color as an explicit manual preference.

#### Scenario: Toggle audit log
- **WHEN** the operator toggles enable_audit_log from on to off
- **THEN** state.settings.general.enable_audit_log becomes false.
