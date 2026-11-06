package dbtarget

import (
	"context"
	"errors"

	datalinkbase "go-gateway/internal/datalink"
	"go-gateway/internal/datalink/schema"
)

type destinationOpenResult struct {
	manager *datalinkbase.DBManager
	err     error
}

// openDestination bounds the wait for openers that do not accept a context.
// Only one opener may remain outstanding per connector; a late connection is
// closed rather than being handed to a worker after its deadline.
func (w *Writer) openDestination(ctx context.Context, connector *schema.DatabaseConnector, config ConnectionConfig) (*datalinkbase.DBManager, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	connectorID, kind := connector.ID, connector.Kind
	w.openMu.Lock()
	if w.opening == nil {
		w.opening = make(map[string]struct{})
	}
	if _, busy := w.opening[connectorID]; busy {
		w.openMu.Unlock()
		return nil, errors.New("previous destination connection attempt is still running")
	}
	w.opening[connectorID] = struct{}{}
	w.openMu.Unlock()

	open := w.open
	if open == nil {
		open = openExternalDBManagerFunc
	}
	result := make(chan destinationOpenResult)
	closeResult := make(chan bool)
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			w.openMu.Lock()
			delete(w.opening, connectorID)
			w.openMu.Unlock()
		}()
		manager, err := open(kind, config)
		select {
		case result <- destinationOpenResult{manager: manager, err: err}:
			if <-closeResult && manager != nil {
				_ = manager.Close()
			}
		case <-ctx.Done():
			if manager != nil {
				_ = manager.Close()
			}
		}
	}()

	select {
	case opened := <-result:
		if err := ctx.Err(); err != nil {
			closeResult <- true
			return nil, err
		}
		if opened.err != nil {
			closeResult <- true
			if opened.manager == nil {
				<-done
			}
			return nil, opened.err
		}
		if opened.manager == nil {
			closeResult <- false
			<-done
			return nil, errors.New("destination opener returned no connection")
		}
		closeResult <- false
		<-done
		return opened.manager, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
