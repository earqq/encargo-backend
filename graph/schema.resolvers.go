package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/earqq/encargo-backend/auth"
	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
	"gopkg.in/mgo.v2/bson"
)

func (r *mutationResolver) CreateCarrier(ctx context.Context, input model.NewCarrier) (*model.Carrier, error) {
	carriers := db.GetCollection("carriers")
	var carrier model.Carrier
	var fields = bson.M{}
	fields["$or"] = []bson.M{
		bson.M{"username": input.Username},
		bson.M{"phone": input.Phone}}
	if err := carriers.Find(fields).One(&carrier); err == nil {
		return &model.Carrier{}, errors.New("Nombre de usuario o Celular ya existe")
	}
	password, _ := HashPassword((input.Password))
	//Generar token
	var Token = auth.GenerateJWT(input.Username, "carrier")
	id := bson.NewObjectId()
	carriers.Insert(bson.M{
		"_id":              bson.ObjectId(id).Hex(),
		"name":             input.Name,
		"state_delivery":   0,
		"username":         input.Username,
		"token":            Token,
		"store_id":         input.StoreID,
		"password":         password,
		"global":           input.Global,
		"current_order_id": 0,
		"message_token":    input.MessageToken,
		"phone":            input.Phone,
		"updated_at":       time.Now().Local(),
	})

	err := carriers.Find(bson.M{"username": input.Username}).One(&carrier)
	if err != nil {
		return &model.Carrier{}, err
	}

	return &carrier, nil
}

func (r *mutationResolver) CreateStore(ctx context.Context, input model.NewStore) (*model.Store, error) {
	var stores model.Store
	storesDB := db.GetCollection("stores")
	var fields = bson.M{}
	fields["$or"] = []bson.M{
		bson.M{"username": input.Username},
		bson.M{"phone": input.Phone}}
	if err := storesDB.Find(fields).One(&stores); err == nil {
		return &model.Store{}, errors.New("Nombre de usuario o Celular ya existe")
	}
	if input.Ruc != nil {
		var fields = bson.M{}
		fields["ruc"] = input.Ruc
		if err := storesDB.Find(fields).One(&stores); err == nil {
			return &model.Store{}, errors.New("Ruc ya existe")
		}
	}
	password, _ := HashPassword((input.Password))
	//Generar token
	var Token = auth.GenerateJWT(input.Username, "store")
	id := bson.NewObjectId()
	storesDB.Insert(bson.M{
		"_id":         bson.ObjectId(id).Hex(),
		"name":        input.Name,
		"phone":       input.Phone,
		"username":    input.Username,
		"ruc":         input.Ruc,
		"token":       Token,
		"password":    password,
		"firebase_id": input.FirebaseID,
		"location":    input.Location,
	})
	if err := storesDB.Find(bson.M{"username": input.Username}).One(&stores); err != nil {
		return &model.Store{}, err
	}
	return &stores, nil
}

func (r *mutationResolver) UpdateCarrier(ctx context.Context, input model.UpdateCarrier) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	var carrier model.Carrier
	carriers := db.GetCollection("carriers")
	if err := carriers.Find(bson.M{"username": userContext.Username}).One(&carrier); err != nil {
		return &model.Carrier{}, err
	}

	var fields = bson.M{}

	update := false
	if input.MessageToken != nil && *input.MessageToken != "" {
		fields["message_token"] = input.MessageToken
		update = true

	}
	if input.Name != nil && *input.Name != "" {
		update = true
		fields["name"] = input.Name
	}
	if input.Global != nil {
		update = true
		fields["name"] = input.Global
	}
	if input.Password != nil && *input.Password != "" {
		update = true
		password, _ := HashPassword(*input.Password)
		fields["password"] = password
	}
	if input.State != nil {
		update = true
		fields["state"] = input.State
	}

	if !update {
		return &model.Carrier{}, errors.New("no fields present for updating data")
	}

	carriers.Update(bson.M{"username": userContext.Username}, bson.M{"$set": fields})
	carriers.Find(bson.M{"username": userContext.Username}).One(&carrier)
	return &carrier, nil
}

