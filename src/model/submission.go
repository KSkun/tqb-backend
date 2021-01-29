package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const colNameSubmission = "submission"

const PointUnknown = -1.0 // 未批改

type Submission struct {
	ID         primitive.ObjectID   `bson:"_id"`
	User       primitive.ObjectID   `bson:"user"`
	Question   primitive.ObjectID   `bson:"question"`
	Time       int64                `bson:"time"`
	File       []primitive.ObjectID `bson:"file"`
	Option     [][]int              `bson:"option"`
	Point      float64              `bson:"point"`
	AnswerTime int                  `bson:"answer_time"`
	IsTimeOut  bool                 `bson:"is_time_out"`
}

func (m *model) GetSubmissionByUser(userID primitive.ObjectID, filter bson.M) ([]Submission, error) {
	c := m.db.Collection(colNameSubmission)
	filter["user"] = userID
	result, err := c.Find(m.ctx, filter)
	if err != nil {
		return nil, err
	}
	var submissionList []Submission
	err = result.All(m.ctx, &submissionList)
	if err != nil {
		return nil, err
	}
	return submissionList, nil
}

func (m *model) AddSubmission(submission Submission) (primitive.ObjectID, error) {
	c := m.db.Collection(colNameSubmission)
	id := primitive.NewObjectID()
	submission.ID = id
	_, err := c.InsertOne(m.ctx, submission)
	return id, err
}

func (m *model) DeleteSubmissionByUser(userID primitive.ObjectID) error {
	c := m.db.Collection(colNameSubmission)
	_, err := c.DeleteMany(m.ctx, bson.M{"user": userID})
	return err
}
