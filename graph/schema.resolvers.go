package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
)

func (r *mutationResolver) CreateCarrier(ctx context.Context, input model.NewCarrier) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateCarrier(ctx context.Context, publicID string, input *model.UpdateCarrier) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carrier(ctx context.Context, publicID string) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string) ([]*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Login(ctx context.Context, username string, password string) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCarrierStats(ctx context.Context, carrierPublicID string) (*model.CarrierStats, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
