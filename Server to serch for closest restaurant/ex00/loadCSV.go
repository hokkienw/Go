 Elpackage main

import (
    "os"
    "fmt"
    "encoding/csv"
    "log"
    "strconv"
    "io"
    "encoding/json"
    "net/http"
    "bytes"
    "strings"
)

type Restaurant struct {
    ID       string  `json:"id"`
    Name     string  `json:"name"`
    Address  string  `json:"address"`
    Phone    string  `json:"phone"`
    Location struct {
        Lat float64 `json:"lat"`
        Lon float64 `json:"lon"`
    } `json:"location"`
}

func main() {
    args := os.Args
    if len(args) != 2 {
        fmt.Println("Usage: ./loadCSV /path/to/file.csv")
        return
    }

    path := args[1]
    file, err := os.Open(path)
    if err != nil {
        log.Printf("File reading error: %v", err)
        return
    }

    defer file.Close()

    reader := csv.NewReader(file)
    reader.Comma = '\t'

    _, err = reader.Read()

    createIndexIfNotExists()
    id := 1
    for {
        row, err := reader.Read()
        if err != nil {
            if err == io.EOF {
                break
            }
            log.Printf("File reading error: %v", err)
            continue
        }

        restaurant := Restaurant{}
        restaurant.ID = row[0]
        restaurant.Name = row[1]
        restaurant.Address = row[2]
        restaurant.Phone = row[3]
        restaurant.Location.Lon, _ = strconv.ParseFloat(row[4], 64)
        restaurant.Location.Lat, _ = strconv.ParseFloat(row[5], 64)

		restaurantJSON, err := json.Marshal(restaurant)
		if err != nil {
			log.Printf("Error JSON: %v", err)
			return
		}
		err = sendDataToElastic(restaurantJSON, id)
        id++
    }
}

func createIndexIfNotExists() {
    url := "http://localhost:9200/places"
    mapping := `
    {
      "mappings": {
        "properties": {
          "name": {
            "type": "text"
          },
          "address": {
            "type": "text"
          },
          "phone": {
            "type": "text"
          },
          "location": {
            "type": "geo_point"
          }
        }
      }
    }`

    req, err := http.NewRequest("PUT", url, strings.NewReader(mapping))
    if err != nil {
        log.Printf("Index creating error: %v", err)
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Index creating error: %v", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
        log.Printf("Index creation failed: %d", resp.StatusCode)
    }
}

func sendDataToElastic(data []byte, id int) error {
    url := fmt.Sprintf("http://localhost:9200/places/_doc/%d", id)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/x-ndjson")
    client := &http.Client{}
    resp, err := client.Do(req)

    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Request fail code: %d", resp.StatusCode)
    }
    return nil
}
