package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"std/vendor/go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Meeting struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title        string             `json:"title,omitmepty" bson:title,omitempty"`
	Participants int64              `json:"participants,omitempty" bson:"participants,omitempty"`
	stime        time.Time          `json:"btime,omitempty" bson:"btime,omitempty"`
	etime        time.Time          `json:"ftime,omitempty" bson:"ftime,omitempty"`
	Name         string             `json:"name,omitempty" bson:"name,omitempty"`
	Email        string             `json:"_email,omitempty" bson:"_email,omitempty"`
	RSVP         string             `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}

var client *mongo.Client

func CreateEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var meeting Meeting
	json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("MeetinScheduler").Collection("meet_Schedule")
	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}

func GetEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var meet []Meeting
	collection := client.Database("MeetingScheduler").Collection("meet_Schedule")
	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var meeting Meeting
		cursor.Decode(&meeting)
		people = append(meet, meeting)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalSErverError)
		response.Write([]byte(`"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(meet) //Secure threading
}

func GetSingleEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	param := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(param["ID"])
	collection := client.Database("MeetingScheduler").Collection("meet_Schedule")
	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalSErverError)
		response.Write([]byte(`"message":"` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
}

func main() {
	fmt.Println("Starting the application...")

	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	client, _ = mongo.Connect(ctx, "mongodb://localhost:8080")
	router := mux.NewRouter()
	router.HandleFunc("/meeting", CreateEndPoint).Methods("POST")
	router.HandleFunc("/meet", GetEndPoint).Methods("GET")
	router.HandleFunc("/meeting/{id}", GetSingleEndPoint).Methods("GET")
	http.ListenAndServe(":12345", router)
}
