package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMethod(t *testing.T) {
	input := `
		let five = 10;
	`

	program := getProgram(t, input)

	assert.Equal(t, "let five = 10", program.String())
}
