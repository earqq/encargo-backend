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
		"store_id":         input.StoreID,
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

func (r *mutationResolver) UpdateCarrier(ctx context.Context, id string, input model.UpdateCarrier) (*model.Carrier, error) {
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
	var stores model.Store
	storesBD := db.GetCollection("stores")
	id := bson.NewObjectId()
	storesBD.Insert(bson.M{
		"_id":         bson.ObjectId(id).Hex(),
		"name":        input.Name,
		"phone":       input.Phone,
		"username":    input.Username,
		"password":    input.Password,
		"firebase_id": input.FirebaseID,
		"location":    input.Location,
	})
	if err := storesBD.Find(bson.M{"_id": bson.ObjectId(id).Hex()}).One(&stores); err != nil {
		return &model.Store{}, err
	}

	return &stores, nil
}

func (r *mutationResolver) DeleteStore(ctx context.Context, id string) (*model.Store, error) {
	var stores model.Store
	storesBD := db.GetCollection("stores")
	if err := storesBD.Find(bson.M{"_id": id}).One(&stores); err != nil {
		return &model.Store{}, errors.New("no existe este negocio")
	}
	if err := storesBD.Remove(bson.M{"_id": id}); err != nil {
		return &model.Store{}, errors.New("error al borrar negocio")
	}
	return &stores, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.NewOrder) (*model.Order, error) {
	fmt.Println("store")
	fmt.Println(input.StoreID)
	var order model.Order
	var ordersBD = db.GetCollection("orders")
	var store model.Store
	var storeDB = db.GetCollection("stores")
	var ExitLocation model.Location
	id := bson.NewObjectId()
	loc := time.FixedZone("UTC-5", -5*60*60)
	t := time.Now().In(loc)
	if err := storeDB.Find(bson.M{"_id": input.StoreID}).One(&store); err != nil {
		return &model.Order{}, errors.New("No existe tienda")
	}
	ExitLocation.Latitude = store.Location.Latitude
	ExitLocation.Longitude = store.Location.Longitude
	ExitLocation.Address = store.Location.Address
	ExitLocation.Name = store.Location.Name

	ordersBD.Insert(bson.M{
		"_id":              bson.ObjectId(id).Hex(),
		"store_id":         input.StoreID,
		"price":            input.Price,
		"date":             t.Format("2006-01-02T15:04:0"),
		"state":            0,
		"client_phone":     input.ClientPhone,
		"detail":           input.Detail,
		"client_name":      input.ClientName,
		"exit_location":    ExitLocation,
		"arrival_location": input.ArrivalLocation,
	})
	if err := ordersBD.Find(bson.M{"_id": bson.ObjectId(id).Hex()}).One(&order); err != nil {
		return &model.Order{}, err
	}

	return &order, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, id string, input model.UpdateOrder) (*model.Order, error) {
	var order model.Order
	var carriersDB = db.GetCollection("carriers")
	var ordersDB = db.GetCollection("orders")
	if err := ordersDB.Find(bson.M{"_id": id}).One(&order); err != nil {
		return &model.Order{}, errors.New("No existe order")
	}

	var fields = bson.M{}
	loc := time.FixedZone("UTC-5", -5*60*60)
	t := time.Now().In(loc)
	update := false
	if input.CarrierID != nil {
		var carrier model.Carrier
		if err := carriersDB.Find(bson.M{"_id": input.CarrierID}).One(&carrier); err != nil {
			return &model.Order{}, errors.New("No existe carrier")
		}
		fields["carrier_id"] = *input.CarrierID
		update = true
		ordersDB.Update(bson.M{"id": id}, bson.M{"$set": fields})
		ordersDB.Find(bson.M{"id": id}).One(&order)
	}
	if input.State != nil {
		fields["state"] = *input.State
		update = true
		var carrier model.Carrier
		if err := carriersDB.Find(bson.M{"_id": input.CarrierID}).One(&carrier); err != nil {
			return &model.Order{}, errors.New("No se encuentra carrier con ese ID")
		}
		if *input.State == 0 {
			fields["carrier_id"] = ""
		}
		if *input.State == 2 {
			fields["departure_date"] = t.Format("2006-01-02T15:04:05")
		}
		if *input.State == 3 {
			fields["delivery_date"] = t.Format("2006-01-02T15:04:05")
		}
	}
	if input.Score != nil {
		update = true
		fields["experience.score"] = *input.Score
		fields["experience.date"] = t.Format("2006-01-02T15:04:05")
	}
	if input.ScoreDescription != nil {
		update = true
		fields["experience.description"] = *input.ScoreDescription
	}
	if !update {
		return &model.Order{}, errors.New("No hay ningun campo para actualizar")
	}
	fields["updated_at"] = time.Now().Local()
	ordersDB.Update(bson.M{"_id": id}, bson.M{"$set": fields})

	ordersDB.Find(bson.M{"_id": id}).One(&order)

	return &order, nil
}

func (r *orderResolver) Carrier(ctx context.Context, obj *model.Order) (*model.Carrier, error) {
	var carrier model.Carrier
	var carriersDB = db.GetCollection("carriers")
	if err := carriersDB.Find(bson.M{"_id": obj.CarrierID}).One(&carrier); err != nil {
		return &model.Carrier{}, nil
	}
	return &carrier, nil
}

