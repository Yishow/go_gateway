# Iteration Rules

Use these rules when the generated UI is unsatisfactory.

## Change prompts first when

- the information architecture is wrong
- the page layout is wrong
- the workflow hierarchy is wrong
- the visual direction is wrong
- the operating tone is wrong

## Refine generated output when

- spacing is off
- panel order is slightly wrong
- labels need adjustment
- controls need local repositioning
- one screen needs polish without changing the system

## Stitch MCP operating rules

- Edit one screen at a time.
- Change one structural goal at a time.
- Prefer `edit_screens` for accepted directions and `generate_screen_from_text` for first-pass exploration.
- If a Stitch edit is slow, wait rather than retrying immediately.
- If the user interrupts a Stitch edit, inspect `get_screen` and `get_project` before sending another edit.
- Do not claim visual success when the latest screenshot has not been reviewed.

## Move to implementation patching only when

- the route map is accepted
- the screen map is accepted
- the generated layouts are directionally correct
- the backend reconnection path is already known
