package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teohrt/cruddyAPI/dbclient"
)

func TestNewService(t *testing.T) {
	svc := New(&dbclient.Config{
		DynamoDBTableName:  "this",
		AWSRegion:          "is",
		AWSSessionEndpoint: "a test",
	})

	assert.IsType(t, serviceImpl{}, svc)
}
