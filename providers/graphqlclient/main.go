package graphqlclient

type Graph struct {
	gqlClient       *Client
	mutationBuilder *MutationBuilder
	queryBuilder    *QueryBuilder
}

func NewGraphql() *Graph {
	gqlClient := NewClient()
	mutationBuilder := NewMutationBuilder()
	queryBuilder := NewQueryBuilder()
	return &Graph{
		gqlClient:       gqlClient,
		mutationBuilder: mutationBuilder,
		queryBuilder:    queryBuilder,
	}
}

func (g *Graph) GetMutationBuilder() *MutationBuilder {
	return g.mutationBuilder
}

func (g *Graph) GetQueryBuilder() *QueryBuilder {
	return g.queryBuilder
}
