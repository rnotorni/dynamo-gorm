package dynamo_gorm

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"reflect"
	"regexp"
	"strings"
)

// HogeFuga example .
type HogeFuga struct {
	Key  string `dgorm:"p"`
	Sort string `dgorm:"s"`
}

func (h *HogeFuga) TableName() string {
	return "hoge"
}

type model struct {
	rv        reflect.Value
	tableName string

	pertitionKey *dgormField
	sortKey      *dgormField

	fields dgormFields
}

func newModel(i any) *model {
	rv := reflect.ValueOf(i)
	fields := newDgormFields(rv)
	return &model{rv: rv, tableName: getTableName(rv), pertitionKey: fields.getPartitionKey(), sortKey: fields.getSortKey(), fields: fields}
}

type dgormField struct {
	value reflect.Value
	tag   dgormTag
}

func newDgormField(rv reflect.Value, t reflect.StructTag) *dgormField {
	return &dgormField{value: rv, tag: newDgormTag(t.Get("gorm"))}
}

func (dt *dgormTag) has(s string) bool {
	_, exist := (*dt)[s]
	return exist
}

func (dt *dgormTag) isPartitionKey() bool {
	return dt.has("p")
}

func (dt *dgormTag) isSortKey() bool {
	return dt.has("s")
}

type dgormFields []*dgormField

func newDgormFields(rv reflect.Value) dgormFields {
	rt := rv.Elem().Type()
	ret := make([]*dgormField, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		ret[i] = newDgormField(rv.Elem().Field(i), rt.Field(i).Tag)
	}
	return ret
}

func (dfs *dgormFields) find(f func(*dgormField) bool) *dgormField {
	for _, df := range *dfs {
		if f(df) {
			return df
		}
	}
	return nil
}

func (dfs *dgormFields) getPartitionKey() *dgormField {
	return dfs.find(func(df *dgormField) bool { return df.tag.isPartitionKey() })
}

func (dfs *dgormFields) getSortKey() *dgormField {
	return dfs.find(func(df *dgormField) bool { return df.tag.isSortKey() })
}

type dgormTag map[string]string

// str is dgorm tag e.x. "p"
func newDgormTag(str string) dgormTag {
	ret := dgormTag{}
	ss := strings.Split(str, ",")
	for _, s := range ss {
		kv := strings.SplitN(s, "=", 2)
		if len(kv) == 2 {
			ret[kv[0]] = kv[1]
		} else {
			ret[kv[0]] = ""
		}
	}
	return ret
}

func getTableName(v reflect.Value) string {
	fValue := v.MethodByName("TableName")
	if !fValue.IsValid() {
		return getTableNameByStructName(v)
	}
	rets := fValue.Call(nil)
	if len(rets) != 1 {
		return getTableNameByStructName(v)
	}
	ret, ok := rets[0].Interface().(string)
	if !ok {
		return getTableNameByStructName(v)
	}
	return ret
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func getTableNameByStructName(v reflect.Value) string {
	ret := matchFirstCap.ReplaceAllString(v.Elem().Type().Name(), "${1}-${2}")
	matchAllCap.ReplaceAllString(ret, "${1}-${2}")
	return strings.ToLower(ret)
}

func (m *model) getQueryInput() (*dynamodb.QueryInput, error) {
	expr, err := expression.NewBuilder().WithKeyCondition(m.keyConditionBuilder()).Build()
	if err != nil {
		return nil, err
	}
	return &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		// TODO: IndexName: GSI or LSI
		TableName: aws.String(m.tableName),
	}, nil
}

func (m *model) keyConditionBuilder() expression.KeyConditionBuilder {
	ret := expression.KeyEqual(expression.Key(m.pertitionKey.tag["name"]), expression.Value(m.pertitionKey.value.Interface()))
	if m.sortKey != nil && !m.sortKey.value.IsZero() {
		ret = ret.And(expression.KeyEqual(expression.Key(m.sortKey.tag["name"]), expression.Value(m.sortKey.value.Interface())))
	}
	return ret
}

// generic使う?
func Query[T any](c *dynamodb.Client, i *T) []T {
	m := newModel(i)
	qi, _ := m.getQueryInput()

	response, _ := c.Query(context.TODO(), qi)
	ret := []T{}
	attributevalue.UnmarshalListOfMaps(response.Items, &ret)
	return ret
}

func main() {
	var client *dybamodb.Client

}
