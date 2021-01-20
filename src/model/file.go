package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colNameFile = "file"

type File struct {
	ID       primitive.ObjectID `bson:"_id"`
	Filename string             `bson:"filename"`
	User     primitive.ObjectID `bson:"user"`
	Time     int64              `bson:"time"`
}

func (m *model) AddFile(file File) (primitive.ObjectID, error) {
	c := m.db.Collection(colNameFile)
	id := primitive.NewObjectID()
	file.ID = id
	_, err := c.InsertOne(m.ctx, file)
	return id, err
}

func (m *model) GetFile(id primitive.ObjectID) (File, error) {
	c := m.db.Collection(colNameFile)
	file := File{}
	err := c.FindOne(m.ctx, bson.M{"_id": id}).Decode(&file)
	return file, err
}
