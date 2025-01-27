package adapters

import (
	"context"
	domain "github.com/falconfan123/gorder/order/domain/order"
	"github.com/falconfan123/gorder/order/entity"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
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

func NewOrderRepositoryMongo(db *mongo.Client) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{db: db}
}

type orderModel struct {
	MongoID     primitive.ObjectID `bson:"_id"`
	ID          string             `bson:"id"`
	CustomerID  string             `json:"customer_id"`
	Status      string             `json:"status"`
	PaymentLink string             `json:"payment_link"`
	Items       []*entity.Item     `bson:"items"`
}

func (r *OrderRepositoryMongo) Create(ctx context.Context, order *domain.Order) (created *domain.Order, err error) {
	defer r.logWithTag("create", err, order, created)
	write := r.marshellToModel(order)
	res, err := r.collection().InsertOne(ctx, write)
	if err != nil {
		return nil, err
	}
	created = order
	created.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (r *OrderRepositoryMongo) logWithTag(tag string, err error, input *domain.Order, result interface{}) {
	func() {
		l := logrus.WithFields(logrus.Fields{
			"tag":         "order_repository_mongo",
			"input_order": input,
			"create_time": time.Now().Unix(),
			"err":         err,
			"result":      result,
		})
		if err != nil {
			l.Infof("%s_fail", tag)
		} else {
			l.Infof("%s_success", tag)
		}
	}()
}

func (r *OrderRepositoryMongo) Get(ctx context.Context, id, customerID string) (got *domain.Order, err error) {
	defer r.logWithTag("get", err, nil, got)
	read := &orderModel{}
	mongoID, _ := primitive.ObjectIDFromHex(id)
	cond := bson.M{"_id": mongoID} //查询语句
	if err = r.collection().FindOne(ctx, cond).Decode(read); err != nil {
		return
	}
	if read == nil {
		return nil, domain.NotFoundError{
			OrderID: id,
		}
	}
	got = r.unmarshal(read)
	return got, nil
}

// Update 先查找对应的order,然后apply updateFn,再写入回去
func (r *OrderRepositoryMongo) Update(ctx context.Context, o *domain.Order, updateFn func(context.Context, *domain.Order) (*domain.Order, error)) (err error) {
	defer r.logWithTag("after_update", err, o, nil)
	if o == nil {
		panic("got nil order")
	}
	//事务
	session, err := r.db.StartSession()
	if err != nil {
		return
	}
	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}
	defer func() {
		if err == nil {
			_ = session.CommitTransaction(ctx)
		} else {
			_ = session.AbortTransaction(ctx)
		}
	}()

	//inside transaction
	oldOrder, err := r.Get(ctx, o.ID, o.CustomerID)
	if err != nil {
		return
	}
	updated, err := updateFn(ctx, o)
	if err != nil {
		return
	}
	mongoID, _ := primitive.ObjectIDFromHex(oldOrder.ID)
	res, err := r.collection().UpdateOne(
		ctx,
		bson.M{"_id": mongoID, "customer_id": oldOrder.CustomerID},
		bson.M{"$set": bson.M{
			"status":       updated.Status,
			"payment_link": updated.PaymentLink,
		}},
	)
	if err != nil {
		return
	}
	r.logWithTag("finish_update", err, o, res)
	return
}

func (r *OrderRepositoryMongo) marshellToModel(order *domain.Order) interface{} {
	return &orderModel{
		MongoID:     primitive.NewObjectID(),
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      order.Status,
		PaymentLink: order.PaymentLink,
		Items:       order.Items,
	}
}

func (r *OrderRepositoryMongo) unmarshal(m *orderModel) *domain.Order {
	return &domain.Order{
		ID:          m.MongoID.Hex(),
		CustomerID:  m.CustomerID,
		Status:      m.Status,
		PaymentLink: m.PaymentLink,
		Items:       m.Items,
	}
}
