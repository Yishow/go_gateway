package sourcerule

import (
	"fmt"

	"go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func assignRuleRevisionID(rule *schema.SourceRule) error {
	if rule == nil {
		return fmt.Errorf("source rule is nil")
	}

	revisionID, err := common.NewUUID()
	if err != nil {
		return fmt.Errorf("建立來源規則 revision id 失敗: %w", err)
	}
	rule.RevisionID = revisionID
	return nil
}
