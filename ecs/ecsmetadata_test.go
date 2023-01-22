package ecs

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientToMetadataEndpoint(t *testing.T) {
	os.Setenv("ECS_CONTAINER_METADATA_URI_V4", "http://example.com")
	client1 := NewClient("http://example.com")
	client2, err := NewClientToMetadataEndpoint()

	assert.Equal(t, nil, err)
	assert.Equal(t, client1, client2)
}

func TestRetriveTaskMetadata(t *testing.T) {
	ctx := context.Background()
	c := Client{}
	c.endpoint = "http://example.com"
	task := TaskMetadata{
		Cluster: "test",
	}
	taskMetadata, err := c.RetriveTaskMetadata(ctx)

	assert.NotEqual(t, task, taskMetadata)
	assert.NotEqual(t, nil, err)
}
