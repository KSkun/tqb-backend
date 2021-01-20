package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colNameQuestion = "question"

type SubQuestion struct {
	Type       int      `bson:"type"`
	Desc       string   `bson:"desc"`
	Option     []string `bson:"option"`
	TrueOption []int    `bson:"true_option"`
	FullPoint  float32  `bson:"full_point"`
	PartPoint  float32  `bson:"part_point"`
}

type NextScene struct {
	Scene  primitive.ObjectID `bson:"scene"`
	Option string             `bson:"option"`
}

type Question struct {
	ID          primitive.ObjectID `bson:"_id"`
	Title       string             `bson:"title"`
	Desc        string             `bson:"desc"`
	SubQuestion []SubQuestion      `bson:"sub_question"`
	Author      string             `bson:"author"`
	Audio       string             `bson:"audio"`
	TimeLimit   int                `bson:"time_limit"`
	NextScene   []NextScene        `bson:"next_scene"`
}

func (m *model) GetQuestionList() ([]Question, error) {
	c := m.db.Collection(colNameQuestion)
	result, err := c.Find(m.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var questionList []Question
	err = result.All(m.ctx, &questionList)
	if err != nil {
		return nil, err
	}
	return questionList, nil
}

func (m *model) GetQuestion(id primitive.ObjectID) (Question, error) {
	c := m.db.Collection(colNameQuestion)
	question := Question{}
	err := c.FindOne(m.ctx, bson.M{"_id": id}).Decode(&question)
	return question, err
}
