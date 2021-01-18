package model

import (
	"context"
	"github.com/KSkun/tqb-backend/config"
	"github.com/go-pg/pg/v10"
	"time"
)

var ErrNotFound error = pg.ErrNoRows

type ChatModel interface {
	AddMessage(msg, association, newerID string) error
	GetMessage(association, newerID string) ([]string, error)
	ClearMessage(association, newerID string) error
	IncrChattingCounter(association, newerID string) error
	DecrChattingCounter(association, newerID string) error
	GetChatCount(association, newerID string) (int, error)
	Close()
}

type Model interface {
	// 关闭数据库连接
	Close()
	// 终止操作，用于如事务的取消
	Abort()
	// TODO: 将Model层的实现列在这里，然后再去实现model结构体中的对应实现
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
