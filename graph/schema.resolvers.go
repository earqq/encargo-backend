package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"github.com/earqq/encargo-backend/auth"
	"github.com/earqq/encargo-backend/db"
	"github.com/earqq/encargo-backend/graph/generated"
	"github.com/earqq/encargo-backend/graph/model"
	"gopkg.in/mgo.v2/bson"
)

func (r *carrierResolver) ActualLocation(ctx context.Context, obj *model.Carrier) (*model.Location, error) {
	return &obj.ActualLocation, nil
}

func (r *carrierResolver) Order(ctx context.Context, obj *model.Carrier) (*model.Order, error) {
	var order model.Order
	if err := r.orders.Find(bson.M{"_id": obj.CurrentOrderID}).One(&order); err != nil {
		return &model.Order{}, nil
	}
	return &order, nil
}

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

func (r *mutationResolver) UpdateCarrier(ctx context.Context, id *string, input model.UpdateCarrier) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	var carrier model.Carrier
	if userContext.UserType == "carrier" {
		if err := r.carriers.Find(bson.M{"username": userContext.Username}).One(&carrier); err != nil {
			return &model.Carrier{}, err
		}
	} else {
		if err := r.carriers.Find(bson.M{"_id": id, "store_id": userContext.ID}).One(&carrier); err != nil {
			return &model.Carrier{}, errors.New("No existe repartidor en la tienda")
		}
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
		fields["global"] = input.Global
	}
	if input.Password != nil && *input.Password != "" {
		update = true
		password, _ := HashPassword(*input.Password)
		fields["password"] = password
	}
	if input.StateDelivery != nil {
		update = true
		fields["state_delivery"] = input.StateDelivery
	}
	if !update {
		return &model.Carrier{}, errors.New("No hay campos por actualizar")
	}
	r.carriers.Update(bson.M{"_id": carrier.ID}, bson.M{"$set": fields})
	r.carriers.Find(bson.M{"_id": carrier.ID}).One(&carrier)
	r.Lock() //Enviando info a tienda sobre carrier actualizado
	topic := r.storeCarriersTopics[carrier.StoreID]
	if topic != nil {
		for _, observer := range topic.Observers {
			observer <- &carrier
		}
	}
	r.Unlock()

	return &carrier, nil
}

func (r *mutationResolver) UpdateCarrierLocation(ctx context.Context, input model.UpdateCarrierLocation) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	var carrier model.Carrier
	if err := r.carriers.Find(bson.M{"_id": userContext.ID}).One(&carrier); err != nil {
		return &model.Carrier{}, errors.New("No existe carrier con el TOKEN")
	}
	var fields = bson.M{}
	update := false
	if input.ActualLocation != nil {
		update = true
		fields["actual_location"] = input.ActualLocation
	}
	if !update {
		return &model.Carrier{}, errors.New("No hay campos por actualizar")
	}
	r.carriers.Update(bson.M{"_id": userContext.ID}, bson.M{"$set": fields})
	r.carriers.Find(bson.M{"_id": userContext.ID}).One(&carrier)
	r.Lock() //Enviando info subscripcion de ubicacion
	topic := r.carrierLocationTopics[carrier.ID]
	if topic != nil {
		for _, observer := range topic.Observers {
			observer <- &carrier
		}
	}
	r.Unlock()
	r.Lock() //Enviando info subscripcion de ubicacion a toda la tienda
	topicStore := r.storeCarriersLocationTopics[carrier.StoreID]
	if topicStore != nil {
		for _, observer := range topicStore.Observers {
			observer <- &carrier
		}
	}
	r.Unlock()
	return &carrier, nil
}

func (r *mutationResolver) DeleteCarrier(ctx context.Context, carrierID *string) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	var store model.Store
	var carrier model.Carrier
	carriersDB := db.GetCollection("carriers")
	storesDB := db.GetCollection("stores")
	if err := storesDB.Find(bson.M{"username": userContext.Username}).One(&store); err != nil {
		return &model.Carrier{}, errors.New("No existe este tienda")
	}
	if err := carriersDB.Find(bson.M{"_id": carrierID, "store_id": store.ID}).One(&carrier); err != nil {
		return &model.Carrier{}, errors.New("No existe este repartidor para esta tienda")
	}
	if err := carriersDB.Remove(bson.M{"_id": carrierID}); err != nil {
		return &model.Carrier{}, errors.New("error al borrar repartidor")
	}
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
	storeDB := db.GetCollection("stores")
	var store model.Store
	var ExitLocation model.Location
	id := bson.NewObjectId()
	loc := time.FixedZone("UTC-5", -5*60*60)
	t := time.Now().In(loc)
	userContext := auth.ForContext(ctx)
	if userContext == nil {
		return &model.Order{}, errors.New("Acceso denegado")
	}
	if err := storeDB.Find(bson.M{"username": userContext.Username}).One(&store); err != nil {
		return &model.Order{}, err
	}
	ExitLocation.Latitude = store.Location.Latitude
	ExitLocation.Longitude = store.Location.Longitude
	ExitLocation.Address = store.Location.Address
	ExitLocation.Reference = store.Location.Reference

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
	r.Lock()
	topic := r.storeOrdersTopics[order.StoreID]
	if topic != nil {
		for _, observer := range topic.Observers {
			observer <- &order
		}
	}
	r.Unlock()
	return &order, nil
}

