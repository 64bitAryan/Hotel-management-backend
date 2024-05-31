package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/64bitAryan/hotel-management/db"
	"github.com/64bitAryan/hotel-management/types"
	"github.com/gofiber/fiber/v2"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "james@foo.com",
		Email:     "james",
		LastName:  "foo",
		Password:  "supersecurepassword",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticate(t *testing.T) {
	// setting up db and tearing down
	tdb := setup(t)
	defer tdb.teardown(t)
	insertTestUser(t, tdb.UserStore)

	// creating fiber application
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthentication)

	// creating params for request
	params := AuthParams{
		Email:    "james@foo.com",
		Password: "supersecurepassword",
	}
	b, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("unable to marshal params")
	}

	// creating the request with above params
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected the res statuscode to be 200 but got %d", resp.StatusCode)
	}

	// decoding the res.body(json) into an GO understandable entity.
	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	fmt.Println(resp)
}
