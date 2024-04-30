package api

import (
	"context"
	"log"
	"testing"

	"github.com/64bitAryan/hotel-management/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

const testdbUri = "mongodb://localhost:27017"

func (tdb *testdb) teardown(t *testing.T) {
	ctx := context.Background()
	if err := tdb.UserStore.Drop(ctx); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdbUri))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
}
