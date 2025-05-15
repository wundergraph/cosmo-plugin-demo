import { ApolloServer } from '@apollo/server';
import { startStandaloneServer } from '@apollo/server/standalone';
import { buildSubgraphSchema } from '@apollo/subgraph';
import { parse } from 'graphql';
import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

// Get current directory
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Sample product data
const products = [
  { id: '1', name: 'Laptop', price: 999.99, description: 'Powerful laptop for developers', userId: '1' },
  { id: '2', name: 'Smartphone', price: 699.99, description: 'Latest smartphone model', userId: '1' },
  { id: '3', name: 'Tablet', price: 399.99, description: 'Portable tablet device', userId: '2' },
  { id: '4', name: 'Headphones', price: 149.99, description: 'Noise-canceling headphones', userId: '3' },
  { id: '5', name: 'Monitor', price: 299.99, description: '4K HD Monitor', userId: '2' },
];

// Read schema from file
const schemaPath = path.join(__dirname, 'schema.graphql');
const typeDefs = parse(fs.readFileSync(schemaPath, 'utf-8'));
console.log(`📄 Schema read from ${schemaPath}`);

// Resolvers
const resolvers = {
  Query: {
    products: () => products,
    product: (_, { id }) => products.find(product => product.id === id),
  },
  Product: {
    // This resolver sets up the connection from Product to User
    user(product) {
      // We just need to return an object with __typename and the id
      // The federation gateway will handle fetching the actual user fields
      return { __typename: "User", id: product.userId };
    },
    // Needed for federation reference resolution
    __resolveReference(reference) {
      return products.find(product => product.id === reference.id);
    }
  },
  User: {
    // This resolver adds products to the User type
    products(user) {
      return products.filter(product => product.userId === user.id);
    },
    // The reference resolver doesn't need to find the full user, just pass it along
    __resolveReference(reference) {
      // This just says "yes, I recognize this reference" and returns it for processing
      return reference;
    }
  },
};

// Create Apollo Server with the federation subgraph schema
const schema = buildSubgraphSchema({ typeDefs, resolvers });
const server = new ApolloServer({
  schema,
});

// Start the server
const { url } = await startStandaloneServer(server, {
  listen: { port: 3011 },
});

console.log(`🚀 Products subgraph ready at ${url}`); 