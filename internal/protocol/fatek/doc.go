/*
Package fatek implements the FATEK FBs PLC ASCII Protocol.

It supports both TCP and Serial (via io.ReadWriteCloser) transports.

Key Features:
	- Strong typing for components (X, Y, D, R, DR, DD, DF, WX/WY/WM/WS/WT/WC, etc.)
  - Automatic LRC calculation and verification
  - Support for continuous reading/writing (Cmd 44-47)
  - Support for mixed/random reading (Cmd 48) with 16/32-bit support
  - Thread-safe client

Usage:

	transport := fatek.NewTCPTransport("192.168.1.5", 500)
	client := fatek.NewClient(transport, 1)

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Read D0-D9 (10 registers)
	data, err := client.ReadRegisters("D", 0, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("D0-D9: %v\n", data)
*/
package fatek
