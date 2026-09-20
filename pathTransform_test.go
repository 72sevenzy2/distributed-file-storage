package main

import "testing"

func TestPathTransformFunc(t *testing.T) {
	key := "uhhihello"
	newk := TransformPathFunc(key)
	expectedKey := "a5637/06f04/aa87e/87aa7/4dbe2/c9742/d0735/9fd40"
	if newk != expectedKey {
		t.Errorf("want %s, got %s", expectedKey, newk)
	}
}
