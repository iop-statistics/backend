package tool

import (
	"fmt"
	"github.com/Lyt99/iop-statistics/model"
	"github.com/Lyt99/iop-statistics/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"runtime"
	"time"
)

const (
	fromTime = 1499097600 // 2017/7/4 0:0:0
)

func MigrateStats() {
	ctx := context.Background()
	_ = model.MongoClient.Connect(ctx)

	db := model.MongoClient.Database(model.DBName)

	now := time.Now().Unix()

	type stats struct {
		Formula model.Formula `bson:"formula"`
		ID      int           `bson:"id"`
		Type    int           `bson:"type"`
		Count   int           `bson:"count"`
		Date    *int          `bson:"date"`
	}

	type aggregateResult struct {
		ID struct {
			Formula model.Formula `bson:"formula"`
			GunID   int           `bson:"gid"`
			EquipID int           `bson:"eid"`
			FairyID int           `bson:"fid"`
		} `bson:"_id"`
		Count int `bson:"count"`
	}

	dictStats := make(map[string]*stats)

	date := 0

	// 人形
	log.Println("Migrating Tdoll records.")
	collectionOld := db.Collection(model.ColTdollRecords)
	collectionNew := db.Collection(model.ColStats)

	_ = collectionNew.Drop(ctx)

	for d := int64(fromTime); d <= now; d += 60 * 60 * 24 {
		var inserts []interface{}

		date = util.GetDate(time.Unix(d, 0))

		log.Printf("Calculating day %d\n", date)

		pipe := mongo.Pipeline{
			{
				{
					"$match",
					bson.M{
						"dev_time": bson.M{
							"$gte": d,
							"$lt":  d + 60*60*24,
						},
					},
				},
			},
			{
				{
					"$group",
					bson.M{
						"_id": bson.M{
							"formula": "$formula",
							"gid":     "$gun_id",
						},
						"count": bson.M{"$sum": 1},
					},
				},
			},
		}

		cur, err := collectionOld.Aggregate(ctx, pipe)
		if err != nil {
			panic(err)
		}

		var r aggregateResult
		for cur.Next(ctx) {
			err = cur.Decode(&r)
			if err != nil {
				panic(err)
			}

			// 将数据加入到dict中
			key := fmt.Sprintf("%v", r.ID)
			ins, ok := dictStats[key]
			if ok {
				dictStats[key].Count += r.Count // 如果存在就加数量
			} else {
				n := new(stats)
				n.ID = r.ID.GunID
				n.Count = r.Count
				n.Formula = r.ID.Formula
				n.Type = model.TypeTdoll
				n.Date = &date

				dictStats[key] = n
				ins = n
			}

			inserts = append(inserts, ins)
		}

		log.Printf("Inserting %d records\n", len(inserts))
		if len(inserts) == 0 {
			log.Println("Skipped")
			continue
		}
		_, err = collectionNew.InsertMany(ctx, inserts)
		if err != nil {
			panic(err)
		}

	}

	dictStats = nil

	runtime.GC() // 手动回收一下内存

	dictStats = make(map[string]*stats)

	// 装备和妖精的
	log.Println("Migrating Fairy&Equip records.")

	collectionOld = db.Collection(model.ColEquipRecords)

	for d := int64(fromTime); d <= now; d += 60 * 60 * 24 {
		var inserts []interface{}

		date = util.GetDate(time.Unix(d, 0))

		log.Printf("Calculating day %d\n", date)

		pipe := mongo.Pipeline{
			{
				{
					"$match",
					bson.M{
						"dev_time": bson.M{
							"$gte": d,
							"$lt":  d + 60*60*24,
						},
					},
				},
			},
			{
				{
					"$group",
					bson.M{
						"_id": bson.M{
							"formula": "$formula",
							"eid":     "$equip_id",
							"fid":     "$fairy_id",
						},
						"count": bson.M{"$sum": 1},
					},
				},
			},
		}

		cur, err := collectionOld.Aggregate(ctx, pipe)
		if err != nil {
			panic(err)
		}

		var r aggregateResult
		for cur.Next(ctx) {
			err = cur.Decode(&r)
			if err != nil {
				panic(err)
			}

			// 将数据加入到dict中
			key := fmt.Sprintf("%v", r.ID)
			ins, ok := dictStats[key]
			if ok {
				dictStats[key].Count += r.Count // 如果存在就加数量
			} else {
				n := new(stats)
				n.Count = r.Count
				n.Formula = r.ID.Formula
				n.Date = &date

				if r.ID.EquipID != 0 { // equip
					n.Type = model.TypeEquip
					n.ID = r.ID.EquipID
				} else {
					n.Type = model.TypeFairy
					n.ID = r.ID.FairyID
				}

				dictStats[key] = n
				ins = n
			}

			inserts = append(inserts, ins)
		}

		log.Printf("Inserting %d records\n", len(inserts))
		if len(inserts) == 0 {
			log.Println("Skipped")
			continue
		}
		_, err = collectionNew.InsertMany(ctx, inserts)
		if err != nil {
			panic(err)
		}
	}

	// 索引，应该只用date和id_with_type
	log.Println("Ensuring index")
	key1 := "id_type_formula_date"

	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"id", 1},
				{"type", 1},
				{"formula", 1},
				{"date", -1},
			},
			Options: &options.IndexOptions{
				Name: &key1,
			},
		},
	}
	c := db.Collection(model.ColStats)
	_, err := c.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	log.Println("Done.")
}

