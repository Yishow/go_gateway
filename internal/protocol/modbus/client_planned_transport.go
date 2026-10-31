package modbus

import (
	"context"
	"encoding/binary"
)

// exactReadTransport validates the envelope before the existing parser can turn
// an exception into unsupported-command evidence. Legacy scalar parsing is not
// changed. All I/O (and actual-send budget accounting) stays in Transport.
type exactReadTransport struct{ Transport }

func (t *exactReadTransport) GetOriginalTransport() any {
	if u, ok := t.Transport.(transportUnwrapper); ok {
		return u.GetOriginalTransport()
	}
	return t.Transport
}
func (t *exactReadTransport) GetNextTransactionID() uint16 { id, _ := t.nextTransactionID(); return id }
func (t *exactReadTransport) nextTransactionID() (uint16, error) {
	if v, ok := t.Transport.(interface{ nextTransactionID() (uint16, error) }); ok {
		return v.nextTransactionID()
	}
	if v, ok := t.Transport.(interface{ GetNextTransactionID() uint16 }); ok {
		return v.GetNextTransactionID(), nil
	}
	return 0, ErrInvalidFrame
}
func (t *exactReadTransport) SendReceive(request []byte) ([]byte, error) {
	response, err := t.sendExact(request)
	if err != nil {
		return nil, err
	}
	if _, rtu := t.GetOriginalTransport().(*RTUTransport); rtu {
		err = validateExactRTUEnvelope(request, response)
	} else {
		err = validateExactTCPEnvelope(request, response)
	}
	if err != nil {
		return nil, err
	}
	return response, nil
}
func validateExactTCPEnvelope(request, response []byte) error {
	header, err := ParseMBAPHeader(response)
	if err != nil {
		return err
	}
	if header.Length < 2 || header.Length > 254 || len(response) != 6+int(header.Length) {
		return ErrInvalidFrame
	}
	if header.TransactionID != binary.BigEndian.Uint16(request[:2]) {
		return ErrTransactionIDMismatch
	}
	if header.UnitID != request[6] {
		return ErrInvalidResponse
	}
	return validateExactResponseFunction(request[7], response[7:])
}
func validateExactRTUEnvelope(request, response []byte) error {
	if len(response) < RTUFrameMinLength || len(response) > 256 {
		return ErrInvalidFrame
	}
	if CalculateCRC16(response[:len(response)-2]) != binary.LittleEndian.Uint16(response[len(response)-2:]) {
		return ErrCRCError
	}
	if response[0] != request[0] {
		return ErrInvalidResponse
	}
	return validateExactResponseFunction(request[1], response[1:len(response)-2])
}
func validateExactResponseFunction(function byte, pdu []byte) error {
	if len(pdu) == 0 {
		return ErrResponseTooShort
	}
	if pdu[0] == function|FuncExceptionOffset {
		if len(pdu) != 2 {
			return ErrInvalidFrame
		}
		return nil
	}
	if pdu[0] != function {
		return ErrInvalidFunctionCode
	}
	return nil
}

// An RTU wrapper must explicitly preserve whole-frame acquisition. Never bypass
// it by unwrapping and sending directly on a different transport authority.
func (t *exactReadTransport) sendExact(request []byte) ([]byte, error) {
	if exact, ok := t.Transport.(interface{ SendReceiveExact([]byte) ([]byte, error) }); ok {
		return exact.SendReceiveExact(request)
	}
	if _, rtu := t.GetOriginalTransport().(*RTUTransport); rtu {
		return nil, ErrInvalidFrame
	}
	return t.Transport.SendReceive(request)
}

func (t *contextTransport) SendReceiveExact(request []byte) ([]byte, error) {
	if err := t.ctx.Err(); err != nil {
		return nil, err
	}
	if exact, ok := t.Transport.(interface {
		SendReceiveExactContext(context.Context, []byte) ([]byte, error)
	}); ok {
		return exact.SendReceiveExactContext(t.ctx, request)
	}
	if _, rtu := t.GetOriginalTransport().(*RTUTransport); rtu {
		return nil, ErrInvalidFrame
	}
	return t.SendReceive(request)
}
