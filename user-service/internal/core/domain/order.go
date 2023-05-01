package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderStatus int

const (
	Created = iota + 1
	Canceled
)

type Order struct {
	Model  `bson:",inline"`
	User   primitive.ObjectID `json:"user" bson:"user"`
	Status OrderStatus        `json:"status" bson:"status"`
	Items  []*OrderItem       `json:"items" bson:"items"`
	Total  float64            `json:"total" bson:"total"`
}

type OrderItem struct {
	Quantity  uint               `json:"quantity" bson:"quantity"`
	ProductID primitive.ObjectID `json:"product_id" bson:"product_id"`
	Product   *Product           `json:"product,omitempty" bson:"-"`
}

type OrderLookedUp struct {
	Order          `bson:",inline"`
	ProductsJoined []*Product `json:"-" bson:"products_joined"`
	UserDetail     *User      `json:"user_detail" bson:"-"`
}

func (o *OrderLookedUp) Format(user *User) *OrderLookedUp {
	for i, item := range o.Items {
		for _, product := range o.ProductsJoined {
			if item.ProductID == product.ID {
				o.Items[i].Product = product
			}
		}
	}

	o.UserDetail = user
	return o
}
