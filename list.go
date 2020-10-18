package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

//Function to find the meeting details using a participants email-id

func FindEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var meeting []Meeting
	Mydatabase = client.Database("MeetingScheduler").Collection("meet_Schedule")
	parameter := mux.Vars(request)
	email, _ := request.Get(parameter["Email"])
	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	curs, err := Mydatabase.Find(ctx, Meeting{Email: email}).Decode(&meeting)
	if err!=nil {
		log.Fatal(err)
	}
	var filtered []bson.M
	if err = filterCursor.All(ctx, &filtered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(filtered)
}

//Function to find meetings which fall between the start and the End Time

func FindRangeEndPoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var meeting []Meeting
	Mydatabase = client.Database("MeetingScheduler").Collection("meet_Schedule")
	parameter := mux.Vars(request)
	btime, ftime := request.Get(parameter["stime","etime"])
	ctx, _ := context.WithTimeout(context.Backgroung(), 10*time.Second)
	curs, err := Mydatabase.Find(ctx, Meeting{stime: btime, etime:ftime}).Decode(&meeting)
	if err!=nil {
		log.Fatal(err)
	}
	var rangefiltered []bson.M
	if err = filterCursor.All(ctx, &rangefiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(rangefiltered)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	Mydatabase = client.Database("MeetingScheduler")
	meetingcollection = Mydatabase.Collection("meet_Schedule")

	cursor, err := meetingcollection.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	defer cursor.Close()

	for cursor.Next(ctx) {
		var people []bson.M
		if err = cursor.All(ctx, &people); err != nil {
			log.Fatal(err)
		}
		fmt.Println(people)
	}

	// To get the list of meetings with a particular start time

	filterCursor, err := meetingcollection.Find(ctx, bson.M{"stime": "11:00:00"})
	if err != nil {
		log.Fatal(err)
	}
	var meetingsfiltered []bson.M
	if err = filterCursor.All(ctx, &meetingsfiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(meetingsfiltered)

	//To list all meetings of a particular participant

	filterCursor, err := meetingcollection.Find(ctx, bson.M{"Name": "Anastesia", "ID": "1234"})
	if err != nil {
		log.Fatal(err)
	}
	var meetingsfiltered []bson.M
	if err = filterCursor.All(ctx, &meetingsfiltered); err != nil {
		log.Fatal(err)
	}
	fmt.Println(meetingsfiltered)

	// using GET method
	router := mux.NewRouter()
	router.HandleFunc("/meeting?participant=<{email}>", FindEndPoint).Methods("GET")
	router.HandleFunc("/meeting?start=<{btime}>&end=<{ftime}>", FindEndPoint).Methods("GET")
}
