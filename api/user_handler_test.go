package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/64bitAryan/hotel-management/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

const (
	testdbUri = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

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
		UserStore: db.NewMongoUserStore(client, dbname),
	}
}

func TestPostUser(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(db.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		Email:     "somefoo@test.com",
		FirstName: "foo",
		LastName:  "bar",
		Password:  "1234567",
	}

	b, _ := json.Marshal(params)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-type", "application/json")
	res, _ := app.Test(req)
	fmt.Println(res.Status)
}
