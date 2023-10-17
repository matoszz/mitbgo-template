# go-template
Template repo for golang graphql apis

## Create Schema

1. Create a new schema by running the following command, replacing `<object>` with your object:
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/ent/schema <object>
```
1. Add fields to the object, for example: 
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

```
1. Run generate commands, this will use `entc` and `gqlgen` to generate the grapqhl api based on the defined schema. At any point, this should be able to be re-run and regenerate all files. 

```bash
make generate
```
1. Uncomment code in `cmd/serve.go`

If for any reason you want to remove all the generated code, you can run:
`make clean`

You will need `setopt extendedglob` enabled, and this was only tested with `zsh`, sorries. 

## References

### Ent

1. https://entgo.io/docs/code-gen/#use-entc-as-a-package