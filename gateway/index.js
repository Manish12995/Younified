const { ApolloServer } = require('apollo-server');
const { ApolloGateway, IntrospectAndCompose } = require("@apollo/gateway");

const gateway = new ApolloGateway({
    supergraphSdl: new IntrospectAndCompose({
        subgraphs: [
            { name: 'union', url: 'http://localhost:4001/graphql' },
            { name: 'user', url: 'http://localhost:4002/graphql' },
            {name: 'comms', url: 'http://localhost:4003/graphql'}
        ]
    })
});

const server = new ApolloServer({
    gateway,

    subscriptions: false,
});

const PORT = 4000; 

server.listen({ port: PORT }).then(({ url }) => {
    console.log(`🚀 Server ready at ${url}`);
});