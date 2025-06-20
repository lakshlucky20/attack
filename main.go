package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"https://maps.googleapis.com/maps/api/geocode/json?address=Paris&key=AIzaSyDsBDwflgnw7DPOCi5DWiAezkAptmwPSLA",
	"https://maps.googleapis.com/maps/api/place/textsearch/json?query=restaurants+in+London&key=AIzaSyDsBDwflgnw7DPOCi5DWiAezkAptmwPSLA",
	"https://maps.googleapis.com/maps/api/directions/json?origin=Toronto&destination=Montreal&key=AIzaSyDsBDwflgnw7DPOCi5DWiAezkAptmwPSLA",
	"https://maps.googleapis.com/maps/api/distancematrix/json?origins=Seattle&destinations=San+Francisco&key=AIzaSyDsBDwflgnw7DPOCi5DWiAezkAptmwPSLA",
}

func sendRequest(url string, id int, wg *sync.WaitGroup, successCounter *int, lock *sync.Mutex) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[FAILED] Req #%d to %s\n", id, url)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		lock.Lock()
		*successCounter++
		lock.Unlock()
		log.Printf("[OK] %s [Req #%d]\n", url, id)
	} else {
		log.Printf("[ERR %d] %s [Req #%d]\n", resp.StatusCode, url, id)
	}
	_ = body // You can parse JSON if needed
}

func attackHandler(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	var lock sync.Mutex
	successCounter := 0

	start := time.Now()

	for _, url := range urls {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go sendRequest(url, i+1, &wg, &successCounter, &lock)
		}
	}

	wg.Wait()
	duration := time.Since(start)
	fmt.Fprintf(w, "All 400 requests finished in %s. Successful: %d\n", duration, successCounter)
}

func main() {
	http.HandleFunc("/attack", attackHandler)
	http.Handle("/", http.FileServer(http.Dir("./static"))) // serve index.html from static/

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