func (r *mutationResolver) DeleteStore(ctx context.Context) (*model.Store, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Store{}, errors.New("Acceso denegado")
	}
	var store model.Store
	storesBD := db.GetCollection("stores")
	if err := storesBD.Find(bson.M{"username": userContext.Username}).One(&store); err != nil {
		return &model.Store{}, errors.New("no existe este tienda")
	}
	if err := storesBD.Remove(bson.M{"username": userContext.Username}); err != nil {
		return &model.Store{}, errors.New("error al borrar tienda")
	}
	return &store, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, input model.NewOrder) (*model.Order, error) {
	var order model.Order
	var ordersBD = db.GetCollection("orders")
	var store model.Store
	var storeDB = db.GetCollection("stores")
	var ExitLocation model.Location
	id := bson.NewObjectId()
	loc := time.FixedZone("UTC-5", -5*60*60)
	t := time.Now().In(loc)
	if input.StoreRuc != nil {
		if err := storeDB.Find(bson.M{"ruc": input.StoreRuc}).One(&store); err != nil {
			return &model.Order{}, errors.New("No existe tienda con ese RUC")
		}
	}
	if input.StoreID != nil {
		if err := storeDB.Find(bson.M{"_id": input.StoreID}).One(&store); err != nil {
			return &model.Order{}, errors.New("No existe tienda con ese ID")
		}
	}
	ExitLocation.Latitude = store.Location.Latitude
	ExitLocation.Longitude = store.Location.Longitude
	ExitLocation.Address = store.Location.Address
	ExitLocation.Name = store.Location.Name

	ordersBD.Insert(bson.M{
		"_id":              bson.ObjectId(id).Hex(),
		"store_id":         store.ID,
		"price":            input.Price,
		"delivery_price":   input.DeliveryPrice,
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
		ordersDB.Update(bson.M{"_id": id}, bson.M{"$set": fields})
		ordersDB.Find(bson.M{"_id": id}).One(&order)
	}
	if input.State != nil {
		fields["state"] = *input.State
		update = true

		var carrierFields = bson.M{}
		carrierFields["state_delivery"] = *input.State
		if *input.State == 0 {
			fields["carrier_id"] = ""
		}
		if *input.State == 2 {
			fields["departure_date"] = t.Format("2006-01-02T15:04:05")
		}
		if *input.State == 3 {
			fields["delivery_date"] = t.Format("2006-01-02T15:04:05")
			carrierFields["state_delivery"] = "0"
		}
		if order.CarrierID != "" {
			carriersDB.Update(bson.M{"_id": order.CarrierID}, bson.M{"$set": carrierFields})
			var carriers []*model.Carrier
			if err := carriersDB.Find(bson.M{"global": 1}).All(&carriers); err != nil {
				return &model.Order{}, errors.New("No hay carriers")
			}
			r.Lock()
			for _, observer := range Observers {
				observer <- carriers
			}
			r.Unlock()
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

func (r *orderResolver) ArrivalLocation(ctx context.Context, obj *model.Order) (*model.Location, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) ActualLocation(ctx context.Context, obj *model.Order) (*model.Location, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *orderResolver) Carrier(ctx context.Context, obj *model.Order) (*model.Carrier, error) {
	var carrier model.Carrier
	var carriersDB = db.GetCollection("carriers")
	if err := carriersDB.Find(bson.M{"_id": obj.CarrierID}).One(&carrier); err != nil {
		return &model.Carrier{}, errors.New("Carrier relacionado al order no existe")
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

func (r *queryResolver) Carrier(ctx context.Context) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	var carrier model.Carrier
	carriersDB := db.GetCollection("carriers")
	if err := carriersDB.Find(bson.M{"username": userContext.Username}).One(&carrier); err != nil {
		return &model.Carrier{}, err
	}
	return &carrier, nil
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string, global *int) ([]*model.Carrier, error) {
	userContext := auth.ForContext(ctx)

	var carriers []*model.Carrier
	var fields = bson.M{}
	if search != nil {
		fields["name"] = bson.M{"$regex": *search, "$options": "i"}
	}
	if global != nil {
		fields["global"] = global
	}
	if userContext != nil && userContext.UserType == "store" {
		fields["store_id"] = userContext.ID
	}
	carriersDB := db.GetCollection("carriers")
	if limit != nil {
		carriersDB.Find(fields).Select(bson.M{"token": 0}).Limit(*limit).Sort("-updated_at").All(&carriers)
	} else {
		carriersDB.Find(fields).Select(bson.M{"token": 0}).Sort("-updated_at").All(&carriers)
	}

	return carriers, nil
}

func (r *queryResolver) LoginCarrier(ctx context.Context, username string, password string) (*model.Carrier, error) {
	var carrier model.Carrier
	var carrierDB = db.GetCollection("carriers")
	if err := carrierDB.Find(bson.M{"username": username}).One(&carrier); err != nil {
		return &model.Carrier{}, err
	}
	match := CheckPasswordHash(password, carrier.Password)
	if !match {
		return &model.Carrier{}, errors.New("Clave incorrecta")
	}
	return &carrier, nil
}

func (r *queryResolver) LoginStore(ctx context.Context, username string, password string) (*model.Store, error) {
	var store model.Store
	var storeDB = db.GetCollection("stores")
	if err := storeDB.Find(bson.M{"username": username}).One(&store); err != nil {
		return &model.Store{}, err
	}
	match := CheckPasswordHash(password, store.Password)
	if !match {
		return &model.Store{}, errors.New("Clave incorrecta")
	}
	return &store, nil
}

func (r *queryResolver) Store(ctx context.Context) (*model.Store, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Store{}, errors.New("Acceso denegado")
	}
	var store model.Store
	storeDB := db.GetCollection("stores")
	if err := storeDB.Find(bson.M{"username": userContext.Username}).One(&store); err != nil {
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
		storesBD.Find(fields).Limit(*limit).Select(bson.M{"token": 0}).Sort("-updated_at").All(&stores)

	} else {
		storesBD.Find(fields).Select(bson.M{"token": 0}).Sort("-updated_at").All(&stores)
	}
	return stores, nil
}

func (r *queryResolver) GetCarrierStats(ctx context.Context) (*model.CarrierStats, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.CarrierStats{}, errors.New("Acceso denegado")
	}
	var carrier model.Carrier
	var carrierDB = db.GetCollection("carriers")
	var ordersDB = db.GetCollection("orders")
	if err := carrierDB.Find(bson.M{"username": userContext.Username}).One(&carrier); err != nil {
		return &model.CarrierStats{}, errors.New("No existe carrier ")
	}
	//Obtener ordenes completada por repartidor
	var ordersCompleteBD []model.Order
	if err := ordersDB.Find(
		bson.M{"carrier_id": carrier.ID,
			"state": 3}).All(&ordersCompleteBD); err != nil {
		return &model.CarrierStats{}, errors.New("error aqui")
	}
	//OBtener pedidos entregados para sacar promedio de score
	var orders []model.Order
	if err := ordersDB.Find(
		bson.M{"carrier_id": carrier.ID,
			"experience.date": bson.M{"$ne": nil}}).All(&orders); err != nil {
		return &model.CarrierStats{}, errors.New("error aca")
	}
	ordersComplete := len(orders)
	if ordersComplete == 0 {
		ordersComplete = 1
	}
	ranking := 0.00
	for i := len(orders) - 1; i >= 0; i-- {
		ranking += float64(orders[i].Experience.Score)
	}
	//Sacar promedio
	average := ranking / float64(ordersComplete)

	var carrierStats = &model.CarrierStats{
		len(ordersCompleteBD),
		average,
	}
	return carrierStats, nil
}

func (r *queryResolver) Order(ctx context.Context, id string) (*model.Order, error) {
	var order model.Order
	var ordersDB = db.GetCollection("orders")
	if err := ordersDB.Find(bson.M{"_id": id}).One(&order); err != nil {
		return &model.Order{}, err
	}

	return &order, nil
}

func (r *queryResolver) Orders(ctx context.Context, input model.FilterOptions) ([]*model.Order, error) {
	userContext := auth.ForContext(ctx)
	var orders []*model.Order
	var fields = bson.M{}
	var orArray = []bson.M{}
	var ordersDB = db.GetCollection("orders")
	if userContext != nil && userContext.UserType == "store" {
		fields["store_id"] = userContext.ID
	}
	if userContext != nil && userContext.UserType == "carrier" {
		fields["carrier_id"] = userContext.ID
	}
	if input.State != nil {
		fields["state"] = input.State
	}
	if input.CarrierID != nil {
		fields["carrier_id"] = input.CarrierID
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
	if len(orArray) > 0 {
		fields["$or"] = orArray
	}

	ordersDB.Find(fields).Limit(input.Limit).Sort("-updated_at").All(&orders)

	return orders, nil
}

func (r *subscriptionResolver) CarriersAvailable(ctx context.Context) (<-chan []*model.Carrier, error) {
	id := RandStringRunes(8)
	event := make(chan []*model.Carrier, 1)

	go func() {
		<-ctx.Done()
		r.Lock()
		delete(Observers, id)
		r.Unlock()
	}()
	r.Lock()
	Observers[id] = event
	r.Unlock()
	return event, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *subscriptionResolver) Order(ctx context.Context) (<-chan *model.Order, error) {
	panic(fmt.Errorf("not implemented"))
}
