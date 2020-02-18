package model

import "go.mongodb.org/mongo-driver/bson"

// TimeRecord 新加人形、装备等的加入时间
type TimeRecord struct {
	Type int `json:"time" bson:"time"`
	ID   int `json:"id" bson:"id"`
	Date int `json:"date" bson:"date"`
}

const (
	ColTimeRecord = "time_record"
)

func (m *model) GetTimeRecordByIDAndType(id, t int) (TimeRecord, error) {
	c := m.db.Collection(ColTimeRecord)
	query := bson.M{
		"id":   id,
		"type": t,
	}

	var ret TimeRecord
	err := c.FindOne(m.ctx, query).Decode(&ret)

	return ret, err
}
