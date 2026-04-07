## 1. Lock the umbrella contract

- [x] 1.1 Confirm the umbrella change remains the parent contract for `rule-driven-foundation`, `tag-review-first-flow`, `database-output-review-apply`, `local-modbus-review-apply`, and `studio-orchestration-cleanup`
- [x] 1.2 Verify the umbrella spec map stays aligned with the 10 affected capability surfaces named in the proposal and specs
- [x] 1.3 Verify the umbrella design still matches the approved workflow contract: `Device -> SourceRule -> Tag review -> Output review/apply`

## 2. Author the executable child changes

- [x] 2.1 Create `rule-driven-foundation` as the owner of `source-rule-runtime`, `device-registry`, `protocol-connectors`, shared `mapping-pipeline` diff rules, and shared candidate/recompute API contracts
- [x] 2.2 Create `tag-review-first-flow` as the owner of `tag-dictionary`, Step 3 review/apply behavior, and tag-specific apply API contracts
- [x] 2.3 Create `database-output-review-apply` as the owner of database output review/apply behavior and database-specific apply API contracts, explicitly as a sibling of Local Modbus rather than its prerequisite
- [x] 2.4 Create `local-modbus-review-apply` as the owner of Local Modbus review/apply behavior and Local Modbus-specific apply API contracts, explicitly as a sibling of Database rather than its prerequisite
- [x] 2.5 Create `studio-orchestration-cleanup` as the owner of `/studio` consolidation, redirect cleanup, and removal of manual-first product assumptions

## 3. Enforce phase ordering and dependency gates

- [x] 3.1 Define the unlock rule that no child change may start implementation until its prerequisite child change has complete artifacts and an accepted inherited contract
- [x] 3.2 Require every child change proposal and design to name its prerequisite change ids and inherited umbrella contracts explicitly
- [x] 3.3 Require `database-output-review-apply` and `local-modbus-review-apply` to inherit the delivered tag review/apply contract instead of redefining candidate semantics
- [x] 3.4 Require `studio-orchestration-cleanup` to start only after earlier child changes have delivered a fully usable rule-driven replacement path

## 4. Govern execution against the umbrella change

- [x] 4.1 Require each child change to reference the umbrella workflow, ownership, and revision rules instead of inventing a parallel orchestration model
- [x] 4.2 Require each child change to define its local rollout and rollback boundary while preserving the umbrella target state
- [x] 4.3 Validate that `/studio` cleanup removes manual-first primary flow assumptions only where a child change has already supplied the rule-driven replacement
