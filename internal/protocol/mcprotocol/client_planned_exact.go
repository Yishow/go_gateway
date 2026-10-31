package mcprotocol

import (
	"encoding/binary"
	"fmt"
)

// exactReadTransport validates before the existing parser can truncate payloads
// or turn an invalid bit nibble into false. It preserves the underlying context
// transport (and its single shared attempt budget), serializer and end-code error.
type exactReadTransport struct {
	Transport
	count int
	bits  bool
}

func (t exactReadTransport) SendReceive(req []byte) ([]byte, error) {
	response, err := t.Transport.SendReceive(req)
	if err != nil {
		return nil, err
	}
	if len(response) < 2 {
		return nil, fmt.Errorf("incomplete MC response")
	}
	if binary.LittleEndian.Uint16(response) != 0 {
		return response, nil
	}
	n := t.count * 2
	if t.bits {
		n = (t.count + 1) / 2
	}
	if len(response) != 2+n {
		return nil, fmt.Errorf("invalid MC response length: expected %d, got %d", 2+n, len(response))
	}
	if t.bits {
		for i, b := range response[2:] {
			if b>>4 > 1 || b&15 > 1 || (i == n-1 && t.count%2 != 0 && b&15 != 0) {
				return nil, fmt.Errorf("invalid MC bit encoding")
			}
		}
	}
	return response, nil
}

// RandomReadExact uses the existing random serializer, address/count checks and
// status parser, rejecting empty, short or excess successful response payloads.
func (c *MCClient) RandomReadExact(items []RandomReadItem) ([]int, error) {
	view := *c
	view.transport = exactReadTransport{Transport: c.transport, count: len(items)}
	return view.RandomRead(items)
}

// BatchReadBitExact additionally requires only 0/1 nibbles and zero odd padding.
func (c *MCClient) BatchReadBitExact(device string, addr, count int) ([]bool, error) {
	view := *c
	view.transport = exactReadTransport{Transport: c.transport, count: count, bits: true}
	return view.BatchReadBit(device, addr, count)
}
