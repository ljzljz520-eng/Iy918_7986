package main

import "testing"

func TestBuildLabel(t *testing.T) {
	if buildLabel() == "" {
		t.Fatal("empty label")
	}
}
