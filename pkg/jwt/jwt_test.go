package jwt

import (
	"fmt"
	"go_webapp/pkg/snowflake"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(snowflake.GenID())
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)
}

func TestParseToken(t *testing.T) {
	token, err := ParseToken("")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)
}
