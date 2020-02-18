# IOP建造统计 数据库设计文档

## Collection列表

- info

  存放诸如更新时间、公告等信息，KV键值对的形式

  model.KVPair

- date

  存放新加入人形/装备/妖精的加入时间
  model.TimeRecord
  
- tdoll_record
  
  单条人形出货记录
  
  model.Record
  
- equip_record

  单条装备出货记录

  model.Record

- tdoll_stats

  汇总后人形出货记录

  model.StatsRecord

- equip_stats

  汇总后装备出货记录

  model.StatsRecord

