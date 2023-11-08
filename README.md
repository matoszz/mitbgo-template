# go-template

Template repo for golang graphql apis

## Getting Started

This repo contains the basis for generating an opinionated Graph API using:

1. [ent](https://entgo.io/) - insane entity mapping tool, definitely not an ORM but kind of an ORM
1. [atlas](https://atlasgo.io/) - Schema generation and migration
1. [gqlgen](https://gqlgen.com/) - Code generation from schema definitions
1. [openfga](https://openfga.dev/) - Authorization

### Dependencies

Setup [Taskfile](https://taskfile.dev/installation/) by following the instructions and using one of the various convenient package managers or installation scripts. You can then simply run `task install` to load the associated dependencies. Nearly everything in this repository assumes you already have a local golang environment setup so this is not included. Please see the associated documentation.

## Usage

### Cleanup 

1. After cloning the repo, you will need to update all occurrences of `go-template` with your repo name. For convenience, a `task` command is included:
```bash
task clean-template
```

### Schema Generation with Ent

1. As the tooling suggests, this is schema driven api development so first up, is defining your schema
2. Create a new schema by running the following command, replacing `<object>` with your object:
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/ent/schema <object> 
```
For example, if you wanted to create a user, organization, and members schema you would run:
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/ent/schema User Organization Member 
```
3. This will generate a file per schema in `internal/ent/schema`
```bash
tree internal/ent/schema 

internal/ent/schema
└── user.go
└── organization.go
└── member.go
```
4. You will add your fields, edges, annotations, etc to this file for each schema. See the [ent schema def docs](https://entgo.io/docs/schema-def) for more details. 

5. For a simple User, this might look something like :
```go
// Fields of the User.
func (User) Fields() []ent.Field {
    return []ent.Field{
        field.Int("age").
            Positive(),
        field.String("name").
            Default("unknown"),
    }
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
```

5. Now that your schema is created, you want to generate your `ent.graphql`, this will contain all your graph `Input` types. The generate commands are setup in the `Taskfile` to make things easier:
```bash
task ent
```
6. This will create a `schema/ent.graphql` file
```
tree schema 

schema
├── ent.graphql
```
As well as generated ent files, including the `openapi.json`: 
```
internal/ent/generated
├── client.go
├── doc.go
├── ent.go
├── enttest
│   └── enttest.go
├── entviz.go
├── gql_collection.go
├── gql_edge.go
├── gql_node.go
├── gql_pagination.go
├── gql_transaction.go
├── gql_where_input.go
├── hook
│   └── hook.go
├── member
│   ├── member.go
│   └── where.go
├── member.go
├── member_create.go
├── member_delete.go
├── member_query.go
├── member_update.go
├── migrate
│   ├── migrate.go
│   └── schema.go
├── mutation.go
├── ogent
│   ├── oas_cfg_gen.go
│   ├── oas_client_gen.go
│   ├── oas_handlers_gen.go
│   ├── oas_interfaces_gen.go
│   ├── oas_json_gen.go
│   ├── oas_middleware_gen.go
│   ├── oas_parameters_gen.go
│   ├── oas_request_decoders_gen.go
│   ├── oas_request_encoders_gen.go
│   ├── oas_response_decoders_gen.go
│   ├── oas_response_encoders_gen.go
│   ├── oas_router_gen.go
│   ├── oas_schemas_gen.go
│   ├── oas_server_gen.go
│   ├── oas_unimplemented_gen.go
│   ├── oas_validators_gen.go
│   ├── ogent.go
│   └── responses.go
├── openapi.json
├── organization
│   ├── organization.go
│   └── where.go
├── organization.go
├── organization_create.go
├── organization_delete.go
├── organization_query.go
├── organization_update.go
├── predicate
│   └── predicate.go
├── runtime
│   └── runtime.go
├── runtime.go
├── schema-viz.html
├── tx.go
├── user
│   ├── user.go
│   └── where.go
├── user.go
├── user_create.go
├── user_delete.go
├── user_query.go
└── user_update.go
```
8. Now you will need to create a `graphql` file per schema that will handle CRUD operations, using the same example this would look like: 
```
tree schema 
schema
├── ent.graphql
└── user.graphql
└── organization.graphql
└── member.graphql
```
To have the files auto generated, use:
```bash
task graph
```

Now, the schema definitions should be ready to go. However, if at any point the schema needs to change, just rerun `make ent` and the ent.graphql and generated files should be updated. 

### Graph API Generation with gqlgen 

1. With the schemas ready, we can now generate the api code: 
```bash
task gqlgen
```
2. This will have created a new `internal/api` directory with a resolver per schema object
```
tree internal/api
internal/api
├── ent.resolvers.go
├── federation.go
├── gen_models.go
├── gen_server.go
├── resolver.go
└── user.resolvers.go
└── organization.resolvers.go
└── member.resolvers.go
```
3. In the resolvers, there will be stubbed out CRUD operations based on the grapqhl schemas. The business logic and permissions checks should go in here:
```go
// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}
```
4. Uncomment the code in `cmd/serve.go` 
5. Uncomment the code in `internal/api/resolver.go`

### Running locally

1. Now that all the code is there, test it using the playground:
```
make run-dev
```
2. Using the default config, you should be able to go to your browser of choice and see the playground: http://localhost:17608/playground
3. Via curl, `http://localhost:17608/query`


### Creating DB Migrations

1. Create DB Migrations with `atlas`:
```bash
task atlas:create
```
