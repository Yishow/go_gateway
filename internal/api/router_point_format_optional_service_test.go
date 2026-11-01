package api

import (
	"go-gateway/internal/datalink/sourcerule"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouterPointFormatInitializationAllowsMissingPointService(t *testing.T) {
	rules := sourcerule.NewService(sourcerule.NewMemoryRepository(), nil, nil, nil)
	require.NotPanics(t, func() { NewRouter(&DatalinkServices{SourceRule: rules}) })
}
