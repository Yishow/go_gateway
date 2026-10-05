package runtime

import (
	"encoding/json"
	"strconv"
)

// MarshalJSON preserves integer digits at the browser boundary. Values outside
// JavaScript's safe integer range travel as decimal strings; native Go event
// subscribers, typed samples and SQL delivery retain their original types.
func (event ValueEvent) MarshalJSON() ([]byte, error) {
	type wireEvent ValueEvent
	wire := wireEvent(event)
	wire.RawValue = browserInteger(event.RawValue)
	wire.TransformedValue = browserInteger(event.TransformedValue)
	return json.Marshal(wire)
}

func browserInteger(value any) any {
	const safeInteger int64 = 1<<53 - 1
	switch number := value.(type) {
	case int:
		if int64(number) > safeInteger || int64(number) < -safeInteger {
			return strconv.FormatInt(int64(number), 10)
		}
	case int64:
		if number > safeInteger || number < -safeInteger {
			return strconv.FormatInt(number, 10)
		}
	case uint:
		if uint64(number) > uint64(safeInteger) {
			return strconv.FormatUint(uint64(number), 10)
		}
	case uint64:
		if number > uint64(safeInteger) {
			return strconv.FormatUint(number, 10)
		}
	}
	return value
}
