package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
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

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        pageParam := r.URL.Query().Get("page")
        page, err := strconv.Atoi(pageParam)
        if err != nil || page < 1 {
            http.Error(w, `{"error": "Invalid 'page' value: `+pageParam+`"}`, http.StatusBadRequest)
            return
        }

        limit := 10
        offset := (page - 1) * limit

        url := fmt.Sprintf("%s/%s/_search?size=%d&from=%d", esURL, indexName, limit, offset)
        resp, _ := http.Get(url)
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
            return
        }

        body, _ := ioutil.ReadAll(resp.Body)
        var data map[string]interface{}
        if err := json.Unmarshal(body, &data); err != nil {
            http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
            return
        }

        hits, ok := data["hits"].(map[string]interface{})["hits"].([]interface{})
        if !ok {
            http.Error(w, `{"error": "Invalid format"}`, http.StatusInternalServerError)
            return
        }

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
                place.Location.Lat, _ = location["lat"].(float64)
                place.Location.Lon, _ = location["lon"].(float64)
            }

            places = append(places, place)
        }

        totalHits, ok := data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
        if !ok {
            http.Error(w, `{"error": "Invalid format"}`, http.StatusInternalServerError)
            return
        }

        lastPage := int(totalHits+float64(limit-1)) / limit

        previousPage := max(1, page-1)
        nextPage := min(lastPage, page+1)

        response := struct {
            Name      string   `json:"name"`
            Total     int      `json:"total"`
            Places    []Place  `json:"places"`
            PrevPage  int      `json:"prev_page"`
            NextPage  int      `json:"next_page"`
            LastPage  int      `json:"last_page"`
        }{
            Name:     "Places",
            Total:    int(totalHits),
            Places:   places,
            PrevPage: previousPage,
            NextPage: nextPage,
            LastPage: lastPage,
        }

        jsonResponse, err := json.Marshal(response)
        if err != nil {
            http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write(jsonResponse)
    })

    http.ListenAndServe(":8888", nil)
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
