import pymongo
import json
import time


def datestr(ts: int) -> int:
    tarray = time.localtime(ts)
    return int(time.strftime('%Y%m%d', tarray))
    
DB_ADDRESS = '127.0.0.1'
DB_NAME = 'gfdb'

mc = pymongo.MongoClient(DB_ADDRESS)
db = mc.get_database(DB_NAME)

# gun_statistics
cold = db.get_collection('tdoll_record')
cnew = db.get_collection('tdoll_stats')


# from_time = 1495987200 # 2017/5/29 0:0:0
from_time = 1499097600 # 2017/7/4 0:0:0

dics = {}

for d in range(from_time, int(time.time()), 60 * 60 * 24):
    print("Calculating day {0}".format(datestr(d)))
    cur = cold.aggregate([
            {
                "$match": {
                    "dev_time": {
                        '$gte': d,
                        '$lt': d + 60 * 60 * 24
                    }
                }
            },
            {
                "$group": {
                    "_id": {
                        "formula": "$formula",
                        "id": "$gun_id"
                        # "eid": "$equip_id",
                        # "fid": "$fairy_id"
                    },
                    "count": {"$sum": 1}
                }
            }
    ])

    cur = list(cur)

    for i in cur:
        key = str(i['_id'])
        if dics.get(key) == None:
            ins = {}
            ins['formula'] = i['_id']['formula']
            ins['id'] = i['_id']['id']
            ins['type'] = 0
            # ins['id'] = i['_id']['eid'] if i['_id']['eid'] != 0 else i['_id']['fid']
            # ins['type'] = 1 if i['_id']['eid'] != 0 else 2 # 1 equip 2 fairy
            ins['count'] = i['count']
            dics[key] = ins
        else:
            dics[key]['count'] += i['count']
            # dics[key].pop('_id')
            
        for i in dics:
            dics[i]['date'] = datestr(d)

    print(len(dics))
    print("inserting")
    for i in dics.values():
        cnew.insert_one(i.copy())
        