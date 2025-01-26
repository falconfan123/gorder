package adapters

import (
	"context"
	domain "github.com/falconfan123/gorder/order/domain/order"
	"github.com/falconfan123/gorder/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type OrderRepositoryMongo struct {
	db *mongo.Client
}

var (
	DBname   = viper.GetString("mongo.db-name")
	Collname = viper.GetString("mongo.coll-name")
)

func (r *OrderRepositoryMongo) collection() *mongo.Collection {
	return r.db.Database(DBname).Collection(Collname)
}

type writeModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `json:"customer_id"`
	Status      string             `json:"status"`
	PaymentLink string             `json:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer func() {
		l := logrus.WithFields(logrus.Fields{
			"create_time":  time.Now().Unix(),
			"err":          err,
			"create_order": created,
		})
		if err != nil {
			l.Info("order_repository_mongo_Create_Fail")
		} else {
			l.Info("order_repository_mongo_Create_Success")
		}
	}()
	write := r.marshellToModel(order)
	res, err := r.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (*domain.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (r *OrderRepositoryMongo) Update(ctx context.Context, o *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) error {
	//TODO implement me
	panic("implement me")
}

func (r *OrderRepositoryMongo) marshellToModel(order *domain.Order) interface{} {
	return &writeModel{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}
