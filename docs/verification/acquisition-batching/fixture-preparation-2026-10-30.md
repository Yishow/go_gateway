# Acquisition batching: fixture preparation (2026-10-30)

## Status and boundaries

Test/data preparation only for
[the batching change](../../../openspec/changes/integrate-protocol-aware-acquisition-batching/design.md).
No production implementation, task completion, PLC access, commit or push.
Source baseline: `6f4a5fc4a1faf09e3d874b438bc6e25e1f7f6b8c`.
The worktree also contains other workers' proposal/publication work; this is not
an assertion about a clean checkout or the final remote main revision.

**22 design vectors passed integrity checks, not planner acceptance.** The
fixture file explicitly labels its expected plans as unexecuted design outputs.
No new planner exists here; JSON does not establish optimized results.

## Fixture contract

Files are under `internal/datalink/optimizer/`:

- `testdata/acquisition_acceptance.json`: versioned vectors; intervals are
  half-open native-unit ranges. `repeat` expands `id` into `id0`, `id1`, etc.,
  using `stride` native units. An absent repeat means one unchanged identity.
- `acceptance_fixtures_test.go`: strict typed JSON schema (unknown fields fail),
  capability/limit checks, unique Point identities, complete scalar containment
  in exactly one expected request, safe read coverage, explicit point errors,
  selection references, and independent numeric decoder oracles.
- `acceptance_native_test.go`: checks original template addresses/offsets with
  existing MC/Modbus/FATEK parsers and native widths with existing schema/FATEK
  component definitions. It is not a replacement address parser.

Vectors cover nested overlap; adjacency; uncovered/explicitly readable/forbidden
Modbus gaps; malformed address; an oversized scalar; 125/126-unit limits;
uncut uint32 at the limit; ABCD/CDAB identities over the same words; all 21
acquisition Points despite three recording selections; exact manual subset;
MC/FATEK supported and unknown random capability; MC 256 / FATEK 65 sparse
native units; FATEK D/DR float32 widths; MC `W1A` hexadecimal offset 26.
The forbidden-gap case retains two requests: [0,1) and [10,21); the forbidden
cell 5 blocks only the first gap, not the authorized second gap.
Configured random limits are fixture conditions, not a claim every PLC supports
those commands. MC adapter Count counts typed values, not expanded native words.

Source anchors: `optimizer/merger.go`, `optimizer/benchmark_test.go`,
`connector/point_read.go`, `connector/adapters/mc3e_read.go` (all under
`internal/datalink/`); `internal/protocol/mcprotocol/client.go` RandomRead,
`internal/protocol/fatek/{address.go,address_parse.go,client.go}`,
`lib/hsllogic/address_parser_{mitsubishi,modbus}.go`.

## Current behavior versus required behavior

Two isolated Go overlay tests (not installed in the repository) asserted the
future contract against the legacy merger. Both failed as expected:

| Input | Observed legacy behavior | Required behavior |
| --- | --- | --- |
| D0 int64 + D1 int16 | One block, start 0, length **2**, offsets wide=0/inner=1 | Preserve union length **4** |
| `not-an-address` int16 | Zero blocks; no point-error result channel | Explicit invalid-address Point error |

These failures are defects, not golden expectations. Existing tests passing does
not disprove them. `DefaultMergerConfig.MaxHoleSize=10` also lacks readable-gap
permission; the old dispatcher has one global byte order. Neither helper should
be wired unchanged into production. Existing collector `scheduler_polling.go`
issues scalar `conn.Read` calls; these planner fixtures do not exercise that path.

## Measured legacy helper baseline

Linux amd64, Go 1.25.5, AMD EPYC 9V74; `GOMAXPROCS=2`, `-p 1`,
`-benchtime=100ms -count=3 -benchmem`. Existing `generatePoints` gives deterministic
mixed D/M addresses, ten-address groups separated by gaps. Benchmarks measure
only `BlockMerger.Merge`; no transport, simulated RTT, scheduler, or PLC I/O.
Benchmark setup is outside `b.Loop`; these are repeated in-process helper calls,
**not cold/warm connection measurements**. Shared-host timings are descriptive.

| Points | Sample | ns/op | B/op | allocs/op |
| ---: | ---: | ---: | ---: | ---: |
| 100 | 1 | 57058 | 68574 | 269 |
| 100 | 2 | 66937 | 68574 | 269 |
| 100 | 3 | 76382 | 68577 | 269 |
| 1000 | 1 | 722527 | 755861 | 2413 |
| 1000 | 2 | 743859 | 755862 | 2413 |
| 1000 | 3 | 669809 | 755857 | 2413 |
| 10000 | 1 | 15844672 | 8342778 | 23826 |
| 10000 | 2 | 13778832 | 8342663 | 23826 |
| 10000 | 3 | 14164817 | 8342746 | 23827 |

No optimized timing, speedup, protocol request count, TCP packet count, latency
percentile or wire-byte measurement is inferred from this helper benchmark.
Expected JSON request counts are initial protocol operations excluding retries;
they are neither measured requests nor TCP packets (segmentation/coalescing can
differ). Adapter request-count evidence is a separate worker's deliverable.

## Reproduction and checks

From the repo root, use the existing Go 1.25.5 binary and a writable cache:

```sh
export GOMAXPROCS=2
# Set GOCACHE and GOPATH to writable/cache-backed paths for your environment.
go test -p 1 ./internal/datalink/optimizer -count=1 -v
go test -p 1 ./internal/datalink/optimizer -run '^$' \
  -bench '^BenchmarkMerger_' -benchtime=100ms -count=3 -benchmem
go vet -p 1 ./internal/datalink/optimizer
git diff --check
bash scripts/check_file_lines.sh
```

Actual environment used the preinstalled `/tmp/gateway-go/go/bin/go`, existing
`/tmp/gateway-gopath`, and dedicated writable `/tmp/batching-fixture-gocache`.
The first attempt using Go's default home cache failed before running tests
(read-only filesystem); changing only the cache location resolved it.

PASS: original optimizer suite; final optimizer suite including 22 fixture
subtests; existing three-size benchmarks; focused vet; diff whitespace and line
limits. The isolated intended-contract probes are **RED**, intentionally absent
from committed tests. No production code or adapter test file was modified.

## Unmet acceptance (do not mark change tasks complete)

Need actual planner/executor comparison against these vectors, per-Point decoder
resolver integration, deterministic compatible-group/revision isolation,
unsupported/timeout/short/CRC/cancel behavior, bit-vector scatter, all scheduler
fairness/backpressure/lease tests, and scoped zero-I/O preview parity. Still need
the prescribed device/Point/RTT matrix, 100 scans, cold/warm connections, actual
attempts/bytes/retries/fallback, per-device p50/p95/freshness/overrun, peak queue
and goroutines, optimized CPU/allocations, and production-entrypoint simulation.
Focused adapter legacy and optimizer package race checks passed during independent
review; full backend/frontend suites and full-suite race, real device/field
acceptance were not run. No claim of completed batching or performance
improvement follows from this note.
