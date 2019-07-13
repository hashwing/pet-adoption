package auth

import (
	"testing"
)

func Test_GetToken(t *testing.T) {
	token := GetToken("test")
	t.Log(token)
}
