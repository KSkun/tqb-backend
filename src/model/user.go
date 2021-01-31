package model

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/KSkun/tqb-backend/util/log"
	"github.com/go-redis/redis/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const colNameUser = "user"

const (
	timePrivateKey  = time.Minute * 15 // 15 min
	keyPrivateKey   = "priv_key:%s"
	timeVerifyID    = time.Minute * 15 // 15 min
	keyVerifyID     = "verify_id:%s"
	timeVerifyEmail = time.Minute // 1 min
	keyVerifyEmail  = "verify_email:%s"
	keyUserInfo     = "user_info:%s"
	timeUserInfo    = time.Minute * 15 // 15 min
)

type User struct {
	ID               primitive.ObjectID   `bson:"_id" json:"-"`
	Username         string               `bson:"username" json:"username"`
	Password         string               `bson:"password" json:"password"`
	Email            string               `bson:"email" json:"email"`
	LastQuestion     primitive.ObjectID   `bson:"last_question" json:"-"`
	LLastQuestion    primitive.ObjectID   `bson:"l_last_question" json:"-"`
	LastScene        primitive.ObjectID   `bson:"last_scene" json:"-"`
	StartTime        int64                `bson:"start_time" json:"-"`
	UnlockedScene    []primitive.ObjectID `bson:"unlocked_scene" json:"-"`
	FinishedQuestion []primitive.ObjectID `bson:"finished_question" json:"-"`
	CompleteCount    int                  `bson:"complete_count" json:"-"`
}

func (m *model) GetUser(id primitive.ObjectID) (User, error) {
	c := m.db.Collection(colNameUser)
	user := User{}
	err := c.FindOne(m.ctx, bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (m *model) GetUserByEmail(email string) (User, bool, error) {
	c := m.db.Collection(colNameUser)
	user := User{}
	err := c.FindOne(m.ctx, bson.M{"email": email}).Decode(&user)
	if err == ErrNotFound {
		return user, false, nil
	}
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

func (m *model) GetUserList() ([]User, error) {
	c := m.db.Collection(colNameUser)
	userList := make([]User, 0)
	result, err := c.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = result.All(m.ctx, &userList)
	if err != nil {
		return nil, err
	}
	return userList, nil
}

func (m *model) AddUser(user User) (primitive.ObjectID, error) {
	c := m.db.Collection(colNameUser)
	userID := primitive.NewObjectID()
	user.ID = userID
	_, err := c.InsertOne(m.ctx, user)
	return userID, err
}

func (m *model) UpdateUser(id primitive.ObjectID, toUpdate bson.M) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, bson.M{"$set": toUpdate})
	return err
}

func (m *model) AddPrivateKey(email string, key *rsa.PrivateKey) error {
	return redisClient.Set(fmt.Sprintf(keyPrivateKey, email),
		x509.MarshalPKCS1PrivateKey(key), timePrivateKey).Err()
}

func (m *model) GetPrivateKey(email string) (*rsa.PrivateKey, bool, error) {
	result := redisClient.Get(fmt.Sprintf(keyPrivateKey, email))
	if result.Err() == redis.Nil {
		return &rsa.PrivateKey{}, false, nil
	}
	if result.Err() != nil {
		return &rsa.PrivateKey{}, false, result.Err()
	}
	key, err := x509.ParsePKCS1PrivateKey([]byte(result.Val()))
	if err != nil {
		return &rsa.PrivateKey{}, false, err
	}
	return key, true, nil
}

func (m *model) AddVerifyID(email string, id string) error {
	err := redisClient.Set(fmt.Sprintf(keyVerifyID, id), email, timeVerifyID).Err()
	if err != nil {
		return err
	}
	err = redisClient.Set(fmt.Sprintf(keyVerifyEmail, email), id, timeVerifyEmail).Err()
	if err != nil {
		return err
	}
	return nil
}

func (m *model) GetVerifyID(id string) (string, bool, error) {
	result := redisClient.Get(fmt.Sprintf(keyVerifyID, id))
	if result.Err() == redis.Nil {
		return "", false, nil
	}
	if result.Err() != nil {
		return "", false, result.Err()
	}
	resultDel := redisClient.Del(fmt.Sprintf(keyVerifyID, id))
	if resultDel.Err() != nil && resultDel.Err() != redis.Nil {
		log.Logger.Error(resultDel.Err())
	}
	return result.Val(), true, nil
}

func (m *model) GetVerifyIDByEmail(email string) (string, bool, error) {
	result := redisClient.Get(fmt.Sprintf(keyVerifyEmail, email))
	if result.Err() == redis.Nil {
		return "", false, nil
	}
	if result.Err() != nil {
		return "", false, result.Err()
	}
	return result.Val(), true, nil
}

func (m *model) SetUserLastScene(id primitive.ObjectID, sceneID primitive.ObjectID) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, bson.M{
		"$set":      bson.M{"last_scene": sceneID},
		"$addToSet": bson.M{"unlocked_scene": sceneID},
	})
	return err
}

