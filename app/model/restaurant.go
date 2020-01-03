package model

import (
	"../config"
	"encoding/json"
	"io"
	"log"
)

// Restaurant schema for db
type Restaurant struct {
	ID   	 	int    	`json:"rest_id"`
	Name 	 	string 	`json:"name"`
	Location 	string 	`json:"location"`
}

// SetData : post stream into db
func (res *Restaurant) SetData(stream io.Reader) FoulError {
	decoder := json.NewDecoder(stream)
	err := decoder.Decode(&res)
	if err != nil {
		panic(err)
	}
	log.Println(res.Name)

	errF := res.insertIntoDB()
	if errF != nil {
		return errF
	}
	return nil
}

func (res Restaurant) insertIntoDB() FoulError {
	_, err := config.DB.Exec("INSERT INTO RESTAURANT(name, location) VALUES (?,?)", res.Name, res.Location)
	if err != nil {
		return &Foul{"DB.EXEC ERROR", err.Error(), TraceCall()}
	}
	return nil
}

// GetData : call getAllRestaurants
func (res *Restaurant) GetData() ([]*Restaurant, FoulError) {

	data, err := getAllRestaurants()
	return data, err

}

// getAllRestaurants : return map with restaurants
func getAllRestaurants() ([]*Restaurant, FoulError) {
	var m []*Restaurant
	rows, err := config.DB.Query("SELECT * FROM RESTAURANT")
	if err != nil {
		errF := Foul{"DB.QUERY ERROR", err.Error(), TraceCall()}
		return m, &errF
	}
	for rows.Next() {
		var res Restaurant

		err := rows.Scan(&res.ID, &res.Name, &res.Location)
		if err != nil {
			return m, &Foul{"DB.SCAN ERROR", err.Error(), TraceCall()}
		}

		m = append(m, &res)
	}

	return m, nil
}
