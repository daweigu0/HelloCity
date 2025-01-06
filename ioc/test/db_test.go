package test

import (
	"HelloCity/ioc"
	"testing"
)

func TestInitDB(t *testing.T) {
	db := ioc.InitDB()
	if db == nil {
		t.Fail()
	}
}
