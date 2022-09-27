package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestTestImplementation(t *testing.T) {
	*wordListMultiplier = 2
	words := [][]byte{[]byte("test"), []byte("with"), []byte("items")}
	var inserted [][]byte
	insert := func(word []byte) bool {
		inserted = append(inserted, word)
		return true
	}
	contains := func([]byte) bool { return true }

	got := testImplementation(words, 0, insert, contains)

	want := filterStats{
		insertFailed: 0,
		tp:           3,
		fp:           1,
		tn:           0,
		fn:           0,
	}
	if !cmp.Equal(want, got,
		cmp.AllowUnexported(filterStats{}),
		cmpopts.IgnoreFields(filterStats{}, "mem")) {
		t.Errorf("testImplementation got %v, want %v", got, want)
		t.Logf("Inserted: %v", inserted)
	}
}
