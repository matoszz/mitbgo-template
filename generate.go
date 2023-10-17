package main

//go:generate go run -mod=mod ./internal/ent/entc.go
//go:generate go run -mod=mod github.com/99designs/gqlgen
//go:generate go run -mod=mod ./gen_schema.go
//go:generate go run -mod=mod github.com/Yamashou/gqlgenc
