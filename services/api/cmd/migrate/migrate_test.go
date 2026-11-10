package main

import (
	"strings"
	"testing"
)

func TestExtractUpSegment(t *testing.T) {
	sampleSQL := `
-- +goose Up
CREATE TABLE test_table (id INT);

-- +goose Down
DROP TABLE test_table;
`
	up := extractUpSegment(sampleSQL)
	if !strings.Contains(up, "CREATE TABLE test_table") {
		t.Errorf("expected Up segment to contain CREATE TABLE, got: %s", up)
	}
	if strings.Contains(up, "DROP TABLE test_table") {
		t.Errorf("expected Up segment NOT to contain DROP TABLE, got: %s", up)
	}
}