func MigrateFormula() {
	ctx := context.Background()
	_ = model.MongoClient.Connect(ctx)

	db := model.MongoClient.Database(model.DBName)

	now := time.Now().Unix()

	type stats struct {
		Formula model.Formula `bson:"formula"`
		Type    int           `bson:"type"`
		Count   int           `bson:"count"`
		Date    *int          `bson:"date"`
	}

	type aggregateResult struct {
		ID struct {
			Formula model.Formula `bson:"formula"`
		} `bson:"_id"`
		Count int `bson:"count"`
	}

	dictStats := make(map[string]*stats)

	date := 0

	// 人形
	log.Println("Migrating Tdoll records.")
	collectionOld := db.Collection(model.ColTdollRecords)
	collectionNew := db.Collection(model.ColFormula)

	_ = collectionNew.Drop(ctx)

	for d := int64(fromTime); d <= now; d += 60 * 60 * 24 {
		var inserts []interface{}

		date = util.GetDate(time.Unix(d, 0))

		log.Printf("Calculating day %d\n", date)

		pipe := mongo.Pipeline{
			{
				{
					"$match",
					bson.M{
						"dev_time": bson.M{
							"$gte": d,
							"$lt":  d + 60*60*24,
						},
					},
				},
			},
			{
				{
					"$group",
					bson.M{
						"_id": bson.M{
							"formula": "$formula",
						},
						"count": bson.M{"$sum": 1},
					},
				},
			},
		}

		cur, err := collectionOld.Aggregate(ctx, pipe)
		if err != nil {
			panic(err)
		}

		var r aggregateResult
		for cur.Next(ctx) {
			err = cur.Decode(&r)
			if err != nil {
				panic(err)
			}

			// 将数据加入到dict中
			key := fmt.Sprintf("%v", r.ID)
			ins, ok := dictStats[key]
			if ok {
				dictStats[key].Count += r.Count // 如果存在就加数量
			} else {
				n := new(stats)
				n.Count = r.Count
				n.Formula = r.ID.Formula
				n.Type = model.TypeTdoll
				n.Date = &date

				dictStats[key] = n
				ins = n
			}

			inserts = append(inserts, ins)
		}

		log.Printf("Inserting %d records\n", len(inserts))
		if len(inserts) == 0 {
			log.Println("Skipped")
			continue
		}
		_, err = collectionNew.InsertMany(ctx, inserts)
		if err != nil {
			panic(err)
		}

	}

	dictStats = nil

	runtime.GC() // 手动回收一下内存

	dictStats = make(map[string]*stats)

	// 装备和妖精的
	log.Println("Migrating Fairy&Equip records.")

	collectionOld = db.Collection(model.ColEquipRecords)

	for d := int64(fromTime); d <= now; d += 60 * 60 * 24 {
		var inserts []interface{}

		date = util.GetDate(time.Unix(d, 0))

		log.Printf("Calculating day %d\n", date)

		pipe := mongo.Pipeline{
			{
				{
					"$match",
					bson.M{
						"dev_time": bson.M{
							"$gte": d,
							"$lt":  d + 60*60*24,
						},
					},
				},
			},
			{
				{
					"$group",
					bson.M{
						"_id": bson.M{
							"formula": "$formula",
						},
						"count": bson.M{"$sum": 1},
					},
				},
			},
		}

		cur, err := collectionOld.Aggregate(ctx, pipe)
		if err != nil {
			panic(err)
		}

		var r aggregateResult
		for cur.Next(ctx) {
			err = cur.Decode(&r)
			if err != nil {
				panic(err)
			}

			// 将数据加入到dict中
			key := fmt.Sprintf("%v", r.ID)
			ins, ok := dictStats[key]
			if ok {
				dictStats[key].Count += r.Count // 如果存在就加数量
			} else {
				n := new(stats)
				n.Count = r.Count
				n.Formula = r.ID.Formula
				n.Date = &date
				n.Type = model.TypeEquip // 实际上是Equip和Fairy

				dictStats[key] = n
				ins = n
			}

			inserts = append(inserts, ins)
		}

		log.Printf("Inserting %d records\n", len(inserts))
		if len(inserts) == 0 {
			log.Println("Skipped")
			continue
		}
		_, err = collectionNew.InsertMany(ctx, inserts)
		if err != nil {
			panic(err)
		}
	}

	// 索引，应该只用date和id_with_type
	log.Println("Ensuring index")
	key1 := "formula_type_date"

	indexModel := []mongo.IndexModel{
		{
			Keys: bson.D{
				{"formula", 1},
				{"type", 1},
				{"date", -1},
			},
			Options: &options.IndexOptions{
				Name: &key1,
			},
		},
	}
	c := db.Collection(model.ColFormula)
	_, err := c.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	log.Println("Done.")
}

