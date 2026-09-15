package handlers

import (
	"fmt"

	"go-gateway/internal/protocol/modbus"
)

func (h *TestHandler) prepareOverrideClient(state *ConnectionState, unitID *byte, station *int, tempSuffix string) (clientToUse, tempClient interface{}, err error) {
	clientToUse = state.Client

	if unitID == nil && station == nil {
		return clientToUse, nil, nil
	}

	if unitID != nil {
		if modbusClient, ok := state.Client.(*modbus.ModbusClient); ok {
			return modbusClient.WithUnitID(*unitID), nil, nil
		}
	}

	tempConfig := make(map[string]interface{})
	for k, v := range state.Config {
		tempConfig[k] = v
	}
	if unitID != nil {
		tempConfig["unitID"] = int(*unitID)
	}
	if station != nil {
		tempConfig["station"] = *station
	}

	tempClient, err = h.createClientWithDebug(state.Protocol, tempConfig, state.ID+tempSuffix)
	if err != nil {
		return nil, nil, fmt.Errorf("創建臨時客戶端失敗: %w", err)
	}
	if err = h.connectClient(tempClient, state.Protocol); err != nil {
		if closeErr := h.closeClient(tempClient, state.Protocol); closeErr != nil {
			fmt.Printf("關閉臨時客戶端失敗: %v\n", closeErr)
		}
		return nil, nil, fmt.Errorf("連線臨時客戶端失敗: %w", err)
	}

	return tempClient, tempClient, nil
}

func (h *TestHandler) closeTempClient(client interface{}, protocol string) {
	if client == nil {
		return
	}
	if closeErr := h.closeClient(client, protocol); closeErr != nil {
		fmt.Printf("關閉臨時客戶端失敗: %v\n", closeErr)
	}
}
