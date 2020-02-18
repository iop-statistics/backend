package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *model) GetFormulaCountBrutally(f Formula, t int, from, to int) (int, error) {
	count := 0

	ids, err := m.GetIDsByFormula(f, t)
	if err != nil {
		return 0, err
	}

	for _, id := range ids {
		c, err := m.GetFormulaAndIDCount(f, id, t, from, to)
		if err != nil {
			return 0, err
		}

		count += c
	}

	return count, nil
}

func (m *model) GetFormulaAndIDCount(f Formula, id int, t int, from, to int) (int, error) {
	countTo, d, err := m.getSingleStatsLeftBoundaryCount(id, t, f, to)
	if err != nil {
		return 0, err
	}

	if d < from { // 没有记录
		return 0, nil
	}

	countFrom, _, err := m.getSingleStatsLeftBoundaryCount(id, t, f, from-1) // 这里要-1，因为左侧是开区间
	if err != nil {
		return 0, err
	}

	return countTo - countFrom, nil
}

// getSingleStatsLeftBoundaryCount returns count date and error
func (m *model) getSingleStatsLeftBoundaryCount(id, t int, f Formula, date int) (int, int, error) {
	if date <= 0 {
		return 0, 0, nil
	}

	query := bson.D{
		{"id", id},
		{"type", t},
		{"formula", f},
		{"date", bson.M{"$lte": date}},
	}

	var r Stats
	c := m.db.Collection(ColStats)

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
