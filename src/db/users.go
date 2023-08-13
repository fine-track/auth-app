package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Email     string             `bson:"email" json:"email"`
	Fullname  string             `bson:"fullname" json:"fullname"`
	Password  string             `bson:"password" json:"password"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}

func (u *User) GetById(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	fmt.Printf("getting user with: %v\n", id)
	if err != nil {
		return err
	}
	err = UsersCol.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByEmail(email string) error {
	err := UsersCol.FindOne(context.TODO(), bson.M{"email": email}).Decode(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) CreateNew() error {
	payload := bson.M{
		"email":      u.Email,
		"fullname":   u.Fullname,
		"password":   u.Password,
		"created_at": primitive.NewDateTimeFromTime(time.Now()),
		"updated_at": primitive.NewDateTimeFromTime(time.Now()),
	}
	user, err := UsersCol.InsertOne(context.TODO(), payload)
	if err != nil {
		return err
	}
	u.ID = user.InsertedID.(primitive.ObjectID)
	fmt.Print(user)
	return nil
}

func (u *User) UpdateByEmail() error {
	p := bson.M{
		"email":    u.Email,
		"fullname": u.Fullname,
		"password": u.Password,
	}
	_, err := UsersCol.UpdateOne(context.TODO(), bson.M{"email": u.Email}, p)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateByID() error {
	p := bson.M{"$set": bson.M{
		"email":    u.Email,
		"fullname": u.Fullname,
		"password": u.Password,
	}}
	_, err := UsersCol.UpdateByID(context.TODO(), u.ID, p)
	if err != nil {
		return err
	}
	return nil
}
