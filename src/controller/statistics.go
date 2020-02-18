package controller

import (
	"fmt"
	"github.com/Lyt99/iop-statistics/model"
	"github.com/Lyt99/iop-statistics/util"
	"github.com/Lyt99/iop-statistics/util/context"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetStatsInfo(c *gin.Context) {
	var ret ResponseInfo

	m := model.GetModel()
	defer m.Close()

	v, err := m.GetInfoByKey(model.KeyInfo)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	ret.Info = v.(string)

	v, err = m.GetInfoByKey(model.KeyLastUpdate)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	ret.LastUpdate = int64(v.(int32))

	v, err = m.GetInfoByKey(model.KeyCountEquip)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	ret.EquipCount = int64(v.(int32))

	v, err = m.GetInfoByKey(model.KeyCountTdoll)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	ret.TdollCount = int64(v.(int32))

	context.Success(c, ret)
}

func GetRecordsByID(c *gin.Context) {
	var param ParamGetID
	if err := c.ShouldBind(&param); err != nil {
		context.Error(c, http.StatusBadRequest, "Bad request", err)
	}

	// 判断查询数据类型(人形、装备、精灵)
	m := model.GetModel()
	defer m.Close()

	// 类型
	var t int

	switch param.Type {
	case "tdoll":
		{
			t = model.TypeTdoll
		}
	case "equip":
		{
			t = model.TypeEquip
		}
	case "fairy":
		{
			t = model.TypeFairy
		}
	default:
		{
			context.Error(c, http.StatusBadRequest, "error type", nil)
			return
		}
	}

	today := util.GetDateToday()
	if param.ToTime == 0 || param.ToTime > today {
		param.ToTime = today
	}

	if param.FromTime == 0 {
		tr, err := m.GetTimeRecordByIDAndType(param.ID, t)
		if err != nil {
			param.FromTime = 0
		} else {
			param.FromTime = tr.Date
		}
	}

	key := fmt.Sprintf("id:%d:%v", param.ToTime, param)
	if context.TryResponseCache(c, key) {
		return
	}

	f, err := m.GetValidFormulasByID(param.ID, t)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	var records StatisticsResult
	records.ActualFrom = param.FromTime
	records.ActualTo = param.ToTime

	for _, v := range f {
		var row StatisticsRow
		row.Formula = v
		row.Count, err = m.GetFormulaAndIDCount(v, param.ID, t, param.FromTime, param.ToTime)
		if err != nil {
			context.Error(c, http.StatusInternalServerError, "internal server error", err)
			return
		}

		if row.Count == 0 { // 不存在
			continue
		}

		row.Total, err = m.GetFormulaCount(v, t, param.FromTime, param.ToTime)
		if err != nil {
			context.Error(c, http.StatusInternalServerError, "internal server error", err)
			return
		}

		countAnother := 0
		if t == model.TypeEquip {
			countAnother, err = m.GetFormulaCount(v, model.TypeFairy, param.FromTime, param.ToTime)
		} else if t == model.TypeFairy {
			countAnother, err = m.GetFormulaCount(v, model.TypeEquip, param.FromTime, param.ToTime)
		}

		if err != nil {
			context.Error(c, http.StatusInternalServerError, "internal server error", err)
			return
		}

		row.Total += countAnother

		records.Data = append(records.Data, row)
	}

	context.SuccessWithCache(c, key, records)
}

func GetRecordsByFormula(c *gin.Context) {
	var param ParamGetFormula
	if err := c.ShouldBind(&param); err != nil {
		context.Error(c, http.StatusBadRequest, "Bad request", err)
	}

	// 判断查询数据类型(人形、装备、精灵)
	m := model.GetModel()
	defer m.Close()

	// 类型
	var t int

	switch param.Type {
	case "tdoll":
		{
			t = model.TypeTdoll
		}
	case "equip":
		{
			t = model.TypeEquip
		}
	case "fairy":
		{
			t = model.TypeFairy
		}
	default:
		{
			context.Error(c, http.StatusBadRequest, "error type", nil)
			return
		}
	}

	today := util.GetDateToday()
	if param.ToTime == 0 || param.ToTime > today {
		param.ToTime = today
	}

	key := fmt.Sprintf("fo:%d:%v", param.ToTime, param)
	if context.TryResponseCache(c, key) {
		return
	}

	f := model.Formula{
		Mp:         param.Mp,
		Ammo:       param.Ammo,
		Mre:        param.Mre,
		Part:       param.Part,
		InputLevel: param.InputLevel,
	}

	var records StatisticsResult
	records.ActualFrom = param.FromTime
	records.ActualTo = param.ToTime

	i, err := m.GetIDsByFormula(f, t)
	if err != nil {
		context.Error(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	for _, v := range i {
		var row StatisticsRowFormula

		row.ID = v
		row.Type = t
		row.Count, err = m.GetFormulaAndIDCount(f, v, t, param.FromTime, param.ToTime)
		if err != nil {
			context.Error(c, http.StatusInternalServerError, "internal server error", err)
			return
		}

		if row.Count == 0 {
			continue
		}

		records.Data = append(records.Data, row)
	}

	if t != model.TypeTdoll { // 合并查询装备和妖精
		if t == model.TypeEquip {
			t = model.TypeFairy
		} else {
			t = model.TypeEquip
		}

		i, err := m.GetIDsByFormula(f, t)
		if err != nil {
			context.Error(c, http.StatusInternalServerError, "internal server error", err)
			return
		}

		for _, v := range i {
			var row StatisticsRowFormula

			row.ID = v
			row.Type = t
			row.Count, err = m.GetFormulaAndIDCount(f, v, t, param.FromTime, param.ToTime)
			if err != nil {
				context.Error(c, http.StatusInternalServerError, "internal server error", err)
				return
			}

			if row.Count == 0 {
				continue
			}

			records.Data = append(records.Data, row)
		}
	}

	context.SuccessWithCache(c, key, records)
}
