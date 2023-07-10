package test

import (
	"testing"
	"tutorial-restfulapi/simple"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	conn, cleanup := simple.InitializedConnection("database")
	assert.NotNil(t, conn)

	cleanup()
}
