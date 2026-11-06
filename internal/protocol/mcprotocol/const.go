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

	// Long Timer / Long Counter Devices (SH-080008; forbidden in 0403 random read)
	DeviceLTC  = DeviceType{"LTC", 0x50, true}   // Long Timer Coil
	DeviceLTS  = DeviceType{"LTS", 0x51, true}   // Long Timer Contact
	DeviceLTN  = DeviceType{"LTN", 0x52, false}  // Long Timer Current Value
	DeviceLCC  = DeviceType{"LCC", 0x54, true}   // Long Counter Coil
	DeviceLCS  = DeviceType{"LCS", 0x55, true}   // Long Counter Contact
	DeviceLCN  = DeviceType{"LCN", 0x56, false}  // Long Counter Current Value
	DeviceLSTC = DeviceType{"LSTC", 0x58, true}  // Long Retentive Timer Coil
	DeviceLSTS = DeviceType{"LSTS", 0x59, true}  // Long Retentive Timer Contact
	DeviceLSTN = DeviceType{"LSTN", 0x5A, false} // Long Retentive Timer Current Value
)

var deviceMap = map[string]DeviceType{
	"M": DeviceM, "X": DeviceX, "Y": DeviceY, "L": DeviceL, "B": DeviceB, "F": DeviceF,
	"SB": DeviceSB, "S": DeviceS, "TS": DeviceTS, "TC": DeviceTC, "CS": DeviceCS, "CC": DeviceCC,
	"D": DeviceD, "W": DeviceW, "R": DeviceR, "TN": DeviceTN, "CN": DeviceCN,
	"SW": DeviceSW, "Z": DeviceZ,
	"LTC": DeviceLTC, "LTS": DeviceLTS, "LTN": DeviceLTN,
	"LCC": DeviceLCC, "LCS": DeviceLCS, "LCN": DeviceLCN,
	"LSTC": DeviceLSTC, "LSTS": DeviceLSTS, "LSTN": DeviceLSTN,
}

// IsLongTimerOrCounter reports whether name is a Mitsubishi long timer or long counter device.
func IsLongTimerOrCounter(name string) bool {
	switch strings.ToUpper(name) {
	case "LTN", "LTS", "LTC", "LSTS", "LSTC", "LSTN", "LCS", "LCC", "LCN":
		return true
	default:
		return false
	}
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
