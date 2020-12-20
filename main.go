package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type User struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name,omitempty" bson:"name,omitempty"`
	DateOfBirth       string             `json:"dateofbirth,omitempty" bson:"dateofbirth,omitempty"`
	PhoneNo           int                `json:"phoneno,omitempty" bson:"phoneno,omitempty"`
	EmailAddress      string             `json:"emailaddress,omitempty" bson:"emailaddress,omitempty"`
	CreationTimeStamp int                `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
}
type Contact struct {
	UserIdOne string `json:"UserIdOne,omitempty" bson:"UserIdOne,omitempty"`
	UserIdTwo string `json:"UserIdTwo,omitempty" bson:"UserIdTwo,omitempty"`
	Timestamp int    `json:"Timestamp ,omitempty" bson:"Timestamp ,omitempty"`
}

// func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {}
// func GetUserEndpoint(response http.ResponseWriter, request *http.Request)    {}
// func GetContactEndpoint(response http.ResponseWriter, request *http.Request) {}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	router := mux.NewRouter()
	router.HandleFunc("/user", CreateUserEndpoint).Methods("POST")
	router.HandleFunc("/user", GetAUserEndpoint).Methods("GET")
	router.HandleFunc("/user/{id}", GetUserEndpoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(request.Body).Decode(&user)
	collection := client.Database("thepolyglotdeveloper").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var user User
	collection := client.Database("thepolyglotdeveloper").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

func GetAUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var user []User
	collection := client.Database("thepolyglotdeveloper").Collection("user")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		user = append(user, user)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
