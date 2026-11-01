package adapters

import "context"

// acquisitionGenerationState keeps delegated RTU/UDP reconnects on the owner.
// A proxy borrows a pointer, never copies the publication mutex or atomic state.
func (c *ModbusTCPConnector) acquisitionGenerationState() *acquisitionGeneration {
	if c.generationOwner != nil {
		return c.generationOwner
	}
	return &c.acquisitionGeneration
}

func (c *ModbusTCPConnector) connectAcquisitionContext(ctx context.Context) error {
	return c.acquisitionGenerationState().connect(ctx, c.client.ConnectContext)
}

// AcquisitionTransportGeneration is the last successful TCP connect publication.
// Reading it performs no I/O and does not acquire the transport operation gate.
func (c *ModbusTCPConnector) AcquisitionTransportGeneration() uint64 {
	return c.acquisitionGenerationState().published.Load()
}

// AcquisitionTransportGeneration includes successful delegated UDP reconnects.
func (c *ModbusUDPConnector) AcquisitionTransportGeneration() uint64 {
	return c.acquisitionGeneration.published.Load()
}

// AcquisitionTransportGeneration includes successful delegated RTU reconnects.
func (c *ModbusRTUConnector) AcquisitionTransportGeneration() uint64 {
	return c.acquisitionGeneration.published.Load()
}
