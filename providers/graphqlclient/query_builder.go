package graphqlclient

import (
	"fmt"
	"strings"
)

// QueryBuilder helps construct GraphQL queries
type QueryBuilder struct {
	name       string
	fields     []string
	arguments  map[string]interface{}
	fragments  []string
	pagination *paginationConfig
}

type paginationConfig struct {
	first int
	after string
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		name:      "",
		fields:    []string{},
		arguments: make(map[string]interface{}),
	}
}

// SetMutationName allows dynamically changing the mutation name after initial creation
func (qb *QueryBuilder) SetQueryName(name string) *QueryBuilder {
	qb.name = name
	return qb
}

// AddField adds a field to the query
func (qb *QueryBuilder) AddField(field string) *QueryBuilder {
	qb.fields = append(qb.fields, field)
	return qb
}

// AddArgument adds an argument to the query
func (qb *QueryBuilder) AddArgument(key string, value interface{}) *QueryBuilder {
	qb.arguments[key] = value
	return qb
}

// AddFragment adds a fragment to the query
func (qb *QueryBuilder) AddFragment(fragment string) *QueryBuilder {
	qb.fragments = append(qb.fragments, fragment)
	return qb
}

// WithPagination adds pagination to the query
func (qb *QueryBuilder) WithPagination(first int, after string) *QueryBuilder {
	qb.pagination = &paginationConfig{
		first: first,
		after: after,
	}
	return qb
}

// Build constructs the final GraphQL query
func (qb *QueryBuilder) Build() string {
	var args []string
	for k, v := range qb.arguments {
		args = append(args, fmt.Sprintf("%s: %v", k, formatValue(v)))
	}

	// Add pagination arguments
	if qb.pagination != nil {
		args = append(args, fmt.Sprintf("first: %d", qb.pagination.first))
		if qb.pagination.after != "" {
			args = append(args, fmt.Sprintf("after: \"%s\"", qb.pagination.after))
		}
	}

	argsStr := ""
	if len(args) > 0 {
		argsStr = fmt.Sprintf("(%s)", strings.Join(args, ", "))
	}

	query := fmt.Sprintf(`
		query %s%s {
			%s%s {
				%s
			}
		}
		%s
	`,
		qb.name,
		argsStr,
		qb.name,
		argsStr,
		strings.Join(qb.fields, " "),
		strings.Join(qb.fragments, "\n"),
	)

	return query
}

// formatValue handles different types for query arguments
func formatValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", val)
	case int, int64, float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%v", val)
	case []string:
		formatted := make([]string, len(val))
		for i, s := range val {
			formatted[i] = fmt.Sprintf("\"%s\"", s)
		}
		return fmt.Sprintf("[%s]", strings.Join(formatted, ", "))
	default:
		return "null"
	}
}
