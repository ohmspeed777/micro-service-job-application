package domain

type User struct {
	Model     `bson:",inline"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"-" bson:"password"`
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
}
