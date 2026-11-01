package point

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatReceiptKeyFixedWireGolden(t *testing.T) {
	for _, tc := range []struct{ name, workspace, request, want string }{
		{"empty", "", "", "point_format_receipt:v1:439083f38956ba51ece90631552c6ea23c5c29570d3d5710e408e77e01ba7375"},
		{"unicode", "工作區🙂", "請求", "point_format_receipt:v1:e30a806273d0a66e879334fddb7b7f5b0c6d08da7cf392df7a95f1b2be52a102"},
		{"escaped", "\"\\\n<>&", "\u0000\t\r", "point_format_receipt:v1:7c4b85175ce2473de95f7c9c378a81a732a5eba35f76c96a9cb56febfa39e220"},
	} {
		t.Run(tc.name, func(t *testing.T) { require.Equal(t, tc.want, formatReceiptKey(tc.workspace, tc.request)) })
	}
}
