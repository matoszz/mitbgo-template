# go-template

Template repo for golang graphql apis

## Getting Started

This repo contains the basis for generating an opinionated Graph API using:

1. [ent](https://entgo.io/) - insane entity mapping tool, definitely not an ORM but kind of an ORM
1. [atlas](https://atlasgo.io/) - Schema generation and migration
1. [gqlgen](https://gqlgen.com/) - Code generation from schema definitions
1. [gqlgenc](https://github.com/Yamashou/gqlgenc) - client building utilities with GraphQL
1. [openfga](https://openfga.dev/) - Authorization
1. [echo](https://echo.labstack.com/) - High performance, extensible, minimalist Go web framework
1. [koanf](github.com/knadh/koanf) - configuration management
1. [viper](https://github.com/spf13/viper) - command line flags / management

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
1. Create a new schema by running the following command, replacing `<object>` with your object:
    ```bash
    task newschema -- <object> 
    ```
    For example, if you wanted to create a user, organization, and members schema you would run:
    ```bash
    task newschema -- User Organization Member 
    ```
1. This will generate a file per schema in `internal/ent/schema`
    ```bash
    tree internal/ent/schema 

    internal/ent/schema
    └── user.go
    └── organization.go
    └── member.go
    ```
1. You will add your fields, edges, annotations, etc to this file for each schema. See the [ent schema def docs](https://entgo.io/docs/schema-def) for more details. 

1. For a simple User, this might look something like :
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

1. Now that your schema is created, you want to generate your `ent.graphql`, this will contain all your graph `Input` types. The generate commands are setup in the `Taskfile` to make things easier:
    ```bash
    task generate
    ```
1. This will create a `schema/ent.graphql` file
    ```
    tree schema 

    schema
    ├── ent.graphql
    └── todo.graphql
    ```
    As well as generated ent files, including the `openapi.json`: 
   
1. This will have created a new `internal/graphapi` directory with a resolver per schema object
    ```
    internal/graphapi
    ├── doc.go
    ├── ent.resolvers.go
    ├── gen_models.go
    ├── gen_server.go
    ├── helpers.go
    ├── resolver.go
    └── todo.resolvers.go
    ```
1. In the resolvers, there will be stubbed out CRUD operations based on the grapqhl schemas. The business logic and permissions checks should go in here:
    ```go
    // CreateTodo is the resolver for the createTodo field.
    func (r *mutationResolver) CreateTodo(ctx context.Context, input CreateUserInput) (*User, error) {
        panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
    }
    ```

### Dependencies

The Datum ecosystem has additional dependencies that were not included in the vanilla template because they will not be required for all services. These include things such as:

1. [Email Manager](https://github.com/datumforge/datum/tree/main/pkg/utils/emails)
1. [Task Manager](https://github.com/datumforge/datum/tree/main/pkg/utils/marionette)
1. [Analytics Manager](https://github.com/datumforge/datum/tree/main/pkg/analytics)

To add these to your project, refer to the implementation in [Datum](https://github.com/datumforge/datum)

1. [Config Setup](https://github.com/datumforge/datum/blob/main/internal/httpserve/serveropts/option.go#L238)
1. [`entc` Setup](https://github.com/datumforge/datum/blob/main/internal/ent/entc.go#L123)
1. [Server Setup](https://github.com/datumforge/datum/blob/main/cmd/serve.go#L73-L80)


### Running locally

1. Now that all the code is there, test it using the playground:
    ```
    task run-dev
    ```
1. Using the default config, you should be able to go to your browser of choice and see the playground: http://localhost:1337/playground
1. Via curl, `http://localhost:1337/query`


### Creating DB Migrations

1. Create DB Migrations with `atlas`:
    ```bash
    task atlas:create
    ```
