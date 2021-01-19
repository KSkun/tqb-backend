package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colNameScene = "scene"

type Scene struct {
	ID           primitive.ObjectID `bson:"_id"`
	FromQuestion primitive.ObjectID `bson:"from_question"`
	NextQuestion primitive.ObjectID `bson:"next_question"`
	Title        string             `bson:"title"`
	Text         string             `bson:"text"`
}

func (m *model) GetSceneList() ([]Scene, error) {
	c := m.db.Collection(colNameScene)
	result, err := c.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var sceneList []Scene
	err = result.All(m.ctx, &sceneList)
	if err != nil {
		return nil, err
	}
	return sceneList, nil
}

func (m *model) GetScene(id primitive.ObjectID) (Scene, error) {
	c := m.db.Collection(colNameScene)
	scene := Scene{}
	err := c.FindOne(m.ctx, bson.M{"_id": id}).Decode(&scene)
	return scene, err
}
