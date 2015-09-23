/*
 Copyright 2015 Crunchy Data Solutions, Inc.
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package backend

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strconv"
)

type Car struct {
	ID    string
	Model string
	Price string
	Year  string
	Brand string
}

func GetCar(dbConn *sql.DB, ID string) (Car, error) {

	fmt.Println("GetCar called with ID=" + ID)
	car := Car{}

	queryStr := fmt.Sprintf("select id, model, price, year, brand from car where id = %s", ID)
	err := dbConn.QueryRow(queryStr).Scan(&car.ID, &car.Model,
		&car.Price, &car.Year, &car.Brand)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("GetCar: no car found with ID " + ID)
		return car, err
	case err != nil:
		fmt.Println("GetCar: error in get " + err.Error())
		return car, err
	}

	return car, err
}

func GetAllCars(dbConn *sql.DB) ([]Car, error) {
	fmt.Println("GetAllCars:called")
	var rows *sql.Rows
	cars := make([]Car, 0)
	var err error
	rows, err = dbConn.Query("select id, model, price, year, brand from car order by brand, model, year")
	if err != nil {
		fmt.Println("GetAllCars: error " + err.Error())
		return cars, err
	}
	defer rows.Close()
	for rows.Next() {
		car := Car{}
		if err = rows.Scan(&car.ID, &car.Model,
			&car.Price, &car.Year, &car.Brand); err != nil {
			fmt.Println("GetAllCars: error " + err.Error())
			return cars, err
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetAllCars: error " + err.Error())
		return cars, err
	}

	return cars, err
}

func AddCar(dbConn *sql.DB, car Car) (string, error) {

	fmt.Println("AddCar called")
	queryStr := fmt.Sprintf("insert into car ( model, price, year, brand) values ( '%s', '%s', '%s', '%s') returning id", car.Model, car.Price, car.Year, car.Brand)

	fmt.Println(queryStr)
	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println(err.Error())
		return "", err
	default:
		fmt.Println("inserted car ID " + strconv.Itoa(carid))
	}

	return strconv.Itoa(carid), err
}

func UpdateCar(dbConn *sql.DB, car Car) error {

	fmt.Println("UpdateServer:called")
	queryStr := fmt.Sprintf("update car set ( model, price, year, brand) = ('%s', '%s', '%s', '%s') where id = %s returning id", car.Model, car.Price, car.Year, car.Brand, car.ID)

	fmt.Println(queryStr)
	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println("UpdateCar: error in query" + err.Error())
		return err
	default:
		fmt.Println("UpdateCar:car updated " + strconv.Itoa(carid))
	}

	return err
}

func DeleteCar(dbConn *sql.DB, id string) error {

	queryStr := fmt.Sprintf("delete from car where  id=%s returning id", id)
	fmt.Println(queryStr)

	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println("DeleteCar: error in query" + err.Error())
		return err
	default:
		fmt.Println("DeleteCar:car deleted " + strconv.Itoa(carid))
	}
	return err
}
