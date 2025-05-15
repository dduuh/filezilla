package mong

import (
	"context"
	"filezilla/internal/configs"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDB(cfg *configs.Config, ctx context.Context) (*mongo.Database, error) {
	opts := options.Client()
	opts.SetAuth(options.Credential{
		Username: cfg.MongoCfg.Username,
		Password: cfg.MongoCfg.Password,
	})

	opts.ApplyURI(cfg.MongoCfg.Uri)

	db, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx, nil); err != nil {
		return nil, err
	}

	mondoDb := db.Database(cfg.MongoCfg.Database)

	return mondoDb, nil
}
