/*
Package mcprotocol implements the Mitsubishi MELSEC Communication Protocol (MC Protocol) 3E Frame (Binary).

Supported Devices: D, W, R, M, X, Y, etc.
Supported Commands: Batch Read/Write (Word/Bit), Random Read.

Usage:

	client := mcprotocol.NewClient("192.168.1.10", 5000)
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Read D100-D104
	data, err := client.BatchReadWord("D", 100, 5)
	
	// Write M0-M3
	err = client.BatchWriteBit("M", 0, []bool{true, false, true, false})

*/
package mcprotocol
