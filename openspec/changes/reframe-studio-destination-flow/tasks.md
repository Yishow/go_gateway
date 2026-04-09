# Tasks: Reframe Studio Destination Flow

## Phase 1: Spec and integration inventory

- [ ] Document current route map for `/studio` and `/test`
- [ ] Document current backend API groups used by `/studio`
- [ ] Document current backend API groups used by `/test`
- [ ] Document current frontend service/hook integration chain

## Phase 2: Workflow redesign

- [ ] Rewrite the `/studio` step model to `Device -> Source -> Tag -> Destination`
- [ ] Define `Destination Hub`
- [ ] Define `Share Hub`
- [ ] Define `Point`, `Tag`, `Tag Group`, and `Delivery Group` ownership

## Phase 3: Workspace redesign

- [ ] Define `Tag Workspace` as `Semantic Refinement Board`
- [ ] Define `Database Workspace` as row-planner-first
- [ ] Define `Local Modbus Workspace` as block-planner-first
- [ ] Define `MQTT Workspace` as message-planner-first

## Phase 4: `/test` redesign

- [ ] Define `/test` as `Field Engineer Debug Console`
- [ ] Preserve independent route ownership
- [ ] Inventory `/test` API dependencies and reconnect strategy

## Phase 5: AI generation package

- [ ] Produce Google Stitch screen-by-screen prompt package
- [ ] Produce Pencil MCP screen structure plan
- [ ] Produce backend reconnection guidance for AI-generated frontend output

## Phase 6: Implementation prep

- [ ] Identify screens that can reuse current hooks/services unchanged
- [ ] Identify screens that need adapter/view-model layer
- [ ] Identify gaps that truly require backend or OpenSpec amendments
