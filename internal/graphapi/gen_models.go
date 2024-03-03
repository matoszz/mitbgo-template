// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package graphapi

import (
	"github.com/datumforge/go-template/internal/ent/generated"
)

// Return response for createTodo mutation
type TodoCreatePayload struct {
	// Created todo
	Todo *generated.Todo `json:"todo"`
}

// Return response for deleteTodo mutation
type TodoDeletePayload struct {
	// Deleted todo ID
	DeletedID string `json:"deletedID"`
}

// Return response for updateTodo mutation
type TodoUpdatePayload struct {
	// Updated todo
	Todo *generated.Todo `json:"todo"`
}
