package api

import (
	"github.com/stretchr/testify/require"
	"go-gateway/internal/datalink/sourcerule"
	"testing"
)

func TestRouterPointFormatInitializationAllowsMissingPointService(t *testing.T) {
	rules := sourcerule.NewService(sourcerule.NewMemoryRepository(), nil, nil, nil)
	require.NotPanics(t, func() { NewRouter(&DatalinkServices{SourceRule: rules}) })
}