func IndexFormula() {
	ctx := context.Background()
	_ = model.MongoClient.Connect(ctx)

	db := model.MongoClient.Database(model.DBName)

	cIndex := db.Collection(model.ColFormulaIndex)

	_ = cIndex.Drop(ctx)
	_, err := cIndex.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"formula": 1,
			"type":    1,
			"count":   -1,
		},
		Options: nil,
	})

	if err != nil {
		panic(err)
	}

	// 人形部分
	log.Println("Indexing tdoll formula")
	index(db, ctx, model.ColStats, model.ColFormulaIndex, model.TypeTdoll)

	// 装备部分
	log.Println("Indexing equip formula")
	index(db, ctx, model.ColStats, model.ColFormulaIndex, model.TypeEquip)

	// 妖精部分
	log.Println("Indexing fairy formula")
	index(db, ctx, model.ColStats, model.ColFormulaIndex, model.TypeFairy)

	_, err = cIndex.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"id":    1,
			"type":  1,
			"count": -1,
		},
		Options: nil,
	})

}

func index(db *mongo.Database, ctx context.Context, from, to string, t int) {
	cStats := db.Collection(from)
	cFormula := db.Collection(to)

	ids, err := cStats.Distinct(ctx, "id", bson.M{"type": t})
	if err != nil {
		panic(err)
	}

	for _, v := range ids {
		log.Printf("Indexing id %d\n", v)
		filter := bson.M{
			"id":   v,
			"type": t,
		}

		fs, err := cStats.Distinct(ctx, "formula", filter)
		if err != nil {
			panic(err)
		}

		for _, f := range fs {
			// 获得该公式在记录中的最大数量
			countQuery := bson.M{
				"id":      v,
				"type":    t,
				"formula": f,
			}

			var countResult model.Stats
			err := cStats.FindOne(ctx, countQuery).Decode(&countResult) // 按照索引来查询的话，应该会查到数量最大的
			if err != nil {
				panic(err)
			}

			query := bson.M{
				"formula": f,
				"type":    t,
			}

			update := bson.M{
				"$addToSet": bson.M{
					"id": v,
				},
				"$inc": bson.M{
					"count": countResult.Count,
				},
			}

			t := true
			_, err = cFormula.UpdateOne(ctx, query, update, &options.UpdateOptions{
				Upsert: &t,
			})

			if err != nil {
				panic(err)
			}
		}
	}

	_, err = cFormula.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.M{
			"id":   1,
			"type": 1,
		},
		Options: nil,
	})

	if err != nil {
		panic(err)
	}
}
