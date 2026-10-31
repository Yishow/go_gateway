package modbus

import (
	"context"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go-gateway/internal/protocol/wirebudget"
)

// Literals are the Modbus Organization V1.1b3 examples, §§6.1–6.4,
// printed pp.12,13,15,16–17; they do not come from the legacy client.
func TestPlannedExactOfficialWireExamples(t *testing.T) {
	for _, rtu := range []bool{false, true} {
		t.Run(fmt.Sprint("rtu=", rtu), func(t *testing.T) {
			p := plannedBasicPeer(rtu)
			coil := []bool{true, false, true, true, false, false, true, true, true, true, false, true, false, true, true, false, true, false, true}
			discrete := []bool{false, false, true, true, false, true, false, true, true, true, false, true, true, false, true, true, true, false, true, false, true, true}
			for i, v := range coil {
				p.bits[plannedCell{1, uint16(19 + i)}] = v
			}
			for i, v := range discrete {
				p.bits[plannedCell{2, uint16(196 + i)}] = v
			}
			p.words[plannedCell{3, 107}] = 555
			p.words[plannedCell{3, 108}] = 0
			p.words[plannedCell{3, 109}] = 100
			p.words[plannedCell{4, 8}] = 10
			c := plannedExactReader(t, NewClient(p, 1).WithContext(wirebudget.With(t.Context(), 4, 0)))
			a, e := c.ReadBitsExact(1, 19, 19)
			require.NoError(t, e)
			require.Equal(t, coil, a)
			b, e := c.ReadBitsExact(2, 196, 22)
			require.NoError(t, e)
			require.Equal(t, discrete, b)
			d, e := c.ReadWordsExact(3, 107, 3)
			require.NoError(t, e)
			require.Equal(t, []uint16{555, 0, 100}, d)
			f, e := c.ReadWordsExact(4, 8, 1)
			require.NoError(t, e)
			require.Equal(t, []uint16{10}, f)
			off := 7
			if rtu {
				off = 1
			}
			for i, want := range [][]byte{{1, 0, 19, 0, 19}, {2, 0, 196, 0, 22}, {3, 0, 107, 0, 3}, {4, 0, 8, 0, 1}} {
				require.Equal(t, want, p.frames[i][off:off+5])
			}
		})
	}
}

func TestPlannedExactRejectsMalformedReadPayload(t *testing.T) {
	cases := []struct {
		name     string
		fn       byte
		quantity uint16
		pdu      []byte
	}{
		{"empty", 3, 2, []byte{3}}, {"word-short", 3, 2, []byte{3, 2, 1, 2}},
		{"word-overlong", 3, 2, []byte{3, 6, 1, 2, 3, 4, 5, 6}}, {"word-extra-tail", 3, 2, []byte{3, 4, 1, 2, 3, 4, 0}},
		{"word-bytecount", 3, 2, []byte{3, 3, 1, 2, 3, 4}}, {"word-count-zero", 4, 2, []byte{4, 0}},
		{"bit-short", 1, 9, []byte{1, 1, 9}}, {"bit-overlong", 2, 9, []byte{2, 3, 9, 1, 0}},
		{"bit-extra-tail", 1, 9, []byte{1, 2, 9, 1, 0}}, {"bit-padding", 2, 9, []byte{2, 2, 9, 0x81}},
	}
	for _, rtu := range []bool{false, true} {
		for _, tc := range cases {
			t.Run(fmt.Sprint(rtu, "/", tc.name), func(t *testing.T) {
				p := plannedBasicPeer(rtu)
				p.mutate = func(r []byte) []byte { return plannedResponsePayload(r, rtu, tc.pdu) }
				c := plannedExactReader(t, NewClient(p, 1).WithContext(wirebudget.With(t.Context(), 1, 0)))
				var err error
				if tc.fn < 3 {
					var v []bool
					v, err = c.ReadBitsExact(tc.fn, 0, tc.quantity)
					require.Nil(t, v)
				} else {
					var v []uint16
					v, err = c.ReadWordsExact(tc.fn, 0, tc.quantity)
					require.Nil(t, v)
				}
				require.Error(t, err)
				require.False(t, IsCommandRejected(err))
				require.Len(t, p.frames, 1)
			})
		}
	}
}

