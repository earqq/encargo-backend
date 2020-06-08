package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
	"gopkg.in/mgo.v2/bson"
)

func (r *mutationResolver) CreateCarrier(ctx context.Context, input model.NewCarrier) (*model.Carrier, error) {
	var user model.Carrier
	var fields = bson.M{}
	fmt.Println("si viene")
	fields["$or"] = []bson.M{
		bson.M{"username": input.Username},
		bson.M{"phone": input.Phone}}
	fmt.Println("si viene2")
	if err := r.carriers.Find(fields).One(&user); err == nil {
		return &model.Carrier{}, errors.New("Nombre de usuario o Celular ya existe")
	}
	fmt.Println("si viene3")

	r.carriers.Insert(bson.M{
		"public_id":        input.PublicID,
		"name":             input.Name,
		"birthdate":        input.Birthdate,
		"state_delivery":   0,
		"username":         input.Username,
		"password":         input.Password,
		"current_order_id": 0,
		"message_token":    input.MessageToken,
		"phone":            input.Phone,
		"updated_at":       time.Now().Local(),
	})
	fmt.Println("si viene4")

	err := r.carriers.Find(bson.M{"username": input.Username}).One(&user)
	if err != nil {
		return &model.Carrier{}, err
	}

	fmt.Println("si viene5")
	return &user, nil
}

func (r *mutationResolver) UpdateCarrier(ctx context.Context, publicID string, input *model.UpdateCarrier) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateStore(ctx context.Context, input model.NewStore) (*model.Store, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteStore(ctx context.Context, publicID string) (*model.Store, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.NewOrder) (*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, publicID string, input model.UpdateOrderInput) (*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Detail(ctx context.Context, obj *model.Order) ([]*model.OrderDetail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Carrier(ctx context.Context, obj *model.Order) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderDetailResolver) Amount(ctx context.Context, obj *model.OrderDetail) (*float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderDetailResolver) Price(ctx context.Context, obj *model.OrderDetail) (*float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderDetailResolver) Description(ctx context.Context, obj *model.OrderDetail) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carrier(ctx context.Context, publicID string) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string) ([]*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carriers2(ctx context.Context, limit *int, search *string) ([]*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carriers3(ctx context.Context, limit *int, search *string) ([]*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Login(ctx context.Context, username string, password string) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetCarrierStats(ctx context.Context, carrierPublicID string) (*model.CarrierStats, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Stores(ctx context.Context, limit *int, search *string) ([]*model.Store, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Orders(ctx context.Context, input model.FilterOptions) ([]*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Order(ctx context.Context, publicID string) (*model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *storeResolver) Location(ctx context.Context, obj *model.Store) (*model.Location, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

// OrderDetail returns generated.OrderDetailResolver implementation.
func (r *Resolver) OrderDetail() generated.OrderDetailResolver { return &orderDetailResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Store returns generated.StoreResolver implementation.
func (r *Resolver) Store() generated.StoreResolver { return &storeResolver{r} }

type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type orderDetailResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type storeResolver struct{ *Resolver }
