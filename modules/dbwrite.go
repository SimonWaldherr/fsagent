package modules

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var dbSafePath = regexp.MustCompile(`^[a-zA-Z0-9_./\-]+$`)
var dbSafeIdent = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

type dbWriteConfig struct {
	Path           string `json:"path"`
	Table          string `json:"table"`
	FilenameColumn string `json:"filenameColumn"`
	ContentColumn  string `json:"contentColumn"`
}

type DBWrite struct{}

func (DBWrite) Name() string {
	return "dbwrite"
}

func (DBWrite) EmptyConfig() interface{} {
	return &dbWriteConfig{}
}

func (DBWrite) Perform(config interface{}, fileName string) error {
	c := config.(*dbWriteConfig)
	if !dbSafePath.MatchString(c.Path) {
		return fmt.Errorf("invalid sqlite path")
	}
	if !dbSafeIdent.MatchString(c.Table) {
		return fmt.Errorf("invalid table name")
	}
	if c.FilenameColumn == "" {
		c.FilenameColumn = "filename"
	}
	if c.ContentColumn == "" {
		c.ContentColumn = "content"
	}
	if !dbSafeIdent.MatchString(c.FilenameColumn) || !dbSafeIdent.MatchString(c.ContentColumn) {
		return fmt.Errorf("invalid column name")
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	if _, err := exec.LookPath("sqlite3"); err != nil {
		return fmt.Errorf("sqlite3 binary not found: %v", err)
	}

	createSQL := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, %s TEXT, %s TEXT);",
		c.Table, c.FilenameColumn, c.ContentColumn,
	)
	if _, err := exec.Command("sqlite3", c.Path, createSQL).CombinedOutput(); err != nil {
		return err
	}

	insertSQL := fmt.Sprintf(
		"INSERT INTO %s (%s, %s) VALUES (%s, %s);",
		c.Table, c.FilenameColumn, c.ContentColumn, dbQuote(fileName), dbQuote(string(content)),
	)
	_, err = exec.Command("sqlite3", c.Path, insertSQL).CombinedOutput()
	return err
}

func dbQuote(v string) string {
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}
