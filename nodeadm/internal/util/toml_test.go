package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	a := `[foo]
		bar = 'baz'
	`
	b := `[foo]
		baz = 'bar'
	`
	expected := `[foo]
		bar = 'baz'
		baz = 'bar'
		`
	merged, err := MergeTOML(Dedent(a), Dedent(b))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, Dedent(expected), *merged)
}
