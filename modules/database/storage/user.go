package storage

import (
	"context"
	"errors"
	"github.com/firat.atmaca/go-mongo/model"
	"github.com/firat.atmaca/go-mongo/modules/config"
	mongoClient "github.com/firat.atmaca/go-mongo/modules/database/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"regexp"
	"time"
)

type UserStorage struct {
	App    *config.AppTools
	Client mongoClient.IClient
}

type IUserStorage interface {
	InsertUser(user *model.User) (bool, int, error)
	VerifyUser(email string) (primitive.M, error)
	UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error)
}

func NewUserStorage(app *config.AppTools, client mongoClient.IClient) *UserStorage {
	return &UserStorage{
		App:    app,
		Client: client,
	}
}

func (g *UserStorage) InsertUser(user *model.User) (bool, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	regMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	ok := regMail.MatchString(user.Email)
	if !ok {
		g.App.ErrorLogger.Println("invalid registered details")
		return false, 0, errors.New("invalid registered details")
	}

	filter := bson.D{{Key: "email", Value: user.Email}}

	var res bson.M
	err := g.Client.Collection("user").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			user.ID = primitive.NewObjectID()
			_, insertErr := g.Client.Collection("user").InsertOne(ctx, user)
			if insertErr != nil {
				g.App.ErrorLogger.Fatalf("cannot add user to the database : %v ", insertErr)
			}
			return true, 1, nil
		}
		g.App.ErrorLogger.Fatal(err)
	}
	return true, 2, nil
}

func (g *UserStorage) VerifyUser(email string) (primitive.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var res bson.M

	filter := bson.D{{Key: "email", Value: email}}
	err := g.Client.Collection("user").FindOne(ctx, filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			g.App.ErrorLogger.Println("no document found for this query")
			return nil, err
		}
		g.App.ErrorLogger.Fatalf("cannot execute the database query perfectly : %v ", err)
	}

	return res, nil
}

func (g *UserStorage) UpdateInfo(userID primitive.ObjectID, tk map[string]string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: tk["t1"]}, {Key: "new_token", Value: tk["t2"]}}}}

	_, err := g.Client.Collection("user").UpdateOne(ctx, filter, update)
	if err != nil {
		return false, err
	}
	return true, nil
}
