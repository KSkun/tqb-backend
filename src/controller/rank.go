package controller

import (
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/KSkun/tqb-backend/util/log"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sort"
	"time"
)

func RankGetList(ctx echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	list, err := m.GetRankList()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get rank list", err)
	}
	return context.Success(ctx, echo.Map{"rank": list})
}

func getUserPoint(m model.Model, userID primitive.ObjectID) (float64, error) {
	submissionList, err := m.GetSubmissionByUser(userID, bson.M{})
	if err != nil {
		return 0, err
	}

	sum := 0.0
	for _, submission := range submissionList {
		if submission.Point == model.PointUnknown {
			continue
		}
		sum += submission.Point
	}
	return sum, nil
}

func RankListWorker() {
	m := model.GetModel()
	defer m.Close()
	for {
		userList, err := m.GetUserList()
		if err != nil {
			log.Logger.Error(err)
			return
		}

		rankList := make([]model.RankEntry, 0)
		for _, user := range userList {
			point, err := getUserPoint(m, user.ID)
			if err != nil {
				log.Logger.Error(err)
				return
			}

			rankList = append(rankList, model.RankEntry{
				Username: user.Username,
				Email:    user.Email,
				Point:    point,
			})
		}
		rankListSorted := model.RankList(rankList)
		sort.Sort(rankListSorted)

		err = m.SaveRankList(rankListSorted)
		if err != nil {
			log.Logger.Error(err)
			return
		}

		time.Sleep(time.Minute * 5) // 每 5 min 刷新一次排行榜
	}
}
