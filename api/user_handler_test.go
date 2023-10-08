package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/fulltimegodev/hotel-reservation/db"
	"github.com/fulltimegodev/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	testmongouri = "mongodb://localhost:27017"
	dbname = "hotel-reservation-test"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown (t *testing.T){
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setUp(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testmongouri))
	
	if err != nil {
		log.Fatal(err)
	}

	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
} 

func TestPostUser (t * testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)

	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams {
		Email: "some@foo.com",
		FirstName: "James",
		LastName: "Foo",
		Password: "aisjaisjiasjiajsiasjia",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("content-Type", "application/json")

	resp, err := app.Test(req)

	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode((&user))

	if len(user.ID) == 0 {
		t.Errorf("Expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("Expecting the EncryptedPassword not to be included in the JSON response")
	}


	if user.FirstName != params.FirstName {
		t.Errorf("Expected username %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName{
		t.Errorf("Expected username %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email{
		t.Errorf("Expected username %s but got %s", params.Email, user.Email)
	}





	fmt.Println(user)

}

