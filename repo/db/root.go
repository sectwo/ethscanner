package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"scanner/env"
	"scanner/types"
	. "scanner/utils"
)

type DB struct {
	env *env.Env

	client *mongo.Client
	db     *mongo.Database

	block *mongo.Collection
	tx    *mongo.Collection
}

type DBImpl interface {
	SaveBlock(b *types.CustomBlock) error
	SaveTxByBulk(model []mongo.WriteModel) error
	SaveTx(t *types.CustomTx) error
}

func NewDB(env *env.Env) (DBImpl, error) {
	d := &DB{
		env: env,
	}

	ctx := Context()
	var err error

	if d.client, err = mongo.Connect(ctx, options.Client().ApplyURI(env.DB.Uri)); err != nil {
		panic(err)
	} else if err = d.client.Ping(ctx, nil); err != nil {
		panic(err)
	} else {
		d.db = d.client.Database(env.DB.DB)
		d.block = d.db.Collection(env.DB.Block)
		d.tx = d.db.Collection(env.DB.Tx)

		return d, nil
	}
}

func (d *DB) SaveBlock(b *types.CustomBlock) error {
	filter := bson.M{"blockNumber": 1}

	if j, err := ToJson(b); err != nil {
		return err
	} else if result, err := d.block.UpdateOne(Context(), filter, bson.M{"$set": j}, options.Update().SetUpsert(true)); err != nil {
		return err
	} else {
		log.Print("Success to save BLock : ", result.UpsertedCount, result.ModifiedCount)
		return nil
	}
}

func (d *DB) SaveTxByBulk(model []mongo.WriteModel) error {
	if result, err := d.tx.BulkWrite(Context(), model); err != nil {
		return err
	} else {
		log.Println("Success to save Tx : ", result.UpsertedCount, result.ModifiedCount)
		return nil
	}
}

func (d *DB) SaveTx(t *types.CustomTx) error {
	filter := bson.M{"hash": 1}

	if j, err := ToJson(t); err != nil {
		return err
	} else if result, err := d.tx.UpdateOne(Context(), filter, bson.M{"$set": j}, options.Update().SetUpsert(true)); err != nil {
		return err
	} else {
		log.Print("Success to save Tx : ", result.UpsertedCount, result.ModifiedCount)
		return nil
	}
}
