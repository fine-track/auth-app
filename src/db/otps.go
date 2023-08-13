package db

import (
	"context"
	"time"

	"github.com/fine-track/auth-app/src/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const OTP_LEN = 6

type OTPSession struct {
	ID        primitive.ObjectID `bson:"_id" json:"_id"`
	Code      string             `bson:"code" json:"code"`
	Email     string             `bson:"email" json:"email"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

func (otp *OTPSession) CreateNew() error {
	doc, err := OTPSessionsCol.InsertOne(context.TODO(), bson.M{
		"code":       utils.GenRandomStr(OTP_LEN),
		"email":      otp.Email,
		"created_at": primitive.NewDateTimeFromTime(time.Now()),
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
