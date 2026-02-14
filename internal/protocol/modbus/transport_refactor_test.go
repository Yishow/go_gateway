package modbus

import (
	"testing"
	"time"

	"go.bug.st/serial"
)

func TestNewRTUTransportDefaults(t *testing.T) {
	transport := NewRTUTransport("COM1", 0, 0, 0, "x", 0)

	if transport.BaudRate != DefaultBaudRate {
		t.Fatalf("預設 BaudRate 錯誤: %d", transport.BaudRate)
	}
	if transport.DataBits != DefaultDataBits {
		t.Fatalf("預設 DataBits 錯誤: %d", transport.DataBits)
	}
	if transport.Timeout != DefaultTimeout*time.Second {
		t.Fatalf("預設 Timeout 錯誤: %v", transport.Timeout)
	}
	if transport.Parity != serial.NoParity {
		t.Fatalf("預設 Parity 錯誤: %v", transport.Parity)
	}
	if transport.StopBits != serial.OneStopBit {
		t.Fatalf("預設 StopBits 錯誤: %v", transport.StopBits)
	}
}

func TestNewRTUTransportParityAndStopBits(t *testing.T) {
	transport := NewRTUTransport("COM2", 19200, 7, 2, "E", 3*time.Second)

	if transport.Parity != serial.EvenParity {
		t.Fatalf("Parity 應為 EvenParity，實際: %v", transport.Parity)
	}
	if transport.StopBits != serial.TwoStopBits {
		t.Fatalf("StopBits 應為 TwoStopBits，實際: %v", transport.StopBits)
	}
	if transport.Timeout != 3*time.Second {
		t.Fatalf("Timeout 應為 3s，實際: %v", transport.Timeout)
	}
}
