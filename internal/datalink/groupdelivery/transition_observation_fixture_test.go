package groupdelivery

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Reflection keeps the initial real SQLite behavioral test executable before
// the additive optional observation API exists. No fake production outcome is used.
func observedDeliveryStore(t *testing.T, store *Store, observations *[]map[string]any) *Store {
	t.Helper()
	method := reflect.ValueOf(store).MethodByName("WithDeliveryObservation")
	if !method.IsValid() {
		return store
	}
	config := reflect.New(method.Type().In(0)).Elem()
	read := config.FieldByName("ScopeRevision")
	read.Set(reflect.MakeFunc(read.Type(), func([]reflect.Value) []reflect.Value {
		return []reflect.Value{reflect.ValueOf("observed-W1"), reflect.ValueOf(true)}
	}))
	observe := config.FieldByName("Observe")
	observe.Set(reflect.MakeFunc(observe.Type(), func(args []reflect.Value) []reflect.Value {
		raw, err := json.Marshal(args[0].Interface())
		require.NoError(t, err)
		var row map[string]any
		require.NoError(t, json.Unmarshal(raw, &row))
		*observations = append(*observations, row)
		return nil
	}))
	return method.Call([]reflect.Value{config})[0].Interface().(*Store)
}

func deliveryObservationRow(t *testing.T, store *Store, offset int) (OutboxItem, Claim) {
	t.Helper()
	effect := enqueueFor(t, store, testKey, testDestination, "entity", at(offset), uint64(offset+1))
	claim, err := store.BeginDelivery(t.Context(), effect, "node/observer-test", time.Minute)
	require.NoError(t, err)
	item, err := store.GetOutbox(t.Context(), effect)
	require.NoError(t, err)
	return item, claim
}
