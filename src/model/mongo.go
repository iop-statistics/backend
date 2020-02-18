package model

import (
	"context"
	"fmt"
	"github.com/Lyt99/iop-statistics/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var (
	MongoClient *mongo.Client
)

func init() {
	var err error
	uri := fmt.Sprintf("mongodb://%s:%s", config.GlobalConfig.Mongo.Host, config.GlobalConfig.Mongo.Port)
	MongoClient, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	// Ping test
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = MongoClient.Connect(ctx)
	if err != nil {
		panic(err)
	}

	err = MongoClient.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
}

const (
	TypeTdoll = 0
	TypeEquip = 1
	TypeFairy = 2

	ColStats        = "stats"
	ColTdollRecords = "tdoll_records"
	ColEquipRecords = "equip_records"
	ColFormula      = "formula"
	ColFormulaIndex = "formula_index"
)

// Formula 公式
type Formula struct {
	Mp         int `bson:"mp" json:"mp"`
	Ammo       int `bson:"ammo" json:"ammo"`
	Mre        int `bson:"mre" json:"mre"`
	Part       int `bson:"part" json:"part"`
	InputLevel int `bson:"input_level" json:"input_level"`
}

// Stats 统计结果
type Stats struct {
	Formula Formula `json:"formula" bson:"formula"`
	ID      int     `json:"id,omitempty" bson:"id"`
	Type    int     `json:"type" bson:"type"`
	Count   int     `json:"count" bson:"count"`
	Date    int     `json:"date,omitempty" bson:"date"`
}

// FormulaStats 公式统计结果
type FormulaStats struct {
	Formula Formula `json:"formula" bson:"formula"`
	Type    int     `json:"type" bson:"type"`
	Count   int     `json:"count" bson:"count"`
	Date    int     `json:"date" bson:"date"`
}

// FormulaIndex
type FormulaIndex struct {
	Formula Formula `bson:"formula"`
	Type    int     `bson:"type"`
	ID      []int   `bson:"id"`
	Count   int     `bson:"count"`
}
