package sourcerule

import (
	"fmt"

	common "go-gateway/internal/datalink/common"
	"go-gateway/internal/datalink/schema"
)

func assignRuleRevisionID(rule *schema.SourceRule) error {
	if rule == nil {
		return fmt.Errorf("source rule is nil")
	}

	if _, valid := ruleStateKey(rule); !valid {
		return validationError("source rule contains an invalid canonical state")
	}

	revisionID, err := common.NewUUID()
	if err != nil {
		return fmt.Errorf("建立來源規則 revision id 失敗: %w", err)
	}
	rule.RevisionID = revisionID
	return nil
}
