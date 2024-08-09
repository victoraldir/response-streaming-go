package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderPage(t *testing.T) {
	sleepCalled := 0

	assert.Equal(t, 4, sleepCalled)

}
