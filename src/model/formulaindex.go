package model

import (
	"github.com/Lyt99/iop-statistics/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *model) GetValidFormulasByID(id, t int) ([]Formula, error) {
	query := bson.D{
		{"id", id},
		{"type", t},
	}

	c := m.db.Collection(ColFormulaIndex)

	cur, err := c.Find(m.ctx, query, &options.FindOptions{
		Sort: bson.M{"count": -1},
		Projection: bson.M{
			"formula": 1,
			"count":   1,
		},
	})

	if err != nil {
		return []Formula{}, err
	}

	var ret []Formula

	for cur.Next(m.ctx) && len(ret) <= config.GlobalConfig.Statistics.MaxRecordCount {
		var idx FormulaIndex
		err := cur.Decode(&idx)
		if err != nil {
			return ret, err
		}

		if idx.Count < config.GlobalConfig.Statistics.RecordThreshold &&
			len(ret) > config.GlobalConfig.Statistics.MinRecordCount {
			break
		}

		ret = append(ret, idx.Formula)
	}

	return ret, nil
}

func (m *model) GetIDsByFormula(f Formula, t int) ([]int, error) {
	query := bson.D{
		{"formula", f},
		{"type", t},
	}

	var res FormulaIndex

	c := m.db.Collection(ColFormulaIndex)

	err := c.FindOne(m.ctx, query).Decode(&res)

	if err == mongo.ErrNoDocuments {
		return []int{}, nil
	}

	return res.ID, err
}
