package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
	"gopkg.in/mgo.v2/bson"
)

func (r *carrierResolver) StoreID(ctx context.Context, obj *model.Carrier) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateCarrier(ctx context.Context, input model.NewCarrier) (*model.Carrier, error) {
	carriers := db.GetCollection("carriers")
	var user model.Carrier
	var fields = bson.M{}
	fields["$or"] = []bson.M{
		bson.M{"username": input.Username},
		bson.M{"phone": input.Phone}}
	if err := carriers.Find(fields).One(&user); err == nil {
		return &model.Carrier{}, errors.New("Nombre de usuario o Celular ya existe")
	}
	id := bson.NewObjectId()
	carriers.Insert(bson.M{
		"_id":              bson.ObjectId(id).Hex(),
		"name":             input.Name,
		"state_delivery":   0,
		"username":         input.Username,
		"password":         input.Password,
		"current_order_id": 0,
		"message_token":    input.MessageToken,
		"phone":            input.Phone,
		"updated_at":       time.Now().Local(),
	})

	err := carriers.Find(bson.M{"username": input.Username}).One(&user)
	if err != nil {
		return &model.Carrier{}, err
	}

	return &user, nil
}

func (r *mutationResolver) UpdateCarrier(ctx context.Context, id string, input *model.UpdateCarrier) (*model.Carrier, error) {
	var user model.Carrier
	carriers := db.GetCollection("carriers")
	if err := carriers.Find(bson.M{"_id": id}).One(&user); err != nil {
		return &model.Carrier{}, err
	}

	var fields = bson.M{}

	update := false
	if input.MessageToken != nil && *input.MessageToken != "" {
		fields["message_token"] = input.MessageToken
		update = true

	}
	if input.Username != nil && *input.Username != "" {
		update = true
		fields["username"] = input.Username
	}
	if input.Name != nil && *input.Name != "" {
		update = true
		fields["name"] = input.Name
	}
	if input.Password != nil && *input.Password != "" {
		update = true
		fields["password"] = input.Password
	}
	if input.State != nil {
		update = true
		fields["state"] = input.State
	}
	if input.Phone != nil && *input.Phone != "" {
		update = true
		fields["phone"] = input.Phone
	}

	if !update {
		return &model.Carrier{}, errors.New("no fields present for updating data")
	}

	carriers.Update(bson.M{"_id": id}, bson.M{"$set": fields})
	carriers.Find(bson.M{"_id": id}).One(&user)

	return &user, nil
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

func (r *orderResolver) PublicID(ctx context.Context, obj *model.Order) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Detail(ctx context.Context, obj *model.Order) ([]*model.OrderDetail, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Carrier(ctx context.Context, obj *model.Order) (*model.Carrier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Store(ctx context.Context, obj *model.Order) (*model.Store, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Carrier(ctx context.Context, id string) (*model.Carrier, error) {
	var user model.Carrier
	carriersDB := db.GetCollection("carriers")
	if err := carriersDB.Find(bson.M{"_id": id}).One(&user); err != nil {
		return &model.Carrier{}, err
	}
	user.ID = bson.ObjectId(user.ID).Hex()
	return &user, nil
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string) ([]*model.Carrier, error) {
	var carriers []*model.Carrier
	var fields = bson.M{}

	carriersDB := db.GetCollection("carriers")
	if limit != nil {
		carriersDB.Find(fields).Limit(*limit).Sort("-updated_at").All(&carriers)

	} else {
		carriersDB.Find(fields).Sort("-updated_at").All(&carriers)
	}

	return carriers, nil
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

// Carrier returns generated.CarrierResolver implementation.
func (r *Resolver) Carrier() generated.CarrierResolver { return &carrierResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Store returns generated.StoreResolver implementation.
func (r *Resolver) Store() generated.StoreResolver { return &storeResolver{r} }

type carrierResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type storeResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *carrierResolver) ID(ctx context.Context, obj *model.Carrier) (*int, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *orderResolver) ID(ctx context.Context, obj *model.Order) (int, error) {
	panic(fmt.Errorf("not implemented"))
}
func (r *storeResolver) ID(ctx context.Context, obj *model.Store) (int, error) {
	panic(fmt.Errorf("not implemented"))
}
