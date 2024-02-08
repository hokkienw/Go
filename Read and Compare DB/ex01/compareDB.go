package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

type Ingredients struct {
	Name  string `json:"ingredient_name" xml:"itemname"`
	Count string `json:"ingredient_count" xml:"itemcount"`
	Unit  string `json:"ingredient_unit" xml:"itemunit"`
}

type Cake struct {
	Name        string       `json:"name" xml:"name"`
	Time        string       `json:"time" xml:"stovetime"`
	Ingredients []Ingredients `json:"ingredients" xml:"ingredients>item"`
}

type Database struct {
	Cakes []Cake `json:"cake" xml:"cake"`
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

func CompareDatabases(oldDB, newDB Database) {
	for _, newCake := range newDB.Cakes {
		cakeFound := false

		for _, oldCake := range oldDB.Cakes {
			if newCake.Name == oldCake.Name {
				cakeFound = true

				if newCake.Time != oldCake.Time {
					fmt.Printf("CHANGED cooking time for cake \"%s\" - \"%s\" instead of \"%s\"\n", newCake.Name, newCake.Time, oldCake.Time)
				}

				oldIngredientsMap := make(map[string]Ingredients)
				for _, ingredient := range oldCake.Ingredients {
					oldIngredientsMap[ingredient.Name] = ingredient
				}

				for _, newIngredient := range newCake.Ingredients {
					oldIngredient, exists := oldIngredientsMap[newIngredient.Name]
					if !exists {
						fmt.Printf("ADDED ingredient \"%s\" for cake \"%s\"\n", newIngredient.Name, newCake.Name)
						continue
					}

					if newIngredient.Unit != oldIngredient.Unit {
						fmt.Printf("CHANGED unit for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", newIngredient.Name, newCake.Name, newIngredient.Unit, oldIngredient.Unit)
					}

					if newIngredient.Count != oldIngredient.Count {
						fmt.Printf("CHANGED unit count for ingredient \"%s\" for cake \"%s\" - \"%s\" instead of \"%s\"\n", newIngredient.Name, newCake.Name, newIngredient.Count, oldIngredient.Count)
					}

					delete(oldIngredientsMap, newIngredient.Name)
				}

				for _, removedIngredient := range oldIngredientsMap {
					fmt.Printf("REMOVED ingredient \"%s\" for cake \"%s\"\n", removedIngredient.Name, newCake.Name)
				}

				break
			}
		}

		if !cakeFound {
			fmt.Printf("ADDED cake \"%s\"\n", newCake.Name)
		}
	}

	for _, oldCake := range oldDB.Cakes {
		cakeFound := false

		for _, newCake := range newDB.Cakes {
			if oldCake.Name == newCake.Name {
				cakeFound = true
				break
			}
		}

		if !cakeFound {
			fmt.Printf("REMOVED cake \"%s\"\n", oldCake.Name)
		}
	}
}


func main() {
	oldDBFile := flag.String("old", "", "Old database XML file")
	newDBFile := flag.String("new", "", "New database JSON file")
	flag.Parse()

	if *oldDBFile == "" || *newDBFile == "" {
		fmt.Println("Print both old and new database files using -old and -new")
		os.Exit(1)
	}

	oldDB, err := ReadXMLData(*oldDBFile)
	if err != nil {
		fmt.Println("Error reading old database:", err)
		os.Exit(1)
	}

	newDB, err := ReadJsonData(*newDBFile)
	if err != nil {
		fmt.Println("Error reading new database:", err)
		os.Exit(1)
	}

	CompareDatabases(oldDB, newDB)
}
