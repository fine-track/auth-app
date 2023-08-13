package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	UserId    primitive.ObjectID `bson:"user_id" json:"user_id"`
	Email     string             `bson:"email" json:"email"`
	IP        string             `bson:"ip" json:"ip"`
	UserAgent string             `bson:"user_agent" json:"user_agent"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

func (s *Session) GetById(id string) error {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	err = SessionsCol.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(s)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) CreateNew() error {
	session, err := SessionsCol.InsertOne(context.TODO(), bson.M{
		"email":      s.Email,
		"user_id":    s.UserId,
		"ip":         s.IP,
		"user_agent": s.UserAgent,
		"created_at": s.CreatedAt,
	})
	if err != nil {
		return err
	}
	s.ID = session.InsertedID.(primitive.ObjectID)
	return nil
}

func (s *Session) Remove() error {
	_, err := SessionsCol.DeleteOne(context.TODO(), bson.M{"_id": s.ID})
	return err
}

func RemoveUserSessions(email string) error {
	_, err := SessionsCol.DeleteMany(context.TODO(), bson.M{"email": email})
	return err
}
