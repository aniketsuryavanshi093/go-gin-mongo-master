package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type SchemaServiceImpl struct {
	schemacollection *mongo.Collection
	ctx              context.Context
}

func NewSchemaService(schemacollection *mongo.Collection, ctx context.Context) SchemaService {
	return &SchemaServiceImpl{
		schemacollection: schemacollection,
		ctx:              ctx,
	}
}
