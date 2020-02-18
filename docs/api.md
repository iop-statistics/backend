# IOP建造统计 API文档

## 基本格式
```json
{
  "success": true,
  "data": {},
  "error": "error reason"
}
```

| 字段 | 说明 |
| ---- | ---- |
| success | 本次请求是否请求到数据并且没出错|
| data | 数据，**其后只列出该字段的数据**|
| error | 错误信息，当success为false的时候存在|

## GET /stats/info

获得基本信息

### Response

```json
{
    "last_update": 1234567890,
    "tdoll_count": 1000,
  	"equip_count": 1000,
    "info": "通知信息"
}
```

| 字段        | 说明                         |
| ----------- | ---------------------------- |
| last_update | 10位时间戳，上次数据更新时间 |
| tdoll_count | 人形数据数                   |
| equip_count | 装备数据数(包括妖精)         |
| info        | 通知信息                     |



## GET /id

根据id获得出货数

### Parameter

x-www-form-urlencoded

| 字段 | 说明                                        |
| ---- | ------------------------------------------- |
| type | 类型，为tdoll/equip/fairy中的一个           |
| id   | 人形/装备/妖精的id                          |
| from | 起始时间(年月日的形式，如20170101)，默认为0 |
| to   | 截止时间(年月日的形式)，默认为当天          |



### Response

```json
{
    "actual_from": 20180101,
    "actual_to": 20180101,
    "data": [
        {
            "id": 1,
            "type": 0,
            "count": 10000,
            "total": 30000,
            "formula": {
                "mp": 130,
                "ammo": 130,
                "mre": 130,
                "part": 130,
                "input_level": 1
            }
        }
	]
}
```

| 字段         | 说明                            |
| ------------ | :------------------------------ |
| actual_from  | 实际起始时间                    |
| actual_to    | 实际结束时间                    |
| data.id      | 人形/装备/妖精的id              |
| data.type    | 类型，0为人形，1为装备，2为妖精 |
| data.count   | 出货数                          |
| data.total   | 该公式出货总数                  |
| data.formula | 公式                            |



## GET /formula

根据公式获得出货数

### Parameter

x-www-form-urlencoded

| 字段        | 说明                                        |
| ----------- | ------------------------------------------- |
| type        | 类型，为tdoll/equip/fairy中的一个           |
| mp          | 人力                                        |
| ammo        | 弹药                                        |
| mre         | 口粮                                        |
| part        | 零件                                        |
| input_level | 档位，0,1,2,3三档                           |
| from        | 起始时间(年月日的形式，如20170101)，默认为0 |
| to          | 截止时间(年月日的形式)，默认为当天          |



### Response

```json
{
    "actual_from": 20180101,
    "actual_to": 20180101,
    "data": [
        {
            "id": 1,
            "type": 0,
            "count": 10000
        }
	]
}
```

| 字段        | 说明                            |
| ----------- | ------------------------------- |
| actual_from | 实际起始时间                    |
| actual_to   | 实际结束时间                    |
| data.id     | 人形/装备/妖精的id              |
| data.type   | 类型，0为人形，1为装备，2为妖精 |
| data.count  | 出货数                          |



## GET /json/<json_name>.json

获得对应数据文件

| 文件名            | 说明             |
| ----------------- | ---------------- |
| gun_info          | 人形数据         |
| equip_info        | 装备数据         |
| fairy_info        | 妖精数据         |
| gun_info_simple   | 简化后的人形数据 |
| equip_info_simple | 简化后的装备数据 |