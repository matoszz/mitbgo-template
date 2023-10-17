# go-template
Template repo for golang graphql apis

## Create Schema

1. Create a new schema by running the following command, replacing `<object>` with your object:
```bash
go run -mod=mod entgo.io/ent/cmd/ent new --target internal/ent/schema <object>
```
1. Add fields to the object
1. Run generate commands:
```bash
make generate
```