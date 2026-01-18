# Code Review Report

## Overview
Review of the last 3 commits covering `datalink`, `scheduler`, and `connector` modules.

**Commits Reviewed:**
- `baea9b7` (HEAD) - `mapping` & `timeseries` implementation
- `50c056a` - `scheduler` & logic
- `16fa813` - `connector` & `schema` setup

## Critical Bugs

### 1. Scheduler Deadlock Risk (Blocking Channel)
**File:** [internal/datalink/collector/scheduler.go](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go)
**Location:** Function [emitValue](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#448-461), Line 450
**Issue:**
The [emitValue](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#448-461) function attempts to write to `valueChan`. If full, it reads one item and attempts to write again. However, the second write is **blocking** if the channel fills up immediately after the read (race condition with other goroutines).
Since [emitValue](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#448-461) is called by [pollPoint](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#405-447), which holds the **Device Lock** (`s.deviceLocks[deviceID]`), a blocked [emitValue](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#448-461) will hang the entire device polling thread. If the consumer of `valueChan` is slow or stopped, all polling threads will eventually hang, causing a total system deadlock.

**Fix Recommendation:**
Use a non-blocking select for the retry, or simply drop the new value if the channel is persistently full.
```go
func (s *Scheduler) emitValue(cv CollectedValue) {
    select {
    case s.valueChan <- cv:
    default:
        // Try to free space
        select {
        case <-s.valueChan:
            // Space freed, try sending again (non-blocking)
            select {
            case s.valueChan <- cv:
            default:
                // Still full, drop data to prevent blocking
                fmt.Println("Warning: Value buffer full, dropping data")
            }
        default:
            // Channel was empty but default triggered? Rare but handle safely
        }
    }
}
```

### 2. Incorrect Register Count for Multi-Word Types
**File:** [internal/datalink/collector/scheduler.go](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go) & [internal/datalink/connector/adapters/fatek.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go)
**Location:** [scheduler.go](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go) Line 411 (`Count: 1`)
**Issue:**
The scheduler always sets `Count: 1` in `ReadRequest`, regardless of the [DataType](file:///c:/AIProject/go_gateway/internal/datalink/schema/models.go#228-229) defined in the Point.
- For 16-bit registers (Modbus/Fatek "D", "R"), `Count: 1` reads 2 bytes (1 word).
- If [DataType](file:///c:/AIProject/go_gateway/internal/datalink/schema/models.go#228-229) is [int32](file:///c:/AIProject/go_gateway/internal/datalink/mapping/service.go#555-562) or `float32`, 4 bytes (2 words) are required.
- If [DataType](file:///c:/AIProject/go_gateway/internal/datalink/schema/models.go#228-229) is `int64` or `float64`, 8 bytes (4 words) are required.
This causes [FatekConnector](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go#21-27) (and likely Modbus) to read only the first word, resulting in incorrect values or data corruption for 32/64-bit types.

**Fix Recommendation:**
Implement a helper `GetRegisterCount(DataType)` in `schema` or `connector` package, and set `req.Count` dynamically in [pollPoint](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#405-447).

```go
// Helper example
func GetRegisterCount(dt schema.DataType) int {
    switch dt {
    case schema.DataTypeInt32, schema.DataTypeUint32, schema.DataTypeFloat32:
        return 2
    case schema.DataTypeInt64, schema.DataTypeUint64, schema.DataTypeFloat64:
        return 4
    default:
        return 1
    }
}
```

### 3. Weak ID Generation
**File:** [internal/datalink/mapping/service.go](file:///c:/AIProject/go_gateway/internal/datalink/mapping/service.go)
**Location:** [generateUUID](file:///c:/AIProject/go_gateway/internal/datalink/mapping/service.go#576-579), Line 576
**Issue:**
```go
func generateUUID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```
Using `UnixNano` as an ID is not safe for a "UUID". It will cause ID collisions if two mappings are created within the same nanosecond (possible in batch imports or concurrent tests) and is predictable.

**Fix Recommendation:**
Use `crypto/rand` or a standard library like `github.com/google/uuid` to generate proper UUIDs (v4).

### 4. FATEK Endianness Logic
**File:** [internal/datalink/connector/adapters/fatek.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go)
**Location:** [convertFatekValue](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go#255-291), Line 272
**Issue:**
```go
return int32(values[0])<<16 | int32(values[1])
```
This assumes `values[0]` (Address N) is the High Word and `values[1]` (Address N+1) is the Low Word (Big-Endian Words).
PLC memory mapping often puts the Lower Word at the lower address (N) and Higher Word at (N+1) (Little-Endian Words).
If FATEK follows Little-Endian Words (Double Word = D0 + size 2 -> D0 is low, D1 is high), this logic is inverted.

**Fix Recommendation:**
Verify FATEK specification. If it is Little-Endian Words (commonly "CDAB" or "BADC" ordering depending on byte endianness too), swap the indices:
```go
// If Low Word First:
return int32(values[1])<<16 | int32(values[0])
```

## Minor Issues
1.  **Unsafe Type Assertions**: In [mapping/service.go](file:///c:/AIProject/go_gateway/internal/datalink/mapping/service.go), lines like `format, _ := params["format"].(string)` will result in empty strings if the key is missing or wrong type, bypassing logic silently.
2.  **Hardcoded Configurations**: [SchedulerConfig](file:///c:/AIProject/go_gateway/internal/datalink/collector/scheduler.go#52-65) (Line 71) sets `MaxConcurrentPerDevice: 1`. While safe for Serial, TCP devices could support concurrency. (Not a bug, but a limitation).
3.  **Error Handling**: [fatek.go](file:///c:/AIProject/go_gateway/internal/datalink/connector/adapters/fatek.go) Line 244 converts `addr` using `Atoi` but doesn't check if the resulting address is within valid range for that symbol.

## Summary
The codebase is structured well ("Clean Architecture"), but the `scheduler` contains a critical concurrency bug and a data integrity bug ([Count](file:///c:/AIProject/go_gateway/internal/datalink/storage/timeseries.go#384-390)) that must be fixed before deploying even a prototype.
