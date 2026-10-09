package adapters

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"strings"
	"testing"

	"go-gateway/internal/datalink/connector"
	"go-gateway/internal/datalink/schema"
	"go-gateway/internal/protocol/mcprotocol"
	"go-gateway/lib/hsllogic"
)

// flakyMCTransport 前 fails 次回傳 EOF（模擬 FX5U「接受連線後直接關閉」的瞬時佔用），之後回正常回應
type flakyMCTransport struct {
	connects int
	fails    int
}

func (t *flakyMCTransport) Connect() error {
	t.connects++
	return nil
}

func (t *flakyMCTransport) Close() error {
	return nil
}

func (t *flakyMCTransport) SendReceive(_ []byte) ([]byte, error) {
	if t.fails > 0 {
		t.fails--
		return nil, io.EOF
	}
	// 回應 body：EndCode(2)=0 + 1 word = 0x1234
	resp := make([]byte, 4)
	binary.LittleEndian.PutUint16(resp[2:], 0x1234)
	return resp, nil
}

func newRetryTestConnector(fails int) *MC3EConnector {
	return &MC3EConnector{
		client:         mcprotocol.NewClientWithTransport(&flakyMCTransport{fails: fails}),
		connected:      true,
		persistentMode: true,
		dataConverter:  hsllogic.NewDataConverter(hsllogic.DataFormatCDAB),
	}
}

func TestMC3EConnector_ReadRetriesAfterEOF(t *testing.T) {
	c := newRetryTestConnector(1)

	res, err := c.Read(context.Background(), connector.ReadRequest{Address: "D100", Count: 1})
	if err != nil {
		t.Fatalf("Read 應在 EOF 後重連重試成功: %v", err)
	}
	if res.Quality != schema.QualityGood {
		t.Fatalf("品質應為 good, got %v", res.Quality)
	}
}

func TestMC3EConnector_ReadReportsPersistentEOF(t *testing.T) {
	c := newRetryTestConnector(100)

	_, err := c.Read(context.Background(), connector.ReadRequest{Address: "D100", Count: 1})
	if err == nil {
		t.Fatal("持續 EOF 應回報錯誤")
	}
	if !strings.Contains(err.Error(), "port") {
		t.Fatalf("錯誤訊息應提示 port 佔用/更換 port，got: %v", err)
	}
	if !errors.Is(err, io.EOF) {
		t.Fatalf("錯誤應保留 EOF 類別，got: %v", err)
	}
}

func TestMC3EConnector_TestConnectionRetriesAfterEOF(t *testing.T) {
	c := newRetryTestConnector(2)

	if err := c.TestConnection(context.Background()); err != nil {
		t.Fatalf("TestConnection 應在 EOF 後重連重試成功: %v", err)
	}
}

func TestMC3EConnector_NonConnectionErrorNotRetried(t *testing.T) {
	c := newRetryTestConnector(0)
	trans := c.client
	_ = trans

	// 位址解析錯誤不屬於連線錯誤，不應重試（也應立即失敗）
	if _, err := c.Read(context.Background(), connector.ReadRequest{Address: "ZZ9999", Count: 1}); err == nil {
		t.Fatal("無效位址應回報錯誤")
	}
}
