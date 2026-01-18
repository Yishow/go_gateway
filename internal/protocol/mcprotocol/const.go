package mcprotocol

import (
	"fmt"
	"strings"
)

// Command Codes (3E Binary)
const (
	CmdBatchRead   = 0x0401
	CmdBatchWrite  = 0x1401
	CmdRandomRead  = 0x0403
	CmdRandomWrite = 0x1402 // Not implemented in MVP but good to have
)

// Subcommand Codes
const (
	SubCmdWord = 0x0000
	SubCmdBit  = 0x0001
)

// Device Types (Binary Codes)
type DeviceType struct {
	Name  string
	Code  byte
	IsBit bool
}

var (
	// Bit Devices
	DeviceM  = DeviceType{"M", 0x90, true}
	DeviceX  = DeviceType{"X", 0x9C, true}
	DeviceY  = DeviceType{"Y", 0x9D, true}
	DeviceL  = DeviceType{"L", 0x92, true}
	DeviceB  = DeviceType{"B", 0xA0, true}
	DeviceF  = DeviceType{"F", 0x93, true}
	DeviceSB = DeviceType{"SB", 0xA1, true}
	DeviceS  = DeviceType{"S", 0x98, true}  // Step Relay
	DeviceTS = DeviceType{"TS", 0xC1, true} // Timer Contact
	DeviceTC = DeviceType{"TC", 0xC0, true} // Timer Coil
	DeviceCS = DeviceType{"CS", 0xC4, true} // Counter Contact
	DeviceCC = DeviceType{"CC", 0xC3, true} // Counter Coil

	// Word Devices
	DeviceD  = DeviceType{"D", 0xA8, false}
	DeviceW  = DeviceType{"W", 0xB4, false}
	DeviceR  = DeviceType{"R", 0xAF, false}
	DeviceTN = DeviceType{"TN", 0xC2, false} // Timer Current Value
	DeviceCN = DeviceType{"CN", 0xC5, false} // Counter Current Value
	DeviceSW = DeviceType{"SW", 0xB5, false}
	DeviceZ  = DeviceType{"Z", 0xCC, false}
)

var deviceMap = map[string]DeviceType{
	"M": DeviceM, "X": DeviceX, "Y": DeviceY, "L": DeviceL, "B": DeviceB, "F": DeviceF,
	"SB": DeviceSB, "S": DeviceS, "TS": DeviceTS, "TC": DeviceTC, "CS": DeviceCS, "CC": DeviceCC,
	"D": DeviceD, "W": DeviceW, "R": DeviceR, "TN": DeviceTN, "CN": DeviceCN,
	"SW": DeviceSW, "Z": DeviceZ,
}

func GetDeviceType(name string) (DeviceType, error) {
	upper := strings.ToUpper(name)
	if d, ok := deviceMap[upper]; ok {
		return d, nil
	}
	return DeviceType{}, fmt.Errorf("unknown device type: %s", name)
}

// Errors
type MCError struct {
	Code int
	Msg  string
}

func (e *MCError) Error() string {
	return fmt.Sprintf("MC Protocol Error 0x%04X: %s", e.Code, e.Msg)
}