func (r *orderResolver) Store(ctx context.Context, obj *model.Order) (*model.Store, error) {
	var store model.Store
	var storesDB = db.GetCollection("stores")
	if err := storesDB.Find(bson.M{"_id": obj.StoreID}).One(&store); err != nil {
		return &model.Store{}, errors.New("Store relacionado al order no existe")
	}
	return &store, nil
}

func (r *queryResolver) Carrier(ctx context.Context, id string) (*model.Carrier, error) {
	var user model.Carrier
	carriersDB := db.GetCollection("carriers")
	if err := carriersDB.Find(bson.M{"_id": id}).One(&user); err != nil {
		return &model.Carrier{}, err
	}
	return &user, nil
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string, storeID string) ([]*model.Carrier, error) {
	var carriers []*model.Carrier
	var fields = bson.M{}
	fields["store_id"] = storeID
	carriersDB := db.GetCollection("carriers")
	if limit != nil {
		carriersDB.Find(fields).Limit(*limit).Sort("-updated_at").All(&carriers)

	} else {
		carriersDB.Find(fields).Sort("-updated_at").All(&carriers)
	}

	return carriers, nil
}

func (r *queryResolver) LoginCarrier(ctx context.Context, username string, password string) (*model.Carrier, error) {
	var user model.Carrier
	var carrierDB = db.GetCollection("carriers")
	if err := carrierDB.Find(bson.M{"username": username, "password": password}).One(&user); err != nil {
		return &model.Carrier{}, err
	}

	return &user, nil
}

func (r *queryResolver) LoginStore(ctx context.Context, username string, password string) (*model.Store, error) {
	var user model.Store
	var storeDB = db.GetCollection("stores")
	if err := storeDB.Find(bson.M{"username": username, "password": password}).One(&user); err != nil {
		return &model.Store{}, err
	}

	return &user, nil
}

func (r *queryResolver) GetCarrierStats(ctx context.Context, carrierID string) (*model.CarrierStats, error) {
	var carrier model.Carrier
	var carrierDB = db.GetCollection("carriers")
	var ordersDB = db.GetCollection("orders")
	if err := carrierDB.Find(bson.M{"public_id": carrierID}).One(&carrier); err != nil {
		return &model.CarrierStats{}, errors.New("No existe carrier ")
	}
	var carrierStats *model.CarrierStats
	var ordersCompleteBD []model.Order
	if err := ordersDB.Find(
		bson.M{"carrier_id": carrierID,
			"state": 3}).All(&ordersCompleteBD); err != nil {
		return &model.CarrierStats{}, err
	}
	carrierStats.Orders = len(ordersCompleteBD)
	var orders []model.Order
	if err := ordersDB.Find(
		bson.M{"carrier_id": carrierID,
			"experience.date": bson.M{"$ne": nil}}).All(&orders); err != nil {
		return &model.CarrierStats{}, err
	}

	ordersComplete := len(orders)
	if ordersComplete == 0 {
		ordersComplete = 1
	}
	ranking := 0.00
	for i := len(orders) - 1; i >= 0; i-- {
		ranking += float64(orders[i].Experience.Score)
	}
	average := ranking / float64(ordersComplete)

	carrierStats.Ranking = average
	return carrierStats, nil
}

func (r *queryResolver) Store(ctx context.Context, id string) (*model.Store, error) {
	var store model.Store
	storeDB := db.GetCollection("stores")
	if err := storeDB.Find(bson.M{"_id": id}).One(&store); err != nil {
		return &model.Store{}, err
	}
	return &store, nil
}

func (r *queryResolver) Stores(ctx context.Context, limit *int, search *string) ([]*model.Store, error) {
	var stores []*model.Store
	var fields = bson.M{}
	storesBD := db.GetCollection("stores")
	if search != nil {
		fields["name"] = bson.M{"$regex": *search, "$options": "i"}
	}
	if limit != nil {
		storesBD.Find(fields).Limit(*limit).Sort("-updated_at").All(&stores)

	} else {
		storesBD.Find(fields).Sort("-updated_at").All(&stores)
	}
	return stores, nil
}

func (r *queryResolver) Orders(ctx context.Context, input model.FilterOptions, storeID string) ([]*model.Order, error) {
	var orders []*model.Order
	var fields = bson.M{}
	var orArray = []bson.M{}
	var ordersDB = db.GetCollection("orders")
	fields["store_id"] = storeID
	if input.ID != nil && *input.ID != "" {
		fields["_id"] = input.ID
	}
	if input.State != nil {
		fields["state"] = input.State
	}
	if input.State1 != nil {
		orArray = append(orArray, bson.M{"state": input.State1})
		orArray = append(orArray, bson.M{"state": input.State2})
	}
	if input.Search != nil {
		orArray = append(orArray, bson.M{"description": bson.M{"$regex": *input.Search, "$options": "i"}},
			bson.M{"client_name": bson.M{"$regex": *input.Search, "$options": "i"}},
			bson.M{"client_phone": bson.M{"$regex": *input.Search, "$options": "i"}})
	}
	if input.CarrierID != nil {
		fields["carrier_id"] = input.CarrierID
	}
	if len(orArray) > 0 {
		fields["$or"] = orArray
	}

	ordersDB.Find(fields).Limit(input.Limit).Sort("-updated_at").All(&orders)

	return orders, nil
}

func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	var ordersDB = db.GetCollection("orders")
	if err := ordersDB.Find(bson.M{"_id": id}).One(&order); err != nil {
		return &model.Order{}, err
	}

	return &order, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
