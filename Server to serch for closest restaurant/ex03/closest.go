package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"bytes"
)

type Place struct {
	ID       string    `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

func main() {
	esURL := "http://localhost:9200"
	indexName := "places"

	http.HandleFunc("/api/recommend", func(w http.ResponseWriter, r *http.Request) {
		latParam := r.URL.Query().Get("lat")
		lonParam := r.URL.Query().Get("lon")

		lat, err := strconv.ParseFloat(latParam, 64)
		lon, err := strconv.ParseFloat(lonParam, 64)
		if err != nil {
			http.Error(w, "Invalid coordinates", http.StatusBadRequest)
			return
		}

		sortQuery := fmt.Sprintf(`{"sort":[{"_geo_distance":{"location":{"lat":%f,"lon":%f},"order":"asc","unit":"km","mode":"min","distance_type":"arc","ignore_unmapped":true}}]}`, lat, lon)

		url := fmt.Sprintf("%s/%s/_search?size=3", esURL, indexName)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(sortQuery)))
		if err != nil {
			log.Printf("Error making HTTP request: %v\n", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Returned code: %d\n", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response: %v\n", err)
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			log.Printf("Error parsing response: %v\n", err)
		}

		hits, _ := data["hits"].(map[string]interface{})["hits"].([]interface{})

		var places []Place
		for _, hit := range hits {
			source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
			if !ok {
				continue
			}

			var place Place
			place.ID,_ = source["id"].(string)
			place.Name, _ = source["name"].(string)
			place.Address, _ = source["address"].(string)
			place.Phone, _ = source["phone"].(string)

			location, ok := source["location"].(map[string]interface{})
			if ok {
				place.Location.Lat = location["lat"].(float64)
				place.Location.Lon = location["lon"].(float64)
			}

			places = append(places, place)
		}

		response := struct {
			Name   string  `json:"name"`
			Places []Place `json:"places"`
		}{
			Name:   "Recommendation",
			Places: places,
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding JSON: %v\n", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	})

	http.ListenAndServe(":8888", nil)
}