func TestPlannedExactRejectsEnvelopeAndForeignExceptions(t *testing.T) {
	for _, rtu := range []bool{false, true} {
		cases := map[string]func([]byte) []byte{
			"empty": func([]byte) []byte { return nil },
			"short": func(r []byte) []byte { return r[:3] },
			"unit": func(r []byte) []byte {
				if rtu {
					r[0] = 2
					binary.LittleEndian.PutUint16(r[len(r)-2:], CalculateCRC16(r[:len(r)-2]))
				} else {
					r[6] = 2
				}
				return r
			},
			"function":           func(r []byte) []byte { return plannedResponsePayload(r, rtu, []byte{4, 4, 1, 2, 3, 4}) },
			"exception-function": func(r []byte) []byte { return plannedResponsePayload(r, rtu, []byte{0x84, 1}) },
			"exception-trailing": func(r []byte) []byte { return plannedResponsePayload(r, rtu, []byte{0x83, 1, 0}) },
			"exception-short":    func(r []byte) []byte { return plannedResponsePayload(r, rtu, []byte{0x83}) },
			"exception-unit": func(r []byte) []byte {
				r = plannedResponsePayload(r, rtu, []byte{0x83, 1})
				if rtu {
					r[0] = 2
					binary.LittleEndian.PutUint16(r[len(r)-2:], CalculateCRC16(r[:len(r)-2]))
				} else {
					r[6] = 2
				}
				return r
			},
		}
		if rtu {
			cases["crc"] = func(r []byte) []byte { r[len(r)-1] ^= 1; return r }
		} else {
			cases["transaction"] = func(r []byte) []byte { r[1]++; return r }
			cases["exception-transaction"] = func(r []byte) []byte { r = plannedResponsePayload(r, false, []byte{0x83, 1}); r[1]++; return r }
			cases["protocol"] = func(r []byte) []byte { r[3] = 1; return r }
			cases["mbap-short"] = func(r []byte) []byte { r[5]--; return r }
			cases["mbap-long"] = func(r []byte) []byte { r[5]++; return r }
		}
		for name, mutate := range cases {
			t.Run(fmt.Sprint(rtu, "/", name), func(t *testing.T) {
				p := plannedBasicPeer(rtu)
				p.mutate = mutate
				c := plannedExactReader(t, NewClient(p, 1))
				v, e := c.ReadWordsExact(3, 0, 2)
				require.Error(t, e)
				require.Nil(t, v)
				require.False(t, IsCommandRejected(e))
				require.Len(t, p.frames, 1)
			})
		}
	}
}

func TestPlannedExactExceptionAndBudgetAuthority(t *testing.T) {
	for _, rtu := range []bool{false, true} {
		t.Run(fmt.Sprint(rtu), func(t *testing.T) {
			p := plannedBasicPeer(rtu)
			p.mutate = func(r []byte) []byte { return plannedResponsePayload(r, rtu, []byte{0x83, 1}) }
			ctx := wirebudget.With(t.Context(), 1, 0)
			c := plannedExactReader(t, NewClient(p, 1).WithContext(ctx))
			v, e := c.ReadWordsExact(3, 0, 2)
			require.Nil(t, v)
			require.True(t, IsCommandRejected(e))
			require.Equal(t, 1, wirebudget.From(ctx).Used())
			_, e = c.ReadWordsExact(3, 0, 2)
			require.ErrorIs(t, e, wirebudget.ErrExhausted)
			require.Len(t, p.frames, 1)
			canceled, cancel := context.WithCancel(t.Context())
			cancel()
			_, e = plannedExactReader(t, NewClient(p, 1).WithContext(canceled)).ReadWordsExact(3, 0, 2)
			require.ErrorIs(t, e, context.Canceled)
			require.Len(t, p.frames, 1)
		})
	}
}

func TestPlannedExactRejectsBoundsBeforeSend(t *testing.T) {
	p := plannedBasicPeer(false)
	c := plannedExactReader(t, NewClient(p, 1))
	for _, tc := range []struct {
		fn   byte
		a, n uint16
		bits bool
	}{{3, 0, 0, false}, {3, 0, 126, false}, {4, 65535, 2, false}, {2, 0, 0, true}, {1, 0, 2001, true}, {1, 65535, 2, true}, {1, 0, 1, false}, {3, 0, 1, true}} {
		if tc.bits {
			v, e := c.ReadBitsExact(tc.fn, tc.a, tc.n)
			require.Error(t, e)
			require.Nil(t, v)
		} else {
			v, e := c.ReadWordsExact(tc.fn, tc.a, tc.n)
			require.Error(t, e)
			require.Nil(t, v)
		}
	}
	require.Empty(t, p.frames)
}
