package model

import "go.mongodb.org/mongo-driver/bson"

type KVPair struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
}

const (
	ColInfo       = "info"
	KeyLastUpdate = "last_update"
	KeyInfo       = "info"
	KeyCountTdoll = "count_tdoll"
	KeyCountEquip = "count_equip"
)

func (m *model) GetInfoByKey(key string) (interface{}, error) {
	c := m.db.Collection(ColInfo)

	query := bson.M{
		"key": key,
	}

	var ret KVPair
	err := c.FindOne(m.ctx, query).Decode(&ret)

	return ret.Value, err
}

func (m *model) GetColCount(col string) (int64, error) {
	c := m.db.Collection(col)

	return c.CountDocuments(m.ctx, bson.M{})
}
