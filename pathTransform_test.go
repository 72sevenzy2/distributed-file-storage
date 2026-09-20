package main

import "testing"

func TestPathTransformFunc(t *testing.T) {
	key := "uhhihello"
	newk := TransformPathFunc(key)
	expectedKey := "a5637/06f04/aa87e/87aa7/4dbe2/c9742/d0735/9fd40"
	expectedOriginalKey := "a563706f04aa87e87aa74dbe2c9742d07359fd40"
	if newk.FileName != expectedKey {
		t.Errorf("want %s, got %s", expectedKey, newk.FileName)
	}
	if newk.Original != expectedOriginalKey {
		t.Errorf("want %s, got %s", expectedOriginalKey, newk.Original)
	}
}
