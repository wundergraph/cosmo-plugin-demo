# Cosmo Router Users Plugin

A Cosmo Router plugin for user management. This plugin provides both local mock users and integration with external user data from JSONPlaceholder.

## Features

- Local mock user management
- External user data fetching from JSONPlaceholder API
- gRPC service implementation for GraphQL resolvers

## Testing

The integration tests validate all functionality of the plugin including:

- User lookup by ID (single and batch)
- User querying (single and all)
- User updates
- External user data fetching

### Running the Tests

To run the integration tests:

```shell
cd cosmo-router/plugins
wgc plugin test users
```

### Test Structure

The tests use:

- In-memory gRPC server for testing the service
- Mock HTTP server for testing external API calls
- Table-driven tests for different test cases
- Cleanup functions to ensure proper test isolation

## Development

This plugin is structured as a standard Go module. The main implementation is in `src/main.go`, with integration tests in `src/main_test.go`.

# Users Plugin for Cosmo Router

This plugin implements a GraphQL Users Subgraph as a Go-based gRPC Router plugin that runs directly inside the Cosmo Router, eliminating the need for a separate service.

## Overview

The Users Plugin provides a GraphQL API for managing user data with the following features:

- User data storage and retrieval (mock implementation)
- User lookup by ID
- List all users
- Update user information
- Integration with external REST APIs (JSONPlaceholder for external users)
- Federation support with `@key` directive

## Implementation

The plugin is implemented as a gRPC service that integrates directly with the Cosmo Router:

- **GraphQL Schema**: Defines the data model with User type, queries, and mutations
- **gRPC Service**: Implements the GraphQL resolvers as gRPC methods
- **Router Integration**: Runs as an embedded plugin within the Cosmo Router
- **REST Integration**: Connects to JSONPlaceholder API to fetch external user data

## Directory Structure

```
users/
├── bin/                # Compiled plugin binaries
├── generated/          # Auto-generated gRPC code from schema
│   ├── service.pb.go
│   ├── service_grpc.pb.go
│   └── service.proto
├── src/                # Source code
│   ├── main.go         # Plugin implementation
│   ├── main_test.go    # Integration tests
│   └── schema.graphql  # GraphQL schema definition
└── go.mod              # Go module dependencies
```

## Development

### Prerequisites

- [Cosmo CLI](https://docs.wundergraph.com/docs/cosmo/cli/overview)

### Building

The plugin can be built with:

```bash
npm run build
```

## Usage

The Users Plugin is automatically loaded by the Cosmo Router when properly configured. The GraphQL API exposes:

### Queries

- `user(id: ID!)`: Get a user by ID
- `users`: List all users
- `externalUser(id: ID!)`: Get an external user by ID from JSONPlaceholder
- `externalUsers`: List all external users from JSONPlaceholder

### Mutations

- `updateUser(id: ID!, input: UserInput!)`: Update user information

## Example GraphQL Queries

```graphql
# Get all users
query {
  users {
    id
    name
    email
    role
  }
}

# Get user by ID
query {
  user(id: "1") {
    id
    name
    email
    role
  }
}

# Update a user
mutation {
  updateUser(id: "1", input: {
    name: "Alice Updated"
    email: "alice.updated@example.com"
    role: ADMIN
  }) {
    id
    name
    email
    role
  }
}

# Get all external users from JSONPlaceholder API
query {
  externalUsers {
    id
    name
    username
    email
    phone
    website
    company {
      name
      catchPhrase
      bs
    }
    address {
      street
      city
      zipcode
      geo {
        lat
        lng
      }
    }
  }
}

# Get external user by ID from JSONPlaceholder API
query {
  externalUser(id: "1") {
    id
    name
    email
    username
    phone
    website
  }
}
```

## REST API Integration

The plugin demonstrates how to integrate external REST APIs into your GraphQL schema:

- **JSONPlaceholder Integration**: The plugin fetches user data from the [JSONPlaceholder](https://jsonplaceholder.typicode.com/) REST API
- **HTTP Client**: Uses Go's standard HTTP client with custom configuration
- **Data Transformation**: Converts JSON responses to GraphQL-compatible types
- **Error Handling**: Properly handles HTTP errors and response parsing
- **Rich Data Structure**: Maps complex JSON objects to GraphQL types (Company, Address, Geo)

### Implementation Details

```go
// HTTP client configuration
httpClient = httpclient.New(
    httpclient.WithBaseURL("https://jsonplaceholder.typicode.com"),
    httpclient.WithTimeout(5*time.Second),
)

// Example of REST API integration in resolvers
func (s *UsersService) QueryExternalUsers(ctx context.Context, req *service.QueryExternalUsersRequest) (*service.QueryExternalUsersResponse, error) {
    // Make HTTP request to external REST API
    httpResp, err := http.Get("https://jsonplaceholder.typicode.com/users")
    
    // Parse response and transform to GraphQL types
    // ...
}
```

## Joined Queries (GraphQL Federation)

```graphql
# Get all products with their associated user
query {
  products {
    id
    name
    price
    description
    user {
      id
      name
      email
      role
    }
  }
}

# Get all users with their associated products
query {
  users {
    id
    name
    email
    role
    products {
      id
      name
      price
      description
    }
  }
}

# Get all external users with their local user information through federation
query {
  externalUsers {
    id
    name
    email
    username
    localUser: _resolveReference(id: $id) {
      id
      role
    }
  }
}
```