package main

import "testing"

func TestSet(t *testing.T) {
	e := set("test", "ok", 100)

	if e != nil {
		t.Error(e)
	}
}

func TestGet(t *testing.T) {
	v := get("test")

	if v != "ok" {
		t.Error(v)
	}
}

func TestDelete(t *testing.T) {
	e := delete("test")

	if e != nil {
		t.Error(e)
	}
}
