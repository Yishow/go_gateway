package recordingplan

import (
	"encoding/json"
	"strings"
)

type successorHistoryFacts struct {
	attributed bool
	roots      map[string]bool
}

// Scan only exact private envelope paths. Facts parsed before a malformed tail
// remain contradictory evidence; duplicate keys cannot hide an earlier match.
// This is negative evidence only, not a permissive decoder for ledger results.
func successorHistoryAttribution(raw, lineageID string, groups map[string]bool) successorHistoryFacts {
	facts := successorHistoryFacts{roots: map[string]bool{}}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	_ = scanSuccessorHistoryValue(decoder, nil, lineageID, groups, &facts)
	return facts
}

func scanSuccessorHistoryValue(decoder *json.Decoder, path []string, lineageID string, groups map[string]bool, facts *successorHistoryFacts) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if value, ok := token.(string); ok {
		if len(path) == 1 && strings.EqualFold(path[0], "reserve_operation_id") && value != "" {
			facts.roots[value] = true
		}
		if len(path) == 2 && strings.EqualFold(path[0], "scope") {
			if strings.EqualFold(path[1], "lineage_id") && value == lineageID {
				facts.attributed = true
			}
			if (strings.EqualFold(path[1], "predecessor_group_id") || strings.EqualFold(path[1], "successor_group_id")) && groups[value] {
				facts.attributed = true
			}
		}
	}
	delimiter, ok := token.(json.Delim)
	if !ok {
		return nil
	}
	switch delimiter {
	case '{':
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return err
			}
			key, ok := keyToken.(string)
			if !ok {
				return ErrSchemaOperationResult
			}
			if err = scanSuccessorHistoryValue(decoder, append(path, key), lineageID, groups, facts); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim('}') {
			return ErrSchemaOperationResult
		}
	case '[':
		for decoder.More() {
			if err = scanSuccessorHistoryValue(decoder, append(path, ""), lineageID, groups, facts); err != nil {
				return err
			}
		}
		end, err := decoder.Token()
		if err != nil {
			return err
		}
		if end != json.Delim(']') {
			return ErrSchemaOperationResult
		}
	default:
		return ErrSchemaOperationResult
	}
	return nil
}
