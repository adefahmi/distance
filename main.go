package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

type DistanceResponse struct {
	Meter float64 `json:"meter"`
	Km     string  `json:"km"`
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert decimal degrees to radians
	lat1, lon1, lat2, lon2 = degToRad(lat1), degToRad(lon1), degToRad(lat2), degToRad(lon2)

	// Radius of the Earth in kilometers
	radius := 6371.0

	// Haversine formula
	dlat := lat2 - lat1
	dlon := lon2 - lon1
	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlon/2)*math.Sin(dlon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := radius * c * 1000.0 // Multiply by 1000 to convert to meters

	return distance
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

func distanceHandler(w http.ResponseWriter, r *http.Request) {
	// Extract latitude and longitude from query parameters
	lat1Str := r.URL.Query().Get("lat1")
	lon1Str := r.URL.Query().Get("lon1")
	lat2Str := r.URL.Query().Get("lat2")
	lon2Str := r.URL.Query().Get("lon2")

	// Convert latitude and longitude to float64
	lat1, err := strconv.ParseFloat(lat1Str, 64)
	if err != nil {
		http.Error(w, "Invalid latitude for point 1", http.StatusBadRequest)
		return
	}

	lon1, err := strconv.ParseFloat(lon1Str, 64)
	if err != nil {
		http.Error(w, "Invalid longitude for point 1", http.StatusBadRequest)
		return
	}

	lat2, err := strconv.ParseFloat(lat2Str, 64)
	if err != nil {
		http.Error(w, "Invalid latitude for point 2", http.StatusBadRequest)
		return
	}

	lon2, err := strconv.ParseFloat(lon2Str, 64)
	if err != nil {
		http.Error(w, "Invalid longitude for point 2", http.StatusBadRequest)
		return
	}

	// Calculate distance
	distance := calculateDistance(lat1, lon1, lat2, lon2)

	// Create response object
	response := DistanceResponse{
		Meter: distance,
		Km:    fmt.Sprintf("%.2f", distance/1000.0),
	}

	// Convert response object to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write JSON response
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/distance", distanceHandler)

	fmt.Println("Starting server on http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", nil))
}
