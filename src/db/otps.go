package db

import (
	"context"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OTP_LEN = 6

func genRandomStr(length int) string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // using all uppercase for better user experience
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type OTPSession struct {
	ID        primitive.ObjectID `bson:"_id"`
	Code      string             `bson:"code"`
	Email     string             `bson:"email"`
	CreatedAt primitive.DateTime `bson:"createdAt"`
}

func (otp *OTPSession) CreateNew() error {
	doc, err := OTPSessionsCol.InsertOne(context.TODO(), bson.M{
		"code":      genRandomStr(OTP_LEN),
		"email":     otp.Email,
		"createdAt": otp.CreatedAt,
	})
	if err != nil {
		return err
	}
	otp.ID = doc.InsertedID.(primitive.ObjectID)
	return nil
}

func (otp *OTPSession) GetByEmail(email string) error {
	err := OTPSessionsCol.FindOne(context.TODO(), bson.M{"email": email}).Decode(otp)
	if err != nil {
		return err
	}
	return nil
}
