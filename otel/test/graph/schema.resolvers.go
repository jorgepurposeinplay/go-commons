package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"
	"github.com/purposeinplay/go-commons/otel/test/graph/generated"
	"github.com/purposeinplay/go-commons/otel/test/graph/model"
)

// GetUser is the resolver for the getUser field.
func (r *queryResolver) GetUser(ctx context.Context, id string) (*model.User, error) {
	if r.GetUserIDFunc != nil {
		return r.GetUserIDFunc(ctx, id)
	}

	return &model.User{
		ID:   id,
		Name: "test",
	}, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