func (m *model) SetUserLastQuestion(id primitive.ObjectID, questionID primitive.ObjectID) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, []bson.M{
		{"$set": bson.M{
			"l_last_question": "$last_question",
			"last_question":   questionID,
			"start_time":      time.Now().Unix(),
		}},
	})
	return err
}

func (m *model) SetUserBackQuestion(id primitive.ObjectID) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, []bson.M{
		{"$set": bson.M{
			"last_question":   "$l_last_question",
			"l_last_question": primitive.ObjectID{},
			"start_time":      time.Now().Unix(),
		}},
	})
	return err
}

func (m *model) AddUserFinishedQuestion(id primitive.ObjectID, questionID primitive.ObjectID) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, bson.M{
		"$addToSet": bson.M{"finished_question": questionID},
	})
	return err
}

func (m *model) UserHasUnlockedScene(id primitive.ObjectID, sceneID primitive.ObjectID) (bool, error) {
	c := m.db.Collection(colNameUser)
	err := c.FindOne(m.ctx, bson.M{
		"_id":            id,
		"unlocked_scene": bson.M{"$elemMatch": bson.M{"$eq": sceneID}},
	}).Err()
	if err == ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *model) UserHasFinishedQuestion(id primitive.ObjectID, questionID primitive.ObjectID) (bool, error) {
	c := m.db.Collection(colNameUser)
	err := c.FindOne(m.ctx, bson.M{
		"_id":               id,
		"finished_question": bson.M{"$elemMatch": bson.M{"$eq": questionID}},
	}).Err()
	if err == ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *model) IncUserCompleteCount(id primitive.ObjectID) error {
	c := m.db.Collection(colNameUser)
	_, err := c.UpdateOne(m.ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"complete_count": 1}})
	return err
}

func (m *model) UserIsAllUnlocked(id primitive.ObjectID) (bool, error) {
	c := m.db.Collection(colNameUser)
	result := c.FindOne(m.ctx, bson.M{"_id": id, "complete_count": bson.M{"$gte": 2}})
	if result.Err() == ErrNotFound {
		return false, nil
	}
	if result.Err() != nil {
		return false, result.Err()
	}
	return true, nil
}

func (m *model) AddTempUserInfo(user User) error {
	userStr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	result := redisClient.Set(fmt.Sprintf(keyUserInfo, user.Email), string(userStr), timeUserInfo)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (m *model) GetTempUserInfo(email string) (User, bool, error) {
	result := redisClient.Get(fmt.Sprintf(keyUserInfo, email))
	if result.Err() == redis.Nil {
		return User{}, false, nil
	}
	if result.Err() != nil {
		return User{}, false, result.Err()
	}
	user := User{}
	err := json.Unmarshal([]byte(result.Val()), &user)
	if err != nil {
		return User{}, false, err
	}
	user.LLastQuestion = primitive.ObjectID{}
	user.LastQuestion = primitive.ObjectID{}
	user.LastScene = primitive.ObjectID{}
	user.CompleteCount = 0
	user.FinishedQuestion = make([]primitive.ObjectID, 0)
	user.UnlockedScene = make([]primitive.ObjectID, 0)
	return user, true, nil
}
