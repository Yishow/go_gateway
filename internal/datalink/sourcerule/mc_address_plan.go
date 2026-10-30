package sourcerule

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"go-gateway/internal/datalink/schema"
)

const (
	mcAddressPlanConflictCode = "mc_address_plan_conflict"
	mcClientMaxAddress        = 0xFFFFFF
)

var errMCAddressPlanConflict = errors.New(mcAddressPlanConflictCode)

type mcAddressPlanConflictError struct {
	ruleID       string
	startAddress string
	reason       string
}

func (e *mcAddressPlanConflictError) Error() string {
	return fmt.Sprintf("%s rule_id=%s start_address=%s: %s", mcAddressPlanConflictCode, e.ruleID, e.startAddress, e.reason)
}

func (e *mcAddressPlanConflictError) Unwrap() []error {
	return []error{errMCAddressPlanConflict, validationError(e.reason)}
}

func (s *Service) listAndValidatePersistedMCLinks(ctx context.Context, rule *schema.SourceRule, protocol schema.ProtocolType) ([]*schema.SourceRuleLink, error) {
	links, err := s.repo.ListLinks(ctx, rule.ID)
	if err != nil {
		return nil, fmt.Errorf("取得來源規則連結失敗: %w", err)
	}
	if err := validatePersistedMCAddressPlan(rule, protocol, links); err != nil {
		return nil, contextualRuleValidationError(rule.ID, rule.StartAddress, protocol, err)
	}
	return links, nil
}

// validatePersistedMCAddressPlan protects persisted MC links from an address
// grammar correction. It intentionally plans from the stored rule, not from a
// pending update, so a name/scale/range patch cannot move existing identities.
func validatePersistedMCAddressPlan(rule *schema.SourceRule, protocol schema.ProtocolType, links []*schema.SourceRuleLink) error {
	if protocol != schema.ProtocolMC3E {
		return nil
	}

	expected, err := correctedMCAddressPlan(rule)
	if err != nil {
		return mcAddressPlanConflict(rule, err.Error())
	}
	persisted := make([]string, 0, len(links))
	for _, link := range links {
		if link == nil {
			return mcAddressPlanConflict(rule, "persisted link is nil")
		}
		persisted = append(persisted, normalizeAddressKey(link.Address))
	}
	slices.Sort(expected)
	slices.Sort(persisted)
	if !slices.Equal(expected, persisted) {
		return mcAddressPlanConflict(rule, fmt.Sprintf("stored=%v persisted=%v", expected, persisted))
	}
	return nil
}

func mcAddressPlanConflict(rule *schema.SourceRule, reason string) error {
	ruleID := "<unknown>"
	startAddress := "<empty>"
	if rule != nil {
		ruleID = strings.TrimSpace(rule.ID)
		startAddress = strings.TrimSpace(strings.ToUpper(rule.StartAddress))
	}
	return &mcAddressPlanConflictError{ruleID: ruleID, startAddress: startAddress, reason: reason}
}

func correctedMCAddressPlan(rule *schema.SourceRule) ([]string, error) {
	if rule == nil {
		return nil, validationError("stored rule is required")
	}
	if rule.Count <= 0 {
		return nil, validationError("stored count must be greater than zero")
	}
	if !rule.DataType.IsValid() {
		return nil, validationError(fmt.Sprintf("stored data_type is invalid: %s", rule.DataType))
	}

	area, suffix, radix, err := splitCorrectedMCAddress(rule.StartAddress)
	if err != nil {
		return nil, err
	}
	value, err := parseCorrectedMCSuffix(suffix, radix)
	if err != nil {
		return nil, err
	}
	skipped, err := parseSkippedAddresses(rule.SkippedAddresses)
	if err != nil {
		return nil, err
	}

	span := getDataTypeCellSpan(rule.DataType)
	maxIndex := (mcClientMaxAddress - value) / span
	if rule.Count-1 > maxIndex {
		return nil, validationError(fmt.Sprintf("MC address plan exceeds client range 0-%X", mcClientMaxAddress))
	}
	addresses := make([]string, 0, rule.Count)
	for index := range rule.Count {
		address := formatCorrectedMCAddress(area, value+index*span, radix)
		if containsAddress(skipped, address) {
			continue
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

func splitCorrectedMCAddress(raw string) (area, suffix string, radix int, err error) {
	address := strings.ToUpper(strings.TrimSpace(raw))
	if len(address) < 2 {
		return "", "", 0, validationError(fmt.Sprintf("invalid stored MC address: %s", raw))
	}
	area = address[:1]
	suffix = address[1:]
	switch area {
	case "X", "Y":
		radix = 8
	case "W", "B":
		radix = 16
	case "D", "M":
		radix = 10
	default:
		return "", "", 0, validationError(fmt.Sprintf("invalid stored MC area: %s", raw))
	}
	return area, suffix, radix, nil
}

func parseCorrectedMCSuffix(raw string, radix int) (int, error) {
	if raw == "" {
		return 0, validationError("invalid stored MC suffix: empty")
	}
	if radix != 8 && radix != 10 && radix != 16 {
		return 0, validationError(fmt.Sprintf("unsupported MC radix: %d", radix))
	}
	for _, ch := range raw {
		valid := ch >= '0' && ch <= '9'
		switch radix {
		case 8:
			valid = valid && ch <= '7'
		case 16:
			valid = valid || (ch >= 'A' && ch <= 'F')
		}
		if !valid {
			return 0, validationError(fmt.Sprintf("invalid stored MC suffix: %s", raw))
		}
	}
	value, err := strconv.ParseUint(raw, radix, 64)
	if err != nil {
		return 0, validationError(fmt.Sprintf("invalid stored MC suffix: %s", raw))
	}
	if value > mcClientMaxAddress {
		return 0, validationError(fmt.Sprintf("MC address offset out of range: %s (valid: 0-%X)", raw, mcClientMaxAddress))
	}
	return int(value), nil
}

func formatCorrectedMCAddress(area string, value, radix int) string {
	return area + strings.ToUpper(strconv.FormatInt(int64(value), radix))
}

func correctedMCOffsetAddress(raw string, delta int) (string, error) {
	area, suffix, radix, err := splitCorrectedMCAddress(raw)
	if err != nil {
		return "", err
	}
	value, err := parseCorrectedMCSuffix(suffix, radix)
	if err != nil {
		return "", err
	}
	if delta > 0 {
		if delta > mcClientMaxAddress || value > mcClientMaxAddress-delta {
			return "", validationError(fmt.Sprintf("MC address offset exceeds client range 0-%X", mcClientMaxAddress))
		}
		return formatCorrectedMCAddress(area, value+delta, radix), nil
	}
	if delta < -value {
		return formatCorrectedMCAddress(area, 0, radix), nil
	}
	return formatCorrectedMCAddress(area, value+delta, radix), nil
}
