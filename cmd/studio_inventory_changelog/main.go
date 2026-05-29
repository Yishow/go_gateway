package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	_ "modernc.org/sqlite"
)

const (
	defaultDBPath = "docs/technical/studio-surface-inventory/changelog.sqlite"
	schemaSQL     = `CREATE TABLE IF NOT EXISTS inventory_changelog (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		changed_at TEXT NOT NULL,
		actor TEXT NOT NULL,
		surface TEXT NOT NULL,
		summary TEXT NOT NULL,
		files_json TEXT NOT NULL,
		reason TEXT NOT NULL DEFAULT '',
		notes TEXT NOT NULL DEFAULT ''
	);`
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "studio inventory changelog: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "init":
		return runInit(args[1:])
	case "add":
		return runAdd(args[1:])
	case "list":
		return runList(args[1:])
	default:
		printUsage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	dbPath := fs.String("db", defaultDBPath, "path to the sqlite changelog database")
	if err := fs.Parse(args); err != nil {
		return err
	}

	db, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureSchema(db); err != nil {
		return err
	}

	fmt.Printf("initialized %s\n", *dbPath)
	return nil
}

func runAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	dbPath := fs.String("db", defaultDBPath, "path to the sqlite changelog database")
	actor := fs.String("actor", "codex", "actor name")
	surface := fs.String("surface", "studio-surface-inventory", "surface or section being updated")
	summary := fs.String("summary", "", "short summary of the change")
	files := fs.String("files", "", "comma-separated files touched")
	reason := fs.String("reason", "", "why the change was made")
	notes := fs.String("notes", "", "extra notes for future maintainers")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*summary) == "" {
		return fmt.Errorf("summary is required")
	}

	db, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureSchema(db); err != nil {
		return err
	}

	filesJSON, err := marshalFiles(*files)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		`INSERT INTO inventory_changelog (
			changed_at, actor, surface, summary, files_json, reason, notes
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		time.Now().UTC().Format(time.RFC3339),
		strings.TrimSpace(*actor),
		strings.TrimSpace(*surface),
		strings.TrimSpace(*summary),
		filesJSON,
		strings.TrimSpace(*reason),
		strings.TrimSpace(*notes),
	)
	if err != nil {
		return fmt.Errorf("insert changelog entry: %w", err)
	}

	fmt.Printf("recorded changelog entry in %s\n", *dbPath)
	return nil
}

func runList(args []string) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	dbPath := fs.String("db", defaultDBPath, "path to the sqlite changelog database")
	limit := fs.Int("limit", 20, "max entries to display")
	if err := fs.Parse(args); err != nil {
		return err
	}

	db, err := openDB(*dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := ensureSchema(db); err != nil {
		return err
	}

	rows, err := db.Query(
		`SELECT changed_at, actor, surface, summary, files_json, reason
		FROM inventory_changelog
		ORDER BY id DESC
		LIMIT ?`,
		*limit,
	)
	if err != nil {
		return fmt.Errorf("query changelog entries: %w", err)
	}
	defer rows.Close()

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(writer, "changed_at\tactor\tsurface\tsummary\tfiles\treason")

	hasRows := false
	for rows.Next() {
		hasRows = true
		var changedAt string
		var actor string
		var surface string
		var summary string
		var filesJSON string
		var reason string
		if err := rows.Scan(&changedAt, &actor, &surface, &summary, &filesJSON, &reason); err != nil {
			return fmt.Errorf("scan changelog entry: %w", err)
		}

		files := decodeFiles(filesJSON)
		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\t%s\t%s\n",
			changedAt,
			actor,
			surface,
			summary,
			strings.Join(files, ", "),
			reason,
		)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate changelog entries: %w", err)
	}
	if !hasRows {
		fmt.Fprintln(writer, "(no entries)\t\t\t\t\t")
	}

	return writer.Flush()
}

func openDB(path string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create changelog directory: %w", err)
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	return db, nil
}

func ensureSchema(db *sql.DB) error {
	if _, err := db.Exec(schemaSQL); err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}

	return nil
}

func marshalFiles(raw string) (string, error) {
	parts := strings.Split(raw, ",")
	files := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		files = append(files, trimmed)
	}

	encoded, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("marshal files list: %w", err)
	}

	return string(encoded), nil
}

func decodeFiles(filesJSON string) []string {
	var files []string
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return []string{filesJSON}
	}

	return files
}

func printUsage() {
	fmt.Println("usage:")
	fmt.Println("  go run ./cmd/studio_inventory_changelog init")
	fmt.Println("  go run ./cmd/studio_inventory_changelog add -summary \"...\" [-surface ...] [-files file1,file2]")
	fmt.Println("  go run ./cmd/studio_inventory_changelog list [-limit 20]")
}
