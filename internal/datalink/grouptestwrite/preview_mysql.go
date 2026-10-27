package grouptestwrite

import (
	"fmt"
	"strings"
	"time"
)

// mysqlPreview derives the real rows the confirmation will write. The table
// owner marker stays gateway; cleanup names exact record IDs, never an owner
// predicate.
func mysqlPreview(p *plan, owner string, at time.Time) (values []PreviewValue, recordIDs []string, cleanup string, err error) {
	rows, err := p.mysqlRows(owner, at)
	if err != nil {
		return nil, nil, "", err
	}
	if len(rows) == 0 {
		return nil, nil, "", unsupported(ReasonGroupNotTestable)
	}
	recordIDs = make([]string, len(rows))
	for i, row := range rows {
		recordIDs[i] = row.RecordID
	}
	types := make(map[string]string, len(p.columns))
	for _, column := range p.columns {
		types[strings.ToLower(column.Name)] = column.DataType
	}
	values = make([]PreviewValue, 0, len(rows[0].Cells))
	for _, cell := range rows[0].Cells {
		values = append(values, PreviewValue{
			Column: cell.Column, Type: previewType(types[strings.ToLower(cell.Column)]), Value: previewValue(cell.Value),
		})
	}
	cleanup = "remove only the previewed record IDs and their matching receipts; retain all other rows"
	return values, recordIDs, cleanup, nil
}

func previewType(sqlType string) string {
	if strings.TrimSpace(sqlType) == "" {
		return "text"
	}
	return strings.ToLower(sqlType)
}

func previewValue(value any) string {
	switch typed := value.(type) {
	case nil:
		return "NULL"
	case []byte:
		return string(typed)
	case time.Time:
		return typed.UTC().Format(time.RFC3339Nano)
	default:
		return fmt.Sprint(typed)
	}
}
