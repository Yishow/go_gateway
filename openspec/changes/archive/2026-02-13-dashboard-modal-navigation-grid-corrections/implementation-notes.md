# Implementation Notes (Exclusions)

This change explicitly excludes the following items per product decision:

1. Original recommendation #10: `conflict deep-link query` (not implemented)
2. Original recommendation #15: `recent write events` (not implemented)
3. Original recommendation #16: `Top3 conflict registers` (not implemented)

Applied scope focuses on:
- Sidebar-to-modal migration (except `/test`)
- Server Memory Grid discoverability
- Dashboard grid compaction and summary
- Typed occupancy correctness (`float32=2`, `int64=4`)
