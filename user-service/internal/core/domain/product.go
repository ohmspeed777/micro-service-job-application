package domain

type Product struct {
	Model `bson:",inline"`
	Price float64 `json:"price" bson:"price"`
	Stock uint    `json:"stock" bson:"stock"`
	Name  string  `json:"name" bson:"name"`
}
