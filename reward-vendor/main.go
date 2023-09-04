package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type Reward struct {
	RewardId  string `json:"rewardId"`
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

var logger *zap.Logger
var rewards []Reward

func CreateReward(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var reward Reward
	_ = json.NewDecoder(r.Body).Decode(&reward)

	logger.Info("creating a reward", zap.Any("reward", reward))

	rewards = append(rewards, reward)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(reward)
}

func main() {
	defer logger.Sync() // Ensure all buffered logs are written

	logger.Info("Starting the reward vender...")

	r := mux.NewRouter()

	// Sample Data
	rewards = append(rewards, Reward{RewardId: "RWD0000", UserId: "U0000", FirstName: "John", LastName: "Doe", Email: "john@example.com"})
	r.HandleFunc("/rewards", CreateReward).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}
