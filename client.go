package dynamo_gorm

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func Query[T any](c *dynamodb.Client, i *T, ctx context.Context) []T {
	m := newModel(i)
	qi, _ := m.getQueryInput()

	response, _ := c.Query(ctx, qi)
	ret := []T{}
	attributevalue.UnmarshalListOfMaps(response.Items, &ret)
	return ret
}
