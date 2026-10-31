package adapters

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"go-gateway/internal/protocol/fatek"
)

// This peer parses actual wire addresses independently of FormatAddress and the
// planner. Missing memory is an error, never an invented zero-valued response.
type fatekPlannedPeer struct {
	mu                sync.Mutex
	memory            map[string]any
	frames, responses [][]byte
	corrupt           func([]byte) []byte
	hook              func()
	err               error
}

func (*fatekPlannedPeer) Connect() error { return nil }
func (*fatekPlannedPeer) Close() error   { return nil }
func (p *fatekPlannedPeer) SendReceive(frame []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.frames = append(p.frames, append([]byte(nil), frame...))
	addresses, bits, err := parseFatekPlannedFrame(frame)
	if err != nil {
		return nil, err
	}
	var payload strings.Builder
	for i, address := range addresses {
		value, ok := p.memory[address]
		if !ok {
			return nil, fmt.Errorf("peer has no value at %s", address)
		}
		if bits[i] == 1 {
			v, ok := value.(bool)
			if !ok {
				return nil, fmt.Errorf("peer bit type at %s", address)
			}
			if v {
				payload.WriteByte('1')
			} else {
				payload.WriteByte('0')
			}
		} else {
			v, ok := value.(int)
			if !ok || v < 0 || uint64(v) > uint64(1)<<uint(bits[i])-1 {
				return nil, fmt.Errorf("peer register type/range at %s", address)
			}
			fmt.Fprintf(&payload, "%0*X", bits[i]/4, v)
		}
	}
	if p.hook != nil {
		p.hook()
	}
	if p.err != nil {
		return nil, p.err
	}
	response := fatekPeerFrame(string(frame[1:3]), string(frame[3:5]), "0"+payload.String())
	if p.corrupt != nil {
		response = p.corrupt(response)
	}
	p.responses = append(p.responses, response)
	return response, nil
}
func fatekPeerFrame(station, command, body string) []byte {
	content := "\x02" + station + command + body
	var sum byte
	for i := range len(content) {
		sum += content[i]
	}
	return []byte(fmt.Sprintf("%s%02X\x03", content, sum))
}
func parseFatekPlannedFrame(frame []byte) ([]string, []int, error) {
	if len(frame) < 10 || frame[0] != 2 || frame[len(frame)-1] != 3 {
		return nil, nil, fmt.Errorf("peer invalid framing")
	}
	var sum byte
	for _, b := range frame[:len(frame)-3] {
		sum += b
	}
	if string(frame[len(frame)-3:len(frame)-1]) != fmt.Sprintf("%02X", sum) {
		return nil, nil, fmt.Errorf("peer invalid checksum")
	}
	count, err := strconv.ParseUint(string(frame[5:7]), 16, 8)
	if err != nil || count < 1 {
		return nil, nil, fmt.Errorf("peer invalid count")
	}
	body := string(frame[7 : len(frame)-3])
	command := string(frame[3:5])
	var addresses []string
	var widths []int
	if command == "48" {
		cost := 0
		for range int(count) {
			symbol, address, width, used, err := parseFatekPeerAddress(body)
			if err != nil {
				return nil, nil, err
			}
			addresses = append(addresses, fmt.Sprint(symbol, address))
			widths = append(widths, width)
			cost += max(1, width/16)
			body = body[used:]
		}
		if body != "" || cost > 64 {
			return nil, nil, fmt.Errorf("peer invalid random length/cost")
		}
	} else if command == "44" || command == "46" {
		symbol, address, width, used, err := parseFatekPeerAddress(body)
		if err != nil {
			return nil, nil, err
		}
		if used != len(body) || (command == "44") != (width == 1) {
			return nil, nil, fmt.Errorf("peer wrong continuous type/length")
		}
		if width == 32 && count > 32 || width == 16 && count > 64 {
			return nil, nil, fmt.Errorf("peer continuous limit")
		}
		step := 1
		if symbol == "DR" || symbol == "DD" {
			step = 2
		}
		if strings.HasPrefix(symbol, "W") {
			step = 16
		}
		if strings.HasPrefix(symbol, "DW") {
			step = 32
		}
		for i := range int(count) {
			addresses = append(addresses, fmt.Sprint(symbol, address+i*step))
			widths = append(widths, width)
		}
	} else {
		return nil, nil, fmt.Errorf("peer unexpected command %s", command)
	}
	return addresses, widths, nil
}
func parseFatekPeerAddress(body string) (string, int, int, int, error) {
	// FATEK Appendix 1 pp3,14: discrete=5, register16=6, register32=7 characters.
	for _, c := range []struct {
		symbol        string
		width, length int
	}{{"DWM", 32, 7}, {"DWX", 32, 7}, {"DWY", 32, 7}, {"DWS", 32, 7}, {"DR", 32, 7}, {"DD", 32, 7}, {"WM", 16, 6}, {"WX", 16, 6}, {"WY", 16, 6}, {"WS", 16, 6}, {"D", 16, 6}, {"R", 16, 6}, {"X", 1, 5}, {"Y", 1, 5}, {"M", 1, 5}} {
		if !strings.HasPrefix(body, c.symbol) {
			continue
		}
		if len(body) < c.length {
			return "", 0, 0, 0, fmt.Errorf("peer short address")
		}
		digits := body[len(c.symbol):c.length]
		for _, b := range digits {
			if b < '0' || b > '9' {
				return "", 0, 0, 0, fmt.Errorf("peer nondecimal address")
			}
		}
		n, err := strconv.Atoi(digits)
		if err != nil {
			return "", 0, 0, 0, err
		}
		return c.symbol, n, c.width, c.length, nil
	}
	return "", 0, 0, 0, fmt.Errorf("peer unsupported address %s", body)
}
func fatekPlannedConnector(peer *fatekPlannedPeer) *FatekConnector {
	return &FatekConnector{client: fatek.NewClient(peer, 1), connected: true, persistentMode: true}
}
