package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Place struct {
	// ID      string  
	Name    string
	Address string
	Phone   string
}

func main() {
	esURL := "http://localhost:9200"
	indexName := "places"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pageParam := r.URL.Query().Get("page")
		page, err := strconv.Atoi(pageParam)
		if err != nil || page < 1 {
			http.Error(w, "Invalid 'page' value: "+pageParam, http.StatusBadRequest)
			return
		}

		limit := 10
		offset := (page - 1) * limit

		url := fmt.Sprintf("%s/%s/_search?size=%d&from=%d", esURL, indexName, limit, offset)
		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		var data map[string]interface{}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		hits, _ := data["hits"].(map[string]interface{})["hits"].([]interface{})

		var places []Place
		for _, hit := range hits {
			source, _ := hit.(map[string]interface{})["_source"].(map[string]interface{})

			var place Place
			// place.ID,_ = source["id"].(string)
			place.Name, _ = source["name"].(string)
			place.Address, _ = source["address"].(string)
			place.Phone, _ = source["phone"].(string)

			places = append(places, place)
		}

		totalHits, _ := data["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)

		lastPage := int(totalHits+float64(limit-1)) / limit

		previousPage := max(1, page-1)
		nextPage := min(lastPage, page+1)

		fmt.Fprintf(w, `<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Places</title>
    <meta name="description" content="">
    <meta name="viewport" content="width=device-width, initial-scale=1">
</head>

<body>
<h5>Total: %d</h5>
<ul>
`, int(totalHits))

		for _, place := range places {
			fmt.Fprintf(w, `
    <li>
        <div>%s</div>
		<div>%s</div>
        <div>%s</div>
    </li>`, place.Name, place.Address, place.Phone)
		}
		fmt.Fprintf(w, `</ul>
<a href="/?page=%d">Previous</a>
<a href="/?page=%d">Next</a>
<a href="/?page=%d">Last</a>
</body>
</html>`, previousPage, nextPage, lastPage)
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
