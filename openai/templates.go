package openai

import (
	"fmt"
	"strings"

	"github.com/k1LoW/tbls/schema"
)

const (
	defaultPromtTmpl = `Answer the questions in the Question assuming the following DDL.
{{ .DatabaseVersion }}

## DDL ( Data Definition Language )

{{ .QuoteStart }}
{{ .DDL }}
{{ .QuoteEnd }}

## Question
{{ .Question }}
`
	defaultQueryPromptTmpl = `Answer the SQL query in the "Explanation of the query to be created" section, assuming the database was created with the following DDL.
{{ .DatabaseVersion }}

## DDL ( Data Definition Language )

{{ .QuoteStart }}
{{ .DDL }}
{{ .QuoteEnd }}

## Explanation of the query to be created
{{ .Question }}
`
)

func generateDDLRoughly(s *schema.Schema) string {
	var ddl string
	for _, t := range s.Tables {
		if t.Type == "VIEW" {
			continue
		}
		ddl += fmt.Sprintf("CREATE TABLE %s (\n", t.Name)
		td := []string{}
		for _, c := range t.Columns {
			d := fmt.Sprintf("  %s %s", c.Name, c.Type)
			if c.Default.String != "" {
				d += fmt.Sprintf(" DEFAULT %s", c.Default.String)
			}
			if c.Nullable {
				d += " NULL"
			} else {
				d += " NOT NULL"
			}
			if c.Comment != "" {
				d += fmt.Sprintf(" COMMENT %q", c.Comment)
			}
			td = append(td, d)
		}
		for _, i := range t.Indexes {
			d := fmt.Sprintf("  %s", i.Def)
			td = append(td, d)
		}
		for _, c := range t.Constraints {
			switch c.Type {
			case "PRIMARY KEY", "UNIQUE KEY":
				continue
			default:
				d := fmt.Sprintf("  CONSTRAINT %s", c.Def)
				td = append(td, d)
			}
		}
		ddl += fmt.Sprintf("%s\n", strings.Join(td, ",\n"))
		ddl += fmt.Sprintf(") COMMENT = %q;\n\n", t.Comment)
	}
	return ddl
}
