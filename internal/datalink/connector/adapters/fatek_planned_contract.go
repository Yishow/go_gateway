package adapters

import (
	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/optimizer"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// Validate every scalar's physical addresses, not just a planner offset. This
// also makes a future planner bug fail closed instead of decoding another point.
func validateFatekPlannedSlice(r optimizer.AcquisitionRequest, p optimizer.PointSlice) error {
	symbol, start, err := fatek.ParseComponentAddress(p.Point.Address)
	if err != nil {
		return err
	}
	component, err := fatek.GetComponentType(symbol)
	if err != nil {
		return err
	}
	scalarBits := 16 * schema.RegisterCountForDataType(p.Point.DataType)
	if p.Point.DataType == schema.DataTypeBool {
		scalarBits = 1
	}
	if p.NativeBits != component.Width || p.Length != (scalarBits+component.Width-1)/component.Width || p.Length != p.Point.NativeWidth || p.Offset < 0 {
		return connector.ErrPlannedRead
	}
	step := component.AddressStep()
	if r.Kind == "random" {
		if r.AddressStep != 0 || p.Offset > len(r.Units)-p.Length {
			return connector.ErrPlannedRead
		}
		for i := range p.Length {
			unit := r.Units[p.Offset+i]
			if unit.Area != symbol || unit.Address != start+i*step || unit.Bits != component.Width {
				return connector.ErrPlannedRead
			}
		}
	} else if r.Kind != "contiguous" || r.AddressStep != step || r.Group.Area != symbol || p.Offset > r.Length-p.Length || r.Start+p.Offset*step != start {
		return connector.ErrPlannedRead
	}
	return nil
}
