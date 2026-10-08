package extracurricularRepository

import (
	"strings"
	"testing"

	"inspirate-consulting/internal/data/postgres/schema"
)

func TestExtracurricularSQLFiles(t *testing.T) {
	t.Parallel()

	for _, filename := range []string{
		"create_extracurricular.sql",
		"list_extracurriculars.sql",
		"update_extracurricular.sql",
	} {
		t.Run(filename, func(t *testing.T) {
			t.Parallel()

			query, err := schema.ReadSQLBaseScript(filename, SqlExtracurricularFiles)
			if err != nil {
				t.Fatalf("read embedded SQL: %v", err)
			}
			if strings.TrimSpace(query) == "" {
				t.Fatal("embedded SQL query is empty")
			}
		})
	}
}
