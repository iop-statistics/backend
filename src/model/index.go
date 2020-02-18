package model

import (
	"context"
	"github.com/Lyt99/iop-statistics/config"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

var DBName = config.GlobalConfig.Mongo.DB

type Model interface {
	GetValidFormulasByID(id, t int) ([]Formula, error) // 这里就要做处理了
	GetIDsByFormula(f Formula, t int) ([]int, error)
	GetFormulaCount(f Formula, t int, from, to int) (int, error)
	GetFormulaAndIDCount(f Formula, id int, t int, from, to int) (int, error)
	GetTimeRecordByIDAndType(id, t int) (TimeRecord, error)
	GetInfoByKey(key string) (interface{}, error)
	GetColCount(col string) (int64, error)
	Close()
}

type model struct {
	ctx context.Context
	db  *mongo.Database
}

func GetModel() Model {
	var ctx context.Context
	if config.GlobalConfig.EnableDebug {
		ctx = context.Background()
	} else {
		ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	}

	_ = MongoClient.Connect(ctx)
	return &model{
		ctx: ctx,
		db:  MongoClient.Database(DBName),
	}
}

func (m *model) Close() {
	m.ctx.Done()
}
