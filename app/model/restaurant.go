package model

import (
	"../config"
	"encoding/json"
	"io"
	"log"
)

// Restaurant schema for db
type Restaurant struct {
	ID   int    `json:"rest_id"`
	Name string `json:"name"`
	Location string `json:"location"`
}

// SetData : post stream into db
func (res *Restaurant) SetData(stream io.Reader) error {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&res)
	if err != nil {
		panic(err)
	}
	log.Println(res.Name)

	err = res.insertIntoDB()
	if err != nil {
		return err
	}
	return nil
}

func (res *Restaurant) insertIntoDB() error {
	_, err := config.DB.Exec("INSERT INTO RESTAURANT(name, location) VALUES (?,?)", res.Name, res.Location)
	if err != nil {
		return err
	}
	return nil
}

// GetData : call getAllRestaurants
func (res *Restaurant) GetData() ([]*Restaurant, error) {

	data, err := getAllRestaurants()
	if err != nil {
		log.Fatal("getAllRestaurants error: ", err)
	}
	return data, err

}

// getAllRestaurants : return map with restaurants
func getAllRestaurants() ([]*Restaurant, error) {
	var m []*Restaurant
	rows, err := config.DB.Query("SELECT * FROM RESTAURANT")
	if err != nil {
		return m, err
	}
	for rows.Next() {
		var res Restaurant

		err := rows.Scan(&res.ID, &res.Name, &res.Location)
		if err != nil {
			return m, err
		}

		m = append(m, &res)
	}

	return m, nil
}
