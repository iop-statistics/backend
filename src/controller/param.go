package controller

import "github.com/Lyt99/iop-statistics/model"

type ParamGetID struct {
	Type     string `form:"type" binding:"required,eq=tdoll|eq=equip|eq=fairy"` //TODO: 这个参数验证有问题
	ID       int    `form:"id" binding:"required"`
	FromTime int    `form:"from" binding:"omitempty"`
	ToTime   int    `form:"to" binding:"omitempty,gtefield=FromTime"`
}

type ParamGetFormula struct {
	Type       string `form:"type" binding:"required,eq=tdoll|eq=equip|eq=fairy"`
	Mp         int    `form:"mp" binding:"required,gte=10,lte=9999"`
	Ammo       int    `form:"ammo" binding:"required,gte=10,lte=9999"`
	Mre        int    `form:"mre" binding:"required,gte=10,lte=9999"`
	Part       int    `form:"part" binding:"required,gte=10,lte=9999"`
	InputLevel int    `form:"input_level" binding:"gte=0,lte=3"`
	FromTime   int    `form:"from" binding:"omitempty"`
	ToTime     int    `form:"to" binding:"omitempty,gtefield=FromTime"`
}

type ResponseInfo struct {
	LastUpdate int64  `json:"last_update"`
	TdollCount int64  `json:"tdoll_count"`
	EquipCount int64  `json:"equip_count"`
	Info       string `json:"info"`
}

// StatisticsRow
type StatisticsRow struct {
	Formula model.Formula `json:"formula,omitempty"`
	Count   int           `json:"count,omitempty"`
	Total   int           `json:"total,omitempty"`
}

// StatisticsResult
type StatisticsResult struct {
	ActualFrom int           `json:"actual_from"`
	ActualTo   int           `json:"actual_to"`
	Data       []interface{} `json:"data"`
}

type StatisticsRowFormula struct {
	ID    int `json:"id,omitempty"`
	Type  int `json:"type,omitempty"`
	Count int `json:"count,omitempty"`
}
