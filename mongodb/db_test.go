package mongodb

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateSuccessMongoConnection(t *testing.T) {
	client, collection, err := NewMongoDBConnection(os.Getenv("MONGO_DB_USER"), os.Getenv("MONGO_DB_PASS"), os.Getenv("MONGO_DB"), os.Getenv("MONGO_DB_COLLECTION"))

	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, client, "Client should not be nill")
	assert.NotNil(t, collection, "Collection should not be nill")
}
