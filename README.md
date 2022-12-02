# dynamo-gorm
WIP dynamo-go

```golang
import (
	dgorm "github.com/.../dynamo-gorm"
)
// HogeFuga example .
type HogeFuga struct {
	Key  string `dgorm:"p"`
	Sort string `dgorm:"s"`
}

func (h *HogeFuga) TableName() string {
	return "hoge"
}

func main() {
	var client *dybamodb.Client
	h := &HogeFuga{Key:"a", Sort:"b"}
	hs := dgorm.Query[HogeFuga](client, &h)
}
```
