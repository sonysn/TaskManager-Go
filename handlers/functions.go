package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	mgm.IDField `bson:",inline"`
	Title       string    `bson:"title"`
	Description string    `bson:"description"`
	Due_Date    time.Time `bson:"due_date"`
	Completed   bool      `bson:"completed"`
}

const redisTaskListKey = "TaskList"

func ListAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//Create an empty array of tasks from the model
	var taskList []Task

	//!IMPLEMENTING REDIS CACHE ON FETCH
	results, errd := checkRedisCache()
	if errd != nil {
		//!IF NOT FOUND IN REDIS CACHE, FETCH FROM MONGODB

		//Get a reference to the Task collection
		taskCollection := mgm.Coll(&Task{})

		// Find all documents in the collection
		cursor, err := taskCollection.Find(context.Background(), bson.M{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// log.Println(err)
			return
		}
		defer cursor.Close(context.TODO()) // Close the cursor when done

		// Iterate through the cursor and decode documents into the Task struct
		for cursor.Next(context.TODO()) {
			var task Task
			if err := cursor.Decode(&task); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				// log.Println(err)
				return
			}
			taskList = append(taskList, task)
		}

		// Check for errors during iteration
		if err := cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// log.Println(err)
			return
		}

		//!SAVE DATA TO REDIS CACHE AS A HASH
		for index, task := range taskList {
			taskNumber := fmt.Sprintf("Task %d", index)
			RedisClient.HSet(context.Background(), redisTaskListKey, taskNumber, task)
		}

		//!DELETE LIST AFTER 1 HOUR
		RedisClient.Expire(context.Background(), redisTaskListKey, 1*time.Hour)

		// Encode the tasks into JSON and send the response
		jsonData, err := json.Marshal(taskList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			// log.Println(err)
			return
		}

		w.Write(jsonData)
		w.WriteHeader(http.StatusOK)
		return
	} else {
		//!ELSE RETURN REDIS RESULTS
		json.NewEncoder(w).Encode(results)
		w.WriteHeader(http.StatusOK)
	}
}

func RetrieveATaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	taskID := r.URL.Query().Get("taskID")

	// Convert the taskID to mongodb ObjectID
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		// log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get a reference to the Task collection
	taskCollection := mgm.Coll(&Task{})

	// Find a document in the collection by ID
	var document Task
	cursor, err := taskCollection.Find(context.Background(), bson.M{"_id": objectID})

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := cursor.Decode(&document); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the task into JSON and send the response
	jsonData, err := json.Marshal(document)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		// log.Println(err)
		return
	}

	w.Write(jsonData)
	w.WriteHeader(http.StatusOK)
}

func CreateATask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var reqBody struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Due_Date    time.Time `json:"due_date"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Get a reference to the Task collection
	taskCollection := mgm.Coll(&Task{})

	//Create a new task
	err = taskCollection.Create(&Task{
		Title:       reqBody.Title,
		Description: reqBody.Description,
		Due_Date:    reqBody.Due_Date,
		Completed:   false,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateATask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	taskID := r.URL.Query().Get("taskID")

	var reqBody struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Due_Date    time.Time `json:"due_date"`
		Completed   bool      `json:"completed"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert the taskID to mongodb ObjectID
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		// log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get a reference to the Task collection
	taskCollection := mgm.Coll(&Task{})

	res := taskCollection.FindOneAndUpdate(context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"title": reqBody.Title, "description": reqBody.Description, "due_date": reqBody.Due_Date, "completed": reqBody.Completed}},
	)

	if res.Err() != nil {
		http.Error(w, res.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteATask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	taskID := r.URL.Query().Get("taskID")

	// Convert the taskID to mongodb ObjectID
	objectID, err := primitive.ObjectIDFromHex(taskID)
	if err != nil {
		// log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get a reference to the Task collection
	taskCollection := mgm.Coll(&Task{})

	res := taskCollection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID})

	if res.Err() != nil {
		http.Error(w, res.Err().Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// This function returns a the list or an error if not found
func checkRedisCache() (map[string]string, error) {
	results, err := RedisClient.HGetAll(context.Background(), redisTaskListKey).Result()
	if err == redis.Nil {
		return nil, err
	}
	return results, nil
}

// Custom Management Command To marks tasks as completed if their due date has passed
func UpdateTaskWorker() {
	//Continue to loop forever
	for {
		//Create an empty array of tasks from the model
		var taskList []Task

		//Get a reference to the Task collection
		taskCollection := mgm.Coll(&Task{})

		// Find all documents in the collection
		cursor, err := taskCollection.Find(context.Background(), bson.M{})
		if err != nil {
			// log.Println(err)
			return
		}
		defer cursor.Close(context.Background()) // Close the cursor when done

		// Iterate through the cursor and decode documents into the Task struct
		for cursor.Next(context.Background()) {
			var task Task
			if err := cursor.Decode(&task); err != nil {
				// log.Println(err)
				return
			}
			taskList = append(taskList, task)
		}

		// Check for errors during iteration
		if err := cursor.Err(); err != nil {
			// log.Println(err)
			return
		}

		for _, task := range taskList {
			//Set Completed to true if the due date has passed
			if hasPassedDueDate(task.Due_Date) {
				// updatedTask := Task{
				// 	Title:       task.Title,
				// 	Description: task.Description,
				// 	Due_Date:    task.Due_Date,
				// 	Completed:   true,
				// }
				update := bson.M{"$set": bson.M{"completed": true}}
				filter := bson.M{"_id": task.ID}

				taskCollection.FindOneAndUpdate(context.Background(), filter, update)
			}

		}
	}
}

func hasPassedDueDate(targetDate time.Time) bool {
	//Get Current Date
	currentDate := time.Now().UTC()
	return currentDate.After(targetDate)
}
