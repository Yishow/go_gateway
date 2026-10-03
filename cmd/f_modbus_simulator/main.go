// Command f_modbus_simulator is a loopback Modbus TCP device for the
// device-to-SQL acceptance run. It serves holding registers backed by a JSON
// file that the harness rewrites to change what the "device" reports, so no
// final API is stubbed and no SQL row is written by anything but the gateway.
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-gateway/internal/virtual/memory"
	modbusserver "go-gateway/internal/virtual/server/modbus"
)

func main() {
	port := flag.Int("port", 15020, "TCP port to listen on (loopback)")
	file := flag.String("registers", "", "JSON file: {\"<zero-based register>\": <uint16>, ...}; re-read every 500ms")
	flag.Parse()

	bank := memory.NewMemoryBank(1 << 17)
	server, err := startSimulator(bank, *port)
	if err != nil {
		log.Fatalf("start simulator: %v", err)
	}
	log.Printf("modbus simulator listening on %s", server.Address())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()
	var last []byte
	for {
		select {
		case <-stop:
			if err := server.Stop(); err != nil {
				log.Printf("stop simulator: %v", err)
			}
			return
		case <-tick.C:
			if *file == "" {
				continue
			}
			data, err := os.ReadFile(*file)
			if err != nil || bytes.Equal(data, last) {
				continue
			}
			var registers map[uint16]uint16
			if err := json.Unmarshal(data, &registers); err != nil {
				log.Printf("ignoring unreadable register file: %v", err)
				continue
			}
			for register, value := range registers {
				if err := bank.WriteWord(int(register)*2, value); err != nil {
					log.Printf("write register %d: %v", register, err)
				}
			}
			last = data
		}
	}
}

func startSimulator(bank *memory.MemoryBank, port int) (*modbusserver.Server, error) {
	server := modbusserver.NewServer(bank)
	err := server.StartWithConfig(context.Background(), modbusserver.Config{
		BindAddress:       "127.0.0.1",
		Port:              port,
		SlaveID:           1,
		CapacityRegisters: bank.Size() / 2,
	})
	return server, err
}
