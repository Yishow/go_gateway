package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/fatek"
)

// The only fake is the PLC wire peer. Production factories, clients and decoders
// remain installed; replies depend on the actual requested address and count.
type wordWirePLC struct {
	mu    sync.Mutex
	words map[int]uint16
	bad   bool
	calls int
}

func (p *wordWirePLC) set(start int, words []uint16) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i, w := range words {
		p.words[start+i] = w
	}
}
func (p *wordWirePLC) read(start, count int) ([]uint16, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.calls++
	words := make([]uint16, count)
	for i := range words {
		words[i] = p.words[start+i]
	}
	return words, p.bad
}
func wordWireServer(t *testing.T, protocol schema.ProtocolType) (*wordWirePLC, string) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	plc := &wordWirePLC{words: map[int]uint16{}}
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
			wg.Go(func() {
				defer conn.Close()
				var err error
				reader := bufio.NewReader(conn)
				for {
					var response []byte
					switch protocol {
					case schema.ProtocolModbusTCP:
						header := make([]byte, 7)
						if _, err = io.ReadFull(reader, header); err != nil {
							return
						}
						body := make([]byte, int(binary.BigEndian.Uint16(header[4:6]))-1)
						if _, err = io.ReadFull(reader, body); err != nil || len(body) != 5 {
							return
						}
						start, count := int(binary.BigEndian.Uint16(body[1:3])), int(binary.BigEndian.Uint16(body[3:5]))
						words, bad := plc.read(start, count)
						payload := make([]byte, 2+len(words)*2)
						payload[0], payload[1] = body[0], byte(len(words)*2)
						for i, w := range words {
							binary.BigEndian.PutUint16(payload[2+2*i:], w)
						}
						if bad {
							payload = payload[:len(payload)-1]
						}
						binary.BigEndian.PutUint16(header[4:6], uint16(len(payload)+1))
						response = append(header, payload...)
					case schema.ProtocolMC3E:
						header := make([]byte, 9)
						if _, err = io.ReadFull(reader, header); err != nil {
							return
						}
						body := make([]byte, binary.LittleEndian.Uint16(header[7:]))
						if _, err = io.ReadFull(reader, body); err != nil || len(body) != 12 {
							return
						}
						start := int(body[6]) | int(body[7])<<8 | int(body[8])<<16
						words, bad := plc.read(start, int(binary.LittleEndian.Uint16(body[10:])))
						payload := make([]byte, 2+2*len(words))
						for i, w := range words {
							binary.LittleEndian.PutUint16(payload[2+2*i:], w)
						}
						if bad {
							payload = payload[:len(payload)-1]
						}
						header[0], header[1] = 0xd0, 0
						binary.LittleEndian.PutUint16(header[7:], uint16(len(payload)))
						response = append(header, payload...)
					case schema.ProtocolFatekFBs:
						frame, e := reader.ReadBytes(fatek.ETX)
						if e != nil || len(frame) < 12 {
							return
						}
						command := string(frame[3:5])
						body := string(frame[5 : len(frame)-3])
						count, e := strconv.ParseInt(body[:2], 16, 32)
						if e != nil {
							return
						}
						address := body[2:]
						wide := strings.HasPrefix(address, "DR") || strings.HasPrefix(address, "DD")
						start, e := strconv.Atoi(strings.TrimLeft(address, "DR"))
						if e != nil {
							return
						}
						nativeCount := int(count)
						if wide {
							nativeCount *= 2
						}
						words, bad := plc.read(start, nativeCount)
						payload := "0"
						if wide {
							for i := 0; i < len(words); i += 2 {
								payload += fmt.Sprintf("%04X%04X", words[i+1], words[i])
							}
						} else {
							for _, w := range words {
								payload += fmt.Sprintf("%04X", w)
							}
						}
						if bad {
							payload = payload[:len(payload)-1]
						}
						response = fatek.BuildFrame(1, command, payload)
					}
					if _, err := conn.Write(response); err != nil {
						return
					}
				}
			})
		}
	})
	t.Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		for _, c := range connections {
			_ = c.Close()
		}
		mu.Unlock()
		wg.Wait()
	})
	port := listener.Addr().(*net.TCPAddr).Port
	return plc, fmt.Sprintf(`{"host":"127.0.0.1","port":%d,"mode":"tcp","station_no":1,"slave_id":1,"timeout":1,"retry_count":0,"data_format":"ABCD"}`, port)
}
