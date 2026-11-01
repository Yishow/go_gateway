package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"

	"go-gateway/internal/datalink/schema"

	"github.com/stretchr/testify/require"
)

type acquisitionWireRequest struct {
	Kind      string
	Addresses []int
}
type acquisitionWireAction struct {
	wait      <-chan struct{}
	malformed bool
	reject    byte
}

// This peer parses real application frames; no production planner/serializer
// computes its requested addresses or replies. Missing memory is not Good zero.
type acquisitionWirePeer struct {
	mu       sync.Mutex
	words    map[int]uint16
	requests []acquisitionWireRequest
	action   func(acquisitionWireRequest) acquisitionWireAction
}

func (p *acquisitionWirePeer) log() []acquisitionWireRequest {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]acquisitionWireRequest(nil), p.requests...)
}
func (p *acquisitionWirePeer) setAction(action func(acquisitionWireRequest) acquisitionWireAction) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.action = action
}
func (p *acquisitionWirePeer) receive(request acquisitionWireRequest) (values []uint16, allPresent bool, responseAction acquisitionWireAction) {
	p.mu.Lock()
	p.requests = append(p.requests, request)
	words := make([]uint16, len(request.Addresses))
	complete := true
	for i, address := range request.Addresses {
		var ok bool
		words[i], ok = p.words[address]
		complete = complete && ok
	}
	action := p.action
	p.mu.Unlock()
	if action != nil {
		return words, complete, action(request)
	}
	return words, complete, acquisitionWireAction{}
}
func newAcquisitionWirePeer(t *testing.T, protocol schema.ProtocolType, words map[int]uint16) (wirePeer *acquisitionWirePeer, config string) {
	t.Helper()
	var listenConfig net.ListenConfig
	listener, err := listenConfig.Listen(t.Context(), "tcp", "127.0.0.1:0")
	require.NoError(t, err)
	peer := &acquisitionWirePeer{words: words}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var connections []net.Conn
	wg.Go(func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			connections = append(connections, conn)
			mu.Unlock()
			wg.Go(func() { defer conn.Close(); peer.serve(t.Context(), conn, protocol) })
		}
	})
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		for _, conn := range connections {
			_ = conn.Close()
		}
		mu.Unlock()
		wg.Wait()
	})
	port := listener.Addr().(*net.TCPAddr).Port
	station := 1
	if protocol == schema.ProtocolMC3E {
		station = 0
	}
	return peer, fmt.Sprintf(`{"host":"127.0.0.1","port":%d,"mode":"tcp","station_no":%d,"slave_id":1,"timeout":5,"retry_count":0,"data_format":"ABCD"}`, port, station)
}
func (p *acquisitionWirePeer) serve(ctx context.Context, conn net.Conn, protocol schema.ProtocolType) {
	reader := bufio.NewReader(conn)
	for {
		frame, err := readAcquisitionWireFrame(reader, protocol)
		if err != nil {
			return
		}
		request, err := parseAcquisitionWireRequest(frame, protocol)
		if err != nil {
			return
		}
		words, complete, action := p.receive(request)
		if action.wait != nil {
			select {
			case <-action.wait:
			case <-ctx.Done():
				return
			}
		}
		response := acquisitionWireResponse(frame, protocol, words, complete, action)
		if _, err = conn.Write(response); err != nil {
			return
		}
	}
}
func readAcquisitionWireFrame(reader *bufio.Reader, protocol schema.ProtocolType) ([]byte, error) {
	if protocol == schema.ProtocolFatekFBs {
		return reader.ReadBytes(3)
	}
	size := 9
	if protocol == schema.ProtocolModbusTCP {
		size = 7
	}
	header := make([]byte, size)
	if _, err := io.ReadFull(reader, header); err != nil {
		return nil, err
	}
	var length int
	if protocol == schema.ProtocolModbusTCP {
		length = int(binary.BigEndian.Uint16(header[4:6])) - 1
	} else {
		length = int(binary.LittleEndian.Uint16(header[7:]))
	}
	if length < 1 || length > 4096 {
		return nil, fmt.Errorf("invalid frame length")
	}
	body := make([]byte, length)
	if _, err := io.ReadFull(reader, body); err != nil {
		return nil, err
	}
	return append(header, body...), nil
}
func parseAcquisitionWireRequest(frame []byte, protocol schema.ProtocolType) (acquisitionWireRequest, error) {
	request := acquisitionWireRequest{Kind: "contiguous"}
	start, count := 0, 0
	switch protocol {
	case schema.ProtocolModbusTCP:
		if len(frame) != 12 || frame[7] != 3 {
			return request, fmt.Errorf("invalid fixture Modbus read")
		}
		start, count = int(binary.BigEndian.Uint16(frame[8:10])), int(binary.BigEndian.Uint16(frame[10:12]))
	case schema.ProtocolMC3E:
		if len(frame) < 17 || frame[0] != 0x50 || binary.LittleEndian.Uint16(frame[13:15]) != 0 {
			return request, fmt.Errorf("invalid fixture MC header")
		}
		command := binary.LittleEndian.Uint16(frame[11:13])
		if command == 0x0403 {
			request.Kind = "random"
			count = int(frame[15])
			if frame[16] != 0 || count < 1 || len(frame) != 17+4*count {
				return request, fmt.Errorf("invalid MC random count")
			}
			for i := range count {
				entry := frame[17+4*i : 21+4*i]
				if entry[3] != 0xa8 {
					return request, fmt.Errorf("unexpected MC device")
				}
				request.Addresses = append(request.Addresses, int(entry[0])|int(entry[1])<<8|int(entry[2])<<16)
			}
			return request, nil
		}
		if command != 0x0401 || len(frame) != 21 || frame[18] != 0xa8 {
			return request, fmt.Errorf("invalid MC contiguous request")
		}
		start, count = int(frame[15])|int(frame[16])<<8|int(frame[17])<<16, int(binary.LittleEndian.Uint16(frame[19:21]))
	case schema.ProtocolFatekFBs:
		if len(frame) < 14 || frame[0] != 2 || frame[len(frame)-1] != 3 {
			return request, fmt.Errorf("invalid FATEK framing")
		}
		var checksum byte
		for _, v := range frame[:len(frame)-3] {
			checksum += v
		}
		if string(frame[len(frame)-3:len(frame)-1]) != fmt.Sprintf("%02X", checksum) {
			return request, fmt.Errorf("invalid FATEK checksum")
		}
		n, err := strconv.ParseUint(string(frame[5:7]), 16, 8)
		if err != nil || n == 0 {
			return request, fmt.Errorf("invalid FATEK count")
		}
		count = int(n)
		body := string(frame[7 : len(frame)-3])
		command := string(frame[3:5])
		if command == "48" {
			request.Kind = "random"
			if len(body) != count*6 {
				return request, fmt.Errorf("invalid FATEK random length")
			}
			for i := range count {
				entry := body[i*6 : i*6+6]
				if entry[0] != 'D' {
					return request, fmt.Errorf("unexpected FATEK area")
				}
				address, err := strconv.Atoi(entry[1:])
				if err != nil {
					return request, err
				}
				request.Addresses = append(request.Addresses, address)
			}
			return request, nil
		}
		if command != "46" || len(body) != 6 || body[0] != 'D' {
			return request, fmt.Errorf("invalid FATEK contiguous request")
		}
		start, err = strconv.Atoi(body[1:])
		if err != nil {
			return request, err
		}
	}
	for i := range count {
		request.Addresses = append(request.Addresses, start+i)
	}
	return request, nil
}
func acquisitionWireResponse(frame []byte, protocol schema.ProtocolType, words []uint16, complete bool, action acquisitionWireAction) []byte {
	switch protocol {
	case schema.ProtocolFatekFBs:
		payload := "0"
		for _, word := range words {
			payload += fmt.Sprintf("%04X", word)
		}
		if !complete {
			payload = "2"
		}
		if action.reject != 0 {
			payload = string(action.reject)
		}
		if action.malformed && len(payload) > 1 {
			payload = payload[:len(payload)-1]
		}
		content := "\x02" + string(frame[1:5]) + payload
		var sum byte
		for i := range len(content) {
			sum += content[i]
		}
		return []byte(fmt.Sprintf("%s%02X\x03", content, sum))
	case schema.ProtocolMC3E:
		header := append([]byte(nil), frame[:9]...)
		header[0], header[1] = 0xd0, 0
		payload := make([]byte, 2+len(words)*2)
		for i, word := range words {
			binary.LittleEndian.PutUint16(payload[2+2*i:], word)
		}
		if !complete {
			payload = []byte{0x5b, 0xc0}
		}
		if action.reject != 0 {
			payload = []byte{0x59, 0xc0}
		}
		if action.malformed && len(payload) > 2 {
			payload = payload[:len(payload)-1]
		}
		binary.LittleEndian.PutUint16(header[7:], uint16(len(payload)))
		return append(header, payload...)
	default:
		header := append([]byte(nil), frame[:7]...)
		payload := make([]byte, 2+len(words)*2)
		payload[0], payload[1] = 3, byte(len(words)*2)
		for i, word := range words {
			binary.BigEndian.PutUint16(payload[2+2*i:], word)
		}
		if !complete {
			payload = []byte{0x83, 2}
		}
		if action.malformed && len(payload) > 2 {
			payload = payload[:len(payload)-1]
		}
		binary.BigEndian.PutUint16(header[4:6], uint16(len(payload)+1))
		return append(header, payload...)
	}
}
func withAcquisitionDeclaration(config, declaration string) string {
	return strings.TrimSuffix(config, "}") + `,"acquisition":` + declaration + `}`
}
