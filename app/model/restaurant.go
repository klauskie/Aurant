package model

import (
	"encoding/json"
	"io"
	"log"

	"app/config"
)

// Restaurant schema for db
type Restaurant struct {
	ID   int    `json:"rest_id"`
	Name string `json:"name"`
}

func (res *Restaurant) SetData(stream io.Reader) {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&res)
	if err != nil {
		panic(err)
	}
	log.Println(res.Name)

	err = insertIntoDB(res.Name)
	if err != nil {
		panic(err)
	}
}

func (res *Restaurant) GetData(resID int) {

	selectFromDB(resID)
}

func insertIntoDB(resName string) error {
	_, err := config.DB.Exec("INSERT INTO restaurant(name) VALUES (?)", resName)
	if err != nil {
		return err
	}
	return nil
}

func selectFromDB(resID int) {
	// Execute the query
	results, err := config.DB.Query("SELECT * FROM restaurant WHERE rest_id = ?", resID)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var res Restaurant
		// for each row, scan the result into our tag composite object
		err = results.Scan(&res.ID, &res.Name)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		log.Printf(res.Name)
	}
}
