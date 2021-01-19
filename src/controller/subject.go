package controller

import (
	"github.com/KSkun/tqb-backend/controller/param"
	"github.com/KSkun/tqb-backend/model"
	"github.com/KSkun/tqb-backend/util/context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SubjectGetList(ctx echo.Context) error {
	m := model.GetModel()
	defer m.Close()

	subjectList, err := m.GetSubjectList()
	if err != nil {
		return context.Error(ctx, http.StatusInternalServerError, "failed to get subject list", err)
	}

	ret := make([]param.ObjRspSubject, 0)
	for _, subject := range subjectList {
		ret = append(ret, param.ObjRspSubject{
			Abbr:       subject.Abbr,
			Name:       subject.Name,
			StartScene: subject.StartScene.Hex(),
		})
	}
	return context.Success(ctx, param.RspSubjectGetList{Subject: ret})
}
