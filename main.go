package main

import (
	"TaskManager/handlers"
	"log"
	"net/http"
)

func main() {
	startService()
}

func startService() {
	connectToDatabases()
	port := "9000"

	//List all Tasks endpoint
	http.HandleFunc("/list", handlers.ListAllTasks)

	//Retrieve a single task by ID endpoint
	http.HandleFunc("/retrieve", handlers.RetrieveATaskByID)

	//Create a new task endpoint
	http.HandleFunc("/create", handlers.CreateATask)

	//Update an existing task endpoint
	http.HandleFunc("/update", handlers.UpdateATask)

	//Delete a task endpoint
	http.HandleFunc("/delete", handlers.DeleteATask)

	log.Println("Server listening on port", port)

	//Start Worker Service
	go handlers.UpdateTaskWorker()

	http.ListenAndServe(":"+port, nil)
}

func connectToDatabases() {
	handlers.ConnectToRedis()
	handlers.ConnectToMongo()
}
