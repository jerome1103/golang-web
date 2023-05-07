package users

type User struct {
	Uid   int    `bson:"uid"`
	Name  string `bson:"name"`
	Phone string `bson:"phone"`
	Text  string `bson:"text"`
}
