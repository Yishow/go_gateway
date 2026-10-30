package mcprotocol

import "go-gateway/internal/protocol/wirebudget"

// contextMutex shares the cancellable operation gate with managed connections.
type contextMutex = wirebudget.Mutex
