package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/k1LoW/tbls/schema"
)

func TestAnalayzeSchema(t *testing.T)	{
	tests := []struct {
		name string
		strOrPath string
		includes []string
		excludes []string
		labels []string
		want *schema.Schema
		wantErr bool
	} {
		{
			name: "analyze string",
			strOrPath: `{"name": "test", "tables": [{"name": "a", "comment": "table a", "columns": [{"name": "id", "type": "int"}]},{"name": "b", "comment": "table b", "columns": [{"name": "title", "type": "varchar"}]}]}`,
			includes: []string{},
			excludes: []string{},
			labels: []string{},
			want: &schema.Schema{
				Name: "test",
				Tables: []*schema.Table{
					{
						Name: "a",
						Comment: "table a",
						Columns: []*schema.Column{
							{
								Name: "id",
								Type: "int",
							},
						},
					},
					{
						Name: "b",
						Comment: "table b",
						Columns: []*schema.Column{
							{
								Name: "title",
								Type: "varchar",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AnalyzeSchema(tt.strOrPath, tt.includes, tt.excludes, tt.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("AnalyzeSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AnalyzeSchema() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
