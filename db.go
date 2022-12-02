package dynamo_gorm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DB .
type DB struct {
	client *dynamodb.Client
	model  *model
	clone  bool
}

// Model .
func (db *DB) Model(i any) *DB {
	ret := db.getInstance()
	ret.model = newModel(i)
	return ret
}

func (db *DB) getInstance() *DB {
	if db.clone {
		return db
	}
	return &DB{client: db.client, clone: true}
}

// Query .
func (db *DB) Query(ctx context.Context, i any) error {
	qi, err := db.model.getQueryInput()
	if err != nil {
		return err
	}

	resp, err := db.client.Query(ctx, qi)
	if err != nil {
		return err
	}

	attributevalue.UnmarshalListOfMaps(resp.Items, i)
	return nil
}
