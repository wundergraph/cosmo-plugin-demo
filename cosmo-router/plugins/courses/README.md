# Courses Plugin

A Cosmo Router plugin for course management and enrollment. This plugin provides a GraphQL API for managing courses, instructors, and student enrollments using gRPC service implementation.

## Features

- Course creation and management
- Instructor course assignments
- In-memory data storage with mock data

## Testing

The tests validate all functionality of the plugin including:

- Course querying (all courses and by ID)
- User lookup by ID (instructors and students)
- Course creation
- Student enrollment in courses

### Running Tests

```bash
cd cosmo-router/plugins/courses
make test
```

## Directory Structure

```
courses/
├── bin/                   # Compiled plugin binaries
├── generated/             # Auto-generated gRPC code from schema
│   ├── service.pb.js
│   ├── service_grpc_pb.js
│   └── service.proto
├── src/                   # Source code
│   ├── plugin.ts          # Plugin implementation
│   ├── plugin.test.ts     # Integration tests
│   ├── plugin-server.ts   # gRPC server setup
│   └── schema.graphql     # GraphQL schema definition
└── package.json           # Dependencies
```

## GraphQL API

### Queries

- `courses`: List all courses
- `course(id: ID!)`: Get a course by ID

### Mutations

- `createCourse(input: CourseInput!)`: Create a new course
- `enrollUser(userId: ID!, courseId: ID!)`: Enroll a user in a course


## Example GraphQL Operations

```graphql
# Get all courses
query {
  courses {
    id
    title
    description
    instructor {
      id
    }
    durationHours
    published
  }
}

# Get a specific course
query {
  course(id: "course-1") {
    id
    title
    description
    instructor {
      id
    }
    durationHours
    published
  }
}

# Create a new course
mutation {
  createCourse(input: {
    title: "Advanced GraphQL"
    description: "Master GraphQL concepts"
    instructorId: "1"
    durationHours: 30
    published: true
  }) {
    id
    title
    description
    instructor {
      id
    }
  }
}

# Enroll a user in a course
mutation {
  enrollUser(userId: "user-1", courseId: "course-1") {
    id
    user {
      id
    }
    course {
      id
      title
    }
    progress
  }
}
```