func (r *mutationResolver) UpdateOrder(ctx context.Context, id string, input model.UpdateOrder) (*model.Order, error) {
	var order model.Order
	if err := r.orders.Find(bson.M{"_id": id}).One(&order); err != nil {
		return &model.Order{}, errors.New("No existe order")
	}

	var fields = bson.M{}
	loc := time.FixedZone("UTC-5", -5*60*60)
	t := time.Now().In(loc)
	update := false
	if input.CarrierID != nil {
		var carrier model.Carrier
		if err := r.carriers.Find(bson.M{"_id": input.CarrierID}).One(&carrier); err != nil {
			return &model.Order{}, errors.New("No existe carrier")
		}
		fields["carrier_id"] = *input.CarrierID
		update = true
		r.orders.Update(bson.M{"_id": id}, bson.M{"$set": fields})
		r.orders.Find(bson.M{"_id": id}).One(&order)
		//Actualizar estado de repartido a producto asignado
		var carrierFields = bson.M{}
		carrierFields["state_delivery"] = 2
		carrierFields["current_order_id"] = order.ID
		r.carriers.Update(bson.M{"_id": input.CarrierID}, bson.M{"$set": carrierFields})
		r.carriers.Find(bson.M{"_id": input.CarrierID}).One(&carrier)
		r.Lock() // Enviar la informacion del carrier actualizado a la tienda
		topic := r.storeCarriersTopics[order.StoreID]
		if topic != nil {
			for _, observer := range topic.Observers {
				observer <- &carrier
			}
		}
		r.Unlock()
		r.Lock() //Enviando info a carrier sobre asignacion
		topicOrder := r.carrierTopics[*input.CarrierID]
		if topicOrder != nil {
			for _, observer := range topicOrder.Observers {
				observer <- &carrier
			}
		}
		r.Unlock()
	}
	if input.State != nil {
		fields["state"] = *input.State
		update = true
		updateCarrier := false
		var carrierFields = bson.M{}
		if *input.State == 0 { //Pedido cancelado
			carrierFields["state_delivery"] = 1 // Actualizar el estado de repartidor a disponible a repartos
			fields["carrier_id"] = ""
			carrierFields["current_order_id"] = ""
			updateCarrier = true
		}
		if *input.State == 2 { // Pedido aceptado
			carrierFields["state_delivery"] = 3 //Cambiar estado de repartidor a llevando producto
			fields["departure_date"] = t.Format("2006-01-02T15:04:05")
			updateCarrier = true
		}
		if *input.State == 3 { // Pedido completado
			fields["delivery_date"] = t.Format("2006-01-02T15:04:05")
			carrierFields["current_order_id"] = ""
			carrierFields["state_delivery"] = 1 // Actualizar el estado de repartidor a disponible a repartos
			updateCarrier = true
		}
		if order.CarrierID != "" && updateCarrier { // Actualizar campos del repartidor
			//Actualizar estado de repartido a producto asignado
			var carrier model.Carrier
			r.carriers.Update(bson.M{"_id": order.CarrierID}, bson.M{"$set": carrierFields})
			r.carriers.Find(bson.M{"_id": order.CarrierID}).One(&carrier)
			r.Lock() // Enviar la informacion del carrier actualizado a la tienda
			topic := r.storeCarriersTopics[order.StoreID]
			if topic != nil {
				for _, observer := range topic.Observers {
					observer <- &carrier
				}
			}
			r.Unlock()
			r.Lock() //Enviando info a carrier sobre asignacion
			topicOrder := r.carrierTopics[order.CarrierID]
			if topicOrder != nil {
				for _, observer := range topicOrder.Observers {
					observer <- &carrier
				}
			}
			r.Unlock()
		}
		r.orders.Update(bson.M{"_id": id}, bson.M{"$set": fields})
		r.orders.Find(bson.M{"_id": id}).One(&order)

		r.Lock()
		topic := r.storeOrdersTopics[order.StoreID]
		if topic != nil {
			for _, observer := range topic.Observers {
				observer <- &order
			}
		}
		r.Unlock()
		r.Lock()
		topicOrder := r.orderTopics[order.ID]
		if topicOrder != nil {
			for _, observer := range topicOrder.Observers {
				observer <- &order
			}
		}
		r.Unlock()
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
	r.orders.Update(bson.M{"_id": id}, bson.M{"$set": fields})
	r.orders.Find(bson.M{"_id": id}).One(&order)

	return &order, nil
}

func (r *orderResolver) ArrivalLocation(ctx context.Context, obj *model.Order) (*model.Location, error) {
	return &obj.ArrivalLocation, nil
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

func (r *queryResolver) Carrier(ctx context.Context, id *string) (*model.Carrier, error) {
	userContext := auth.ForContext(ctx)
	carriersDB := db.GetCollection("carriers")
	var carrier model.Carrier
	if userContext == nil {
		return &model.Carrier{}, errors.New("Acceso denegado")
	}
	if id != nil {
		if err := carriersDB.Find(bson.M{"_id": id}).Select(bson.M{"token": 0}).One(&carrier); err != nil {
			return &model.Carrier{}, errors.New("No existe carrier con ese id")
		}
	} else if err := carriersDB.Find(bson.M{"username": userContext.Username}).One(&carrier); err != nil {
		return &model.Carrier{}, errors.New("No existe carrier con ese token")
	}
	return &carrier, nil
}

func (r *queryResolver) Carriers(ctx context.Context, limit *int, search *string, global *bool, stateDelivery *int) ([]*model.Carrier, error) {
	userContext := auth.ForContext(ctx)

	var carriers []*model.Carrier
	var fields = bson.M{}
	if search != nil {
		fields["name"] = bson.M{"$regex": *search, "$options": "i"}
	}
	if global != nil {
		fields["global"] = *global
	}
	if stateDelivery != nil {
		fields["state_delivery"] = stateDelivery
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
	} else if userContext != nil && userContext.UserType == "carrier" {
		fields["carrier_id"] = userContext.ID
	} else {
		return []*model.Order{}, errors.New("Acceso denegado")
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
	if input.Ids != nil {
		fields["_id"] = bson.M{"$in": input.Ids}
	}
	if len(orArray) > 0 {
		fields["$or"] = orArray
	}

	ordersDB.Find(fields).Limit(input.Limit).Sort("-updated_at").All(&orders)

	return orders, nil
}

func (r *subscriptionResolver) StoreCarriers(ctx context.Context, storeID string) (<-chan *model.Carrier, error) {
	r.Lock()
	topic := r.storeCarriersTopics[storeID]
	if topic == nil {
		topic = &StoreCarriersTopic{
			Key:       storeID,
			Observers: map[string]chan *model.Carrier{},
		}
		r.storeCarriersTopics[storeID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Carrier, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

func (r *subscriptionResolver) StoreOrders(ctx context.Context, storeID string) (<-chan *model.Order, error) {
	r.Lock()
	topic := r.storeOrdersTopics[storeID]
	if topic == nil {
		topic = &StoreOrdersTopic{
			Key:       storeID,
			Observers: map[string]chan *model.Order{},
		}
		r.storeOrdersTopics[storeID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Order, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

func (r *subscriptionResolver) Order(ctx context.Context, orderID string) (<-chan *model.Order, error) {
	r.Lock()
	topic := r.orderTopics[orderID]
	if topic == nil {
		topic = &OrderTopic{
			Key:       orderID,
			Observers: map[string]chan *model.Order{},
		}
		r.orderTopics[orderID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Order, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

func (r *subscriptionResolver) Carrier(ctx context.Context, carrierID string) (<-chan *model.Carrier, error) {
	r.Lock()
	topic := r.carrierTopics[carrierID]
	if topic == nil {
		topic = &CarrierTopic{
			Key:       carrierID,
			Observers: map[string]chan *model.Carrier{},
		}
		r.carrierTopics[carrierID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Carrier, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

func (r *subscriptionResolver) CarrierLocation(ctx context.Context, carrierID string) (<-chan *model.Carrier, error) {
	r.Lock()
	topic := r.carrierLocationTopics[carrierID]
	if topic == nil {
		topic = &CarrierLocationTopic{
			Key:       carrierID,
			Observers: map[string]chan *model.Carrier{},
		}
		r.carrierLocationTopics[carrierID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Carrier, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

func (r *subscriptionResolver) StoreCarrierLocation(ctx context.Context, storeID string) (<-chan *model.Carrier, error) {
	r.Lock()
	topic := r.storeCarriersLocationTopics[storeID]
	if topic == nil {
		topic = &StoreCarriersLocationTopic{
			Key:       storeID,
			Observers: map[string]chan *model.Carrier{},
		}
		r.storeCarriersLocationTopics[storeID] = topic
	}
	r.Unlock()
	id := RandStringRunes(8)
	event := make(chan *model.Carrier, 1)
	go func() {
		<-ctx.Done()
		r.Lock()
		delete(topic.Observers, id)
		r.Unlock()
	}()

	r.Lock()
	topic.Observers[id] = event
	r.Unlock()
	return event, nil
}

// Carrier returns generated.CarrierResolver implementation.
func (r *Resolver) Carrier() generated.CarrierResolver { return &carrierResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Order returns generated.OrderResolver implementation.
func (r *Resolver) Order() generated.OrderResolver { return &orderResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type carrierResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type orderResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
