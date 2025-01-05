package graphqlclient

// GraphQLError represents an error returned by a GraphQL endpoint
type GraphQLError struct {
	Message    string                 `json:"message"`
	Locations  []ErrorLocation        `json:"locations"`
	Path       []string               `json:"path"`
	Extensions map[string]interface{} `json:"extensions"`
}

// ErrorLocation represents the location of an error in the GraphQL query
type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func (e GraphQLError) Error() string {
	return e.Message
}
