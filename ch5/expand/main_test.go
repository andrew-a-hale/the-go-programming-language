package main

import (
	"strings"
	"testing"
)

func TestExpand(t *testing.T) {
	c := " $foo - $foo foo"
	e := " FOO - FOO foo"
	if g := Expand(c, strings.ToUpper); g != e {
		t.Errorf(`expected %s, got %s`, e, g)
	}
}