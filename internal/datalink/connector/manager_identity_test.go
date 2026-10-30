package connector

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/schema"
)

func TestGenerationSerialPhysicalAliases(t *testing.T) {
	dir := t.TempDir()
	port := filepath.Join(dir, "port")
	require.NoError(t, os.WriteFile(port, nil, 0600))
	alias := filepath.Join(dir, "alias")
	require.NoError(t, os.Symlink(port, alias))
	for _, tt := range []struct {
		name, a, b string
		typ        schema.ProtocolType
	}{
		{"case mode", `{"mode":"SERIAL","serial_port":"COM9"}`, `{"mode":"serial","serial_port":"COM9"}`, schema.ProtocolFatekFBs},
		{"windows prefix", `{"serial_port":"COM9"}`, `{"serial_port":"\\\\.\\com9"}`, schema.ProtocolModbusRTU},
		{"symlink", `{"serial_port":"` + port + `"}`, `{"serial_port":"` + alias + `"}`, schema.ProtocolModbusRTU},
	} {
		t.Run(tt.name, func(t *testing.T) {
			registerGeneration(t, tt.typ, func() Protocol { return &generationProtocol{} })
			cm := NewConnectionManager(DefaultConnectionManagerConfig())
			_, e := cm.GetOrCreate(t.Context(), "a", tt.typ, tt.a)
			require.NoError(t, e)
			_, e = cm.GetOrCreate(t.Context(), "b", tt.typ, tt.b)
			require.Error(t, e)
		})
	}
}

func TestGenerationConfigCanonicalNestedAndInvalid(t *testing.T) {
	typ := schema.ProtocolType(t.Name())
	registerGeneration(t, typ, func() Protocol { return &generationProtocol{} })
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	a, e := cm.GetOrCreate(t.Context(), "a", typ, `{"auth":{"user":"x","password":" secret "},"topics":["one","two"],"id":9007199254740992}`)
	require.NoError(t, e)
	b, e := cm.GetOrCreate(t.Context(), "a", typ, `{"id":9007199254740992,"topics":["one","two"],"auth":{"password":" secret ","user":"x"}}`)
	require.NoError(t, e)
	require.Same(t, a, b)
	c, e := cm.GetOrCreate(t.Context(), "a", typ, `{"auth":{"user":"x","password":" secret "},"topics":["one","two"],"id":9007199254740993}`)
	require.NoError(t, e)
	require.NotSame(t, a, c)
	for _, bad := range []string{"null", "[]", "{} {}", `{"password":"secret"`} {
		_, e := cm.GetOrCreate(t.Context(), "a", typ, bad)
		require.Error(t, e)
		require.NotContains(t, e.Error(), "secret")
	}
	current, e := cm.GetOrCreate(t.Context(), "a", typ, c.Config)
	require.NoError(t, e)
	require.Same(t, c, current)
}

type generationCloseFailure struct{ generationProtocol }

func (p *generationCloseFailure) Close() error {
	p.closed.Add(1)
	return errors.New("private-port secret-credential")
}
func TestGenerationSerialCleanupFailureKeepsOwnership(t *testing.T) {
	for _, dialFail := range []bool{false, true} {
		t.Run(fmt.Sprint(dialFail), func(t *testing.T) {
			typ := schema.ProtocolModbusRTU
			p := &generationCloseFailure{}
			if dialFail {
				p.dial = func(context.Context, string) error { return errors.New("private-port secret-credential") }
			}
			registerGeneration(t, typ, func() Protocol { return p })
			cm := NewConnectionManager(DefaultConnectionManagerConfig())
			_, err := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9"}`)
			if !dialFail {
				require.NoError(t, err)
				err = cm.Close("a")
			}
			require.Error(t, err)
			require.NotContains(t, err.Error(), "private-port")
			require.NotContains(t, err.Error(), "secret-credential")
			_, err = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
			require.Error(t, err)
			require.EqualValues(t, 1, p.closed.Load())
		})
	}
}

func TestGenerationSerialReplacementAbortsIfOldCloseFails(t *testing.T) {
	typ := schema.ProtocolModbusRTU
	var calls int
	p := &generationCloseFailure{}
	registerGeneration(t, typ, func() Protocol { calls++; return p })
	cm := NewConnectionManager(DefaultConnectionManagerConfig())
	_, err := cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9","baud_rate":9600}`)
	require.NoError(t, err)
	_, err = cm.GetOrCreate(t.Context(), "a", typ, `{"serial_port":"COM9","baud_rate":19200}`)
	require.Error(t, err)
	require.Equal(t, 1, calls)
	_, err = cm.GetOrCreate(t.Context(), "b", typ, `{"serial_port":"COM9"}`)
	require.Error(t, err)
	require.Equal(t, 1, calls)
}
