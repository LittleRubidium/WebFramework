package jwt

import (
	"fmt"
	"testing"
)

func TestGetToken(t *testing.T) {
	token := GetToken(1)
	if token == "" {
		t.Fatal("cannot get token")
	}
	fmt.Println("success: " + token)
}
