package analyzer

import (
	"fmt"

	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
	"github.com/k1LoW/tbls/config"
)

func AnalyzeSchema(path string, includes []string, excludes []string, labels []string) (s *schema.Schema, err error) {
	dsn := config.DSN{URL: path}
	s, err = datasource.Analyze(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze schema: %w", err)
	}

	if err := s.Filter(&schema.FilterOption{
		Include:       includes,
		Exclude:       excludes,
		IncludeLabels: labels,
	}); err != nil {
		return nil, fmt.Errorf("failed to filter schema: %w", err)
	}

	return s, nil
}
