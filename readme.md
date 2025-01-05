# Younified Backend

## Overview

Younified is a comprehensive union management system built with a microservices architecture, utilizing GraphQL for API communication and Apollo Federation for gateway services.

## Project Structure

```
younified-backend/
│
├── contracts/                 # GraphQL contract definitions
│   ├── union/
│   │   ├── graph/
│   │   │   └── union.graphql  # Union service GraphQL schema
│   │   └── model/
│   │       ├── model.go       # Union data models
│   │       └── model_gen.go   # Generated GraphQL models
│   │
│   └── user/
│       ├── graph/
│       │   └── user.graphql   # User service GraphQL schema
│       └── model/
│           ├── model.go       # User data models
│           └── model_gen.go   # Generated GraphQL models
│
├── gateway/                   # Apollo Federated GraphQL Gateway
│   ├── index.js               # Gateway server configuration
│   └── package.json           # NPM dependencies
│
├── providers/                 # Shared providers and utilities
│   ├── database/
│   │   ├── mongodb.manager.go # MongoDB connection and management
│   │   └── redis.go           # Redis client handler
│   │
│   ├── graphqlclient/         # Custom GraphQL client
│       ├── cache.go           # Caching implementation
│       ├── client.go          # GraphQL client instance
│       ├── mutation_builder.go# Mutation string builder
│       ├── query_builder.go   # Query string builder
│       └── types.go           # Error structs
│
├── services/
│       ├── unionService/      # Union management microservice
│       │   ├── internal/
│       │   │   ├── controller/
│       │   │   │   └── union.controller.go
│       │   │   ├── repository/
│       │   │   │   ├── mongo_repository.go
│       │   │   │   └── redis_repository.go
│       │   │   └── resolvers/
│       │   │       ├── union.resolvers.go
│       │   │       ├── federation.go
│       │   │       ├── generated.go
│       │   │       └── resolver.go
│       │   ├── .env
│       │   ├── go.mod
│       │   ├── gqlgen.yml
│       │   └── main.go
│       │
│       └── userService/       # User management microservice
│           ├── internal/
│           │   ├── controller/
│           │   │   └── user.controller.go
│           │   ├── repository/
│           │   │   ├── mongo_repository.go
│           │   │   └── redis_repository.go
│           │   └── resolvers/
│           │       ├── user.resolvers.go
│           │       ├── federation.go
│           │       ├── generated.go
│           │       └── resolver.go
│           ├── .env
│           ├── go.mod
│           ├── gqlgen.yml
│           └── main.go
```

## Services Overview

### Union Service
Responsible for:
- Creating unions
- Approving unions
- Assigning managers and admins
- Adding first and default users
- Creating union-specific databases on MongoDB Atlas

### User Service
Responsible for:
- User registration
- User approval
- Bulk user uploading
- Role-Based Access Control (RBAC)

### Communications Service
Responsible for:
- Text message blasts
- Email blasts
- Peer-to-peer calls
- Notifications

## Prerequisites

- Go (1.23.2)
- Node.js (18+)
- MongoDB Atlas account
- Redis
- Apollo Federation (knowledge)

## Setup

### Environment Configuration

1. Clone the repository
```bash
git clone https://github.com/SpeedAD-CIS/younified-backend.git
cd younified-backend
```

2. Set up environment variables
- Create `.env` files in each service directory
- Configure MongoDB connection strings
- Configure Redis connection
- Set up authentication credentials

### Database Setup

1. Create a MongoDB Atlas cluster
2. Configure connection in `providers/database/mongodb.manager.go`
3. Set up Redis instance

### Running Services

#### Union Service
```bash
cd services/unionService
go mod tidy
go run main.go
```

#### User Service
```bash
cd services/userService
go mod tidy
go run main.go
```

#### Gateway
```bash
cd gateway
npm install
npm start
```

## Technologies

- Backend: Go
- GraphQL: gqlgen
- API Gateway: Apollo Federation
- Database: MongoDB Atlas
- Caching: Redis
- Inter-service Communication: GraphQLClient

## Architecture

- Microservices architecture
- GraphQL for API design
- Apollo Federation for gateway
- Separation of concerns
- Modular service design

## Contributing

1. Clone the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## License

Created with ❤️ by Speed 😎 & Nishid

This project is proudly crafted by Speed for Younified

**Licensing Terms:**
- Use it, break it, make it better
- Give credit where it's due

## Contact & Contributers

agam.d@cisinlabs.com
nishid.p@cisinlabs.com

Project Link: [https://github.com/SpeedAD-CIS/younified-backend](https://github.com/SpeedAD-CIS/younified-backend)