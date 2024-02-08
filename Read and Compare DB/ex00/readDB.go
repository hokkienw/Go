package main

import (
	"os"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"flag"
)

type Ingredients struct {
	Name string `json:"ingredient_name" xml:"itemname"`
	Count string `json:"ingredient_count" xml:"itemcount"`
	Unit string `json:"ingredient_unit" xml:"itemunit"`
}

type Cake struct {
	Name string `json:"name" xml:"name"`
	Time string `json:"time" xml:"stovetime"`
	Ingredients []Ingredients `json:"ingredients" xml:"ingredients>item"`
}

type Database struct {
    Cakes []Cake `json:"cake" xml:"cake"`
}

func ReadJsonData(filename string) (Database, error) {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return Database{}, err
    }

    var db Database
    if err := json.Unmarshal(file, &db); err != nil {
        return Database{}, err
    }

    return db, nil
}

func ReadXMLData(filename string) (Database, error) {
    file, err := ioutil.ReadFile(filename)
    if err != nil {
        return Database{}, err
    }

    var db Database
    if err := xml.Unmarshal(file, &db); err != nil {
        return Database{}, err
    }

    return db, nil
}

func main() {

	filename := flag.String("f", "", "Database file to read")
    flag.Parse()

	if *filename == "" {
        fmt.Println("Print filename using -f")
        os.Exit(1)
    }

    if strings.HasSuffix(*filename, ".json") {
        db, err := ReadJsonData(*filename)
        if err != nil {
            fmt.Println("Error JSON", err)
            os.Exit(1)
        }

        fmt.Printf("%+v\n", db)
    } else if strings.HasSuffix(*filename, ".xml") {
        db, err := ReadXMLData(*filename)
        if err != nil {
            fmt.Println("Error XML", err)
            os.Exit(1)
        }

        fmt.Printf("%+v\n", db)
    } else {
        fmt.Println("Error file")
        os.Exit(1)
    }
}