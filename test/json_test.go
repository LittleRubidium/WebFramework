package test

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"testing"
)

func TestJson(t *testing.T) {
	m := map[string]interface{}{
		"title":   "bot",
		"user_id": 1,
	}
	j, _ := json.Marshal(m)
	fmt.Println(string(j))
	p := gjson.ParseBytes(j)
	fmt.Println(p)
}
