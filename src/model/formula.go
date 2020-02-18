package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *model) GetFormulaCount(f Formula, t, from, to int) (int, error) {
	if t == TypeFairy {
		t = TypeEquip
	}

	countTo, d, err := m.getSingleFormulaLeftBoundaryCount(f, t, to)
	if err != nil {
		return 0, err
	}

	if d < from { // 没有记录
		return 0, nil
	}

	countFrom, _, err := m.getSingleFormulaLeftBoundaryCount(f, t, from-1) // 这里要-1，因为左侧是开区间
	if err != nil {
		return 0, err
	}

	return countTo - countFrom, nil
}

func (m *model) getSingleFormulaLeftBoundaryCount(f Formula, t, date int) (int, int, error) {
	if date <= 0 {
		return 0, 0, nil
	}

	query := bson.D{
		{"formula", f},
		{"type", t},
		{"date", bson.M{"$lte": date}},
	}

	var r Stats
	c := m.db.Collection(ColFormula)

	err := c.FindOne(m.ctx, query, &options.FindOneOptions{
		//Sort: bson.M{"date": -1},
	}).Decode(&r)

	if err == mongo.ErrNoDocuments {
		return 0, 0, nil
	}

	if err != nil {
		return 0, 0, err
	}

	return r.Count, r.Date, nil
}
