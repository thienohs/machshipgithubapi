package graph

import (
	"net/http"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	RetrieveUsersHandler func(w http.ResponseWriter, r *http.Request)
}
