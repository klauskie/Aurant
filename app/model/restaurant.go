package model

import (
	"encoding/json"
	"io"
	"log"
	"strconv"

	"app/config"
)

// Restaurant schema for db
type Restaurant struct {
	ID   int    `json:"rest_id"`
	Name string `json:"name"`
}

// SetData : post stream into db
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

func insertIntoDB(resName string) error {
	_, err := config.DB.Exec("INSERT INTO restaurant(name) VALUES (?)", resName)
	if err != nil {
		return err
	}
	return nil
}

// GetData : call getAllRestaurants
func (res *Restaurant) GetData() ([]byte, error) {

	data, err := getAllRestaurants()
	if err != nil {
		log.Fatal("getAllRestaurants error: ", err)
	}
	output, err2 := json.Marshal(data)
	if err2 != nil {
		log.Fatal("Encoding error: ", err2)
	}
	return output, err2

}

// getAllRestaurants : return map with restaurants
func getAllRestaurants() (map[string]Restaurant, error) {
	m := make(map[string]Restaurant)
	rows, err := config.DB.Query("SELECT * FROM restaurant")
	if err != nil {
		return m, err
	}
	for rows.Next() {
		var res Restaurant

		err := rows.Scan(&res.ID, &res.Name)
		if err != nil {
			return m, err
		}

		strID := strconv.Itoa(res.ID)
		m[strID] = res
	}

	return m, nil
}
