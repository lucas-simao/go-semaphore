package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

type Response struct {
	TotalTime  float64
	Goroutines int
	Text       string
}

var (
	semaphoreLimit = 5 // limit of goroutines that can be created
	semaphore      = make(chan int, semaphoreLimit)
)

func main() {
	http.HandleFunc("/with-semaphore", WithSemaphore)
	http.HandleFunc("/without-semaphore", WithoutSemaphore)

	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Panic(err)
	}
}

func WithSemaphore(w http.ResponseWriter, r *http.Request) {
	var goroutinesBeforeProcess = runtime.NumGoroutine()

	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil {
		fmt.Fprintf(w, "error to convert string to int: %v", err)
		return
	}

	start := time.Now()
	var goroutinesCount int

	for i := 0; i < quantity; i++ {
		semaphore <- 1
		go func(i int) {
			Process(i)
			<-semaphore
		}(i)
		if runtime.NumGoroutine() > goroutinesCount {
			goroutinesCount = runtime.NumGoroutine()
		}
	}

	GenerateResponse(w, start, goroutinesCount-goroutinesBeforeProcess)
}

func WithoutSemaphore(w http.ResponseWriter, r *http.Request) {
	var goroutinesBeforeProcess = runtime.NumGoroutine()

	quantity, err := strconv.Atoi(r.URL.Query().Get("quantity"))
	if err != nil {
		fmt.Fprintf(w, "error to convert string to int: %v", err)
		return
	}

	start := time.Now()
	var goroutinesCount int

	for i := 0; i < quantity; i++ {
		go func(i int) {
			Process(i)
		}(i)
		if runtime.NumGoroutine() > goroutinesCount {
			goroutinesCount = runtime.NumGoroutine()
		}
	}

	GenerateResponse(w, start, goroutinesCount-goroutinesBeforeProcess)
}

func Process(i int) {
	fmt.Printf("quantity: %d\n", i)
	time.Sleep(1 * time.Second)
}

func GenerateResponse(w http.ResponseWriter, start time.Time, goroutines int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(Response{
		TotalTime:  time.Since(start).Seconds(),
		Goroutines: goroutines,
		Text:       fmt.Sprintf("total time execution %.0fs and created %d goroutines", time.Since(start).Seconds(), goroutines),
	})
}
