package dbtarget

import (
	"context"
	"database/sql"
	"strings"
	"unicode"
)

const sqlServerManagedCheckConstraintsQuery = `
	SELECT cc.name, cc.definition, cc.is_disabled, cc.is_not_trusted, cc.parent_column_id
	FROM sys.check_constraints AS cc
	WHERE cc.parent_object_id = @p1
	ORDER BY cc.object_id`

func readSQLServerManagedProvenanceCheck(ctx context.Context, db sqlServerManagedCatalogQuerier, objectID, provenanceColumnID int64) (bool, error) {
	rows, err := db.QueryContext(ctx, sqlServerManagedCheckConstraintsQuery, objectID)
	if err != nil {
		return false, err
	}
	checked := false
	for rows.Next() {
		var name string
		var definition sql.NullString
		var disabled, notTrusted bool
		var parentColumnID int64
		if err := rows.Scan(&name, &definition, &disabled, &notTrusted, &parentColumnID); err != nil {
			_ = rows.Close()
			return false, err
		}
		belongsToProvenance := parentColumnID == 0 || parentColumnID == provenanceColumnID
		if provenanceColumnID > 0 && definition.Valid && !disabled && !notTrusted && belongsToProvenance &&
			sqlServerCanonicalProvenanceCheck(definition.String) {
			checked = true
		}
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return false, err
	}
	if err := rows.Close(); err != nil {
		return false, err
	}
	return checked, nil
}

func sqlServerCanonicalProvenanceCheck(definition string) bool {
	var normalized strings.Builder
	quoted := false
	for index := 0; index < len(definition); index++ {
		char := definition[index]
		if char == '\'' {
			normalized.WriteByte(char)
			if quoted && index+1 < len(definition) && definition[index+1] == '\'' {
				normalized.WriteByte(definition[index+1])
				index++
				continue
			}
			quoted = !quoted
			continue
		}
		if !quoted {
			if char == '[' || char == ']' || unicode.IsSpace(rune(char)) {
				continue
			}
		}
		normalized.WriteByte(char)
	}
	expression := strings.ToLower(sqlServerStripOuterParentheses(normalized.String()))
	return expression == "isjson(provenance)=1" || expression == "isjson(provenance)=(1)"
}

func sqlServerCanonicalDefault(definition string) string {
	expression := sqlServerStripOuterParentheses(strings.TrimSpace(definition))
	if len(expression) >= 2 && (expression[0] == 'N' || expression[0] == 'n') && expression[1] == '\'' {
		expression = expression[1:]
	}
	if len(expression) < 2 || expression[0] != '\'' || expression[len(expression)-1] != '\'' {
		return strings.TrimSpace(expression)
	}
	body := expression[1 : len(expression)-1]
	var canonical strings.Builder
	for index := 0; index < len(body); index++ {
		if body[index] != '\'' {
			canonical.WriteByte(body[index])
			continue
		}
		if index+1 >= len(body) || body[index+1] != '\'' {
			return strings.TrimSpace(expression)
		}
		canonical.WriteByte('\'')
		index++
	}
	return canonical.String()
}

func sqlServerStripOuterParentheses(expression string) string {
	expression = strings.TrimSpace(expression)
	for len(expression) >= 2 && expression[0] == '(' && expression[len(expression)-1] == ')' {
		depth := 0
		quoted := false
		wrapsAll := true
		for index := 0; index < len(expression); index++ {
			char := expression[index]
			if char == '\'' {
				if quoted && index+1 < len(expression) && expression[index+1] == '\'' {
					index++
					continue
				}
				quoted = !quoted
				continue
			}
			if quoted {
				continue
			}
			switch char {
			case '(':
				depth++
			case ')':
				depth--
				if depth < 0 || (depth == 0 && index != len(expression)-1) {
					wrapsAll = false
				}
			}
		}
		if !wrapsAll || quoted || depth != 0 {
			break
		}
		expression = strings.TrimSpace(expression[1 : len(expression)-1])
	}
	return expression
}
