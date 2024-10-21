package analyzer

import (
	"fmt"
	"strings"

	"github.com/k1LoW/tbls/config"
	"github.com/k1LoW/tbls/datasource"
	"github.com/k1LoW/tbls/schema"
)

func AnalyzeSchema(strOrPath string, includes []string, excludes []string, labels []string) (s *schema.Schema, err error) {
	if strings.HasPrefix(strOrPath, "{") {
		s, err = datasource.AnalyzeJSONStringOrFile(strOrPath)
	} else {
		dsn := config.DSN{URL: strOrPath}
		s, err = datasource.Analyze(dsn)
	}
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
