# OpenFGA Tutorial: Fine-Grained Authorization with Zanzibar

This repository demonstrates the concepts of Fine-Grained Authorization (FGA) using OpenFGA, Google's Zanzibar-inspired authorization system. The project implements a Document Management System scenario to showcase how relationship-based access control works in practice.

## Overview

OpenFGA is an open-source authorization engine that implements Google's Zanzibar model for fine-grained authorization. This tutorial demonstrates:

- **Relationship-based Access Control**: Define complex authorization relationships between users, teams, departments, and documents
- **Hierarchical Permissions**: Implement nested organizational structures with inherited permissions
- **Authorization Modeling**: Create and deploy authorization models using OpenFGA's DSL
- **Runtime Authorization Checks**: Perform real-time permission checks for access control decisions

## Project Structure

```
├── main.go                 # Main application demonstrating OpenFGA concepts
├── docker-compose.yaml     # OpenFGA server and PostgreSQL setup
├── go.mod                  # Go module dependencies
└── go.sum                  # Go module checksums
```

## Authorization Model

The project implements a Document Management System with the following entities and relationships:

### Entity Types

- **User**: Individual users in the system
- **Team**: Groups of users with shared responsibilities
- **Department**: Organizational units containing multiple teams
- **Document**: Files with ownership and editing permissions

### Relationships

- **Team Members**: Users can be members of teams
- **Department Structure**: Teams belong to departments, and department membership is inherited from team membership
- **Document Ownership**: Users can own documents (owners have full access)
- **Document Editing**: Users, team members, or department members can have editing permissions

### Permission Inheritance

The model demonstrates sophisticated permission inheritance:
- Document owners automatically have editing permissions
- Team members inherit permissions granted to their team
- Department members inherit permissions through their team memberships

## Getting Started

### Prerequisites

- Go 1.24+ installed
- Docker and Docker Compose
- Basic understanding of authorization concepts

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd fga-tutorial-1
   ```

2. **Start OpenFGA server**
   ```bash
   docker-compose up -d
   ```
   This starts:
   - PostgreSQL database (port 5432)
   - OpenFGA server (port 8080 for HTTP, 8081 for gRPC)
   - OpenFGA Playground (port 3000)

3. **Install Go dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the tutorial**
   ```bash
   go run main.go
   ```

## What the Code Demonstrates

The `main.go` file walks through a complete OpenFGA workflow:

### 1. Store Creation
Creates a new authorization store for the Document Management System.

### 2. Authorization Model Definition
Defines the relationship model with four entity types and their relationships, including complex inheritance patterns.

### 3. Relationship Data Population
Creates sample relationships:
- Alice owns `doc-001`
- Engineering team members can edit `doc-001`
- Bob is a member of the engineering team
- Engineering team belongs to product department
- Product department members can edit `doc-002`

### 4. Authorization Checks
Performs real-time checks to determine:
- Can Bob edit `doc-001`? ✅ (via team membership)
- Can Bob edit `doc-002`? ✅ (via department membership through team)
- Can Bob edit `doc-003`? ❌ (no relationship exists)

## Key Concepts Demonstrated

### Relationship-Based Access Control (ReBAC)
Instead of traditional role-based access control, the system uses relationships between entities to determine permissions.

### Transitive Relationships
Shows how permissions can be inherited through multiple levels of relationships (user → team → department → document).

### Union Permissions
Demonstrates how multiple permission paths can lead to the same access (direct assignment OR ownership).

### Computed Relationships
Illustrates how department membership is computed from team relationships using `tupleToUserset`.

## OpenFGA Playground

Access the OpenFGA Playground at `http://localhost:3000` to:
- Visualize the authorization model
- Test authorization queries interactively
- Explore relationship graphs
- Debug permission inheritance

## Learning Resources

- [OpenFGA Documentation](https://openfga.dev/docs)
- [Zanzibar Paper](https://research.google/pubs/pub48190/)
- [Fine-Grained Authorization Concepts](https://openfga.dev/docs/concepts)

## Next Steps

To extend this tutorial:
- Add more complex organizational hierarchies
- Implement time-based permissions
- Add audit logging for authorization decisions
- Integrate with your application's authentication system
- Explore batch authorization checks for performance optimization

## License

This project is provided as an educational resource for learning OpenFGA and fine-grained authorization concepts.
