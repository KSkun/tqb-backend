package model

import (
	"context"
	"crypto/rsa"
	"github.com/KSkun/tqb-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var ErrNotFound error = mongo.ErrNoDocuments

type Model interface {
	// 关闭数据库连接
	Close()
	// 终止操作，用于如事务的取消
	Abort()
	// user
	GetUser(id primitive.ObjectID) (User, error)
	GetUserByEmail(email string) (User, bool, error)
	AddUser(user User) (primitive.ObjectID, error)
	UpdateUser(id primitive.ObjectID, toUpdate bson.M) error
	AddPrivateKey(email string, key *rsa.PrivateKey) error
	GetPrivateKey(email string) (*rsa.PrivateKey, bool, error)
	AddVerifyID(email string, id string) error
	GetVerifyID(id string) (string, bool, error)
	GetVerifyIDByEmail(email string) (string, bool, error)
	SetUserLastScene(id primitive.ObjectID, sceneID primitive.ObjectID) error
	UserHasUnlockedScene(id primitive.ObjectID, sceneID primitive.ObjectID) (bool, error)
	// subject
	GetSubjectList() ([]Subject, error)
	// scene
	GetSceneList() ([]Scene, error)
	GetScene(id primitive.ObjectID) (Scene, error)
}

type model struct {
	dbTrait
	ctx   context.Context
	abort bool
}

func GetModel() Model {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if config.C.Debug {
		ctx = context.Background()
	}

	ret := &model{
		dbTrait: getDBTx(ctx),
		ctx:     ctx,
		abort:   false,
	}

	return ret
}
