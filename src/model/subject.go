package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colNameSubject = "subject"

type Subject struct {
	ID         primitive.ObjectID `bson:"_id"`
	Abbr       string             `bson:"abbr"`
	Name       string             `bson:"name"`
	StartScene primitive.ObjectID `bson:"start_scene"`
}

func (m *model) GetSubjectList() ([]Subject, error) {
	c := m.db.Collection(colNameSubject)
	result, err := c.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var subjectList []Subject
	err = result.All(m.ctx, &subjectList)
	if err != nil {
		return nil, err
	}
	return subjectList, nil
}
