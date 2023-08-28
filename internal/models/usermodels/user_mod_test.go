package usermodels

import (
	"fmt"
	"testing"
)

func TestXMUserWithAccount(t *testing.T) {
	xmUser, err := XMUserWithAccount("123")
	if err != nil {
		t.Errorf(err.Error())
	}
	fmt.Printf("%+v\n", xmUser)
}
