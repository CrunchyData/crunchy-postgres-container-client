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

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-pg-client/backend"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func init() {
	fmt.Println("before parsing in init")
	flag.Parse()

}

var CPMDIR = "/var/cpm/"
var CPMBIN = CPMDIR + "bin/"

func main() {

	fmt.Println("at top of backend main")

	var err error

	handler := rest.ResourceHandler{
		PreRoutingMiddlewares: []rest.Middleware{
			&rest.CorsMiddleware{
				RejectNonCorsRequests: false,
				OriginValidator: func(origin string, request *rest.Request) bool {
					return true
				},
				AllowedMethods: []string{"DELETE", "GET", "POST", "PUT"},
				AllowedHeaders: []string{
					"Accept", "Content-Type", "X-Custom-Header", "Origin"},
				AccessControlAllowCredentials: true,
				AccessControlMaxAge:           3600,
			},
		},
		EnableRelaxedContentType: true,
	}

	err = handler.SetRoutes(
		&rest.Route{"POST", "/car/add", backend.AddCar},
		&rest.Route{"POST", "/car/update", backend.UpdateCar},
		&rest.Route{"GET", "/car/list", backend.GetAllCars},
		&rest.Route{"GET", "/car/:ID.", backend.GetCar},
		&rest.Route{"POST", "/car/delete", backend.DeleteCar},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":13001", &handler))
}

type Car struct {
	ID    string
	Model string
	Price string
	Year  string
	Brand string
}

func GetCar(w rest.ResponseWriter, r *rest.Request) {
	ID := r.PathParam("ID")
	if ID == "" {
		fmt.Println("GetCar: ID not found in request")
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("GetCar called with ID=" + ID)
	car := Car{}

	queryStr := fmt.Sprintf("select id, model, price, year, brand from car where id = %s", ID)
	err := dbConn.QueryRow(queryStr).Scan(&car.ID, &car.Model,
		&car.Price, &car.Year, &car.Brand)
	switch {
	case err == sql.ErrNoRows:
		fmt.Println("GetCar: no car found with ID " + ID)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case err != nil:
		fmt.Println("GetCar: error in get " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&car)
}

func GetAllCars(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("GetAllCars:called")
	var rows *sql.Rows
	var err error
	rows, err = dbConn.Query("select id, model, price, year, brand) from car order by name")
	if err != nil {
		fmt.Println("GetAllCars: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	cars := make([]Car, 0)
	for rows.Next() {
		car := Car{}
		if err = rows.Scan(&car.ID, &car.Model,
			&car.Price, &car.Year, &car.Brand); err != nil {
			fmt.Println("GetAllCars: error " + err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cars = append(cars, car)
	}
	if err = rows.Err(); err != nil {
		fmt.Println("GetAllCars: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&cars)
}

func AddCar(w rest.ResponseWriter, r *rest.Request) {

	car := Car{}
	err = r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("AddCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("AddCar called")
	queryStr := fmt.Sprintf("insert into car ( model, price, year, brand) values ( '%s', '%s', '%s', '%s') returning id", car.Model, car.Price, car.Year, car.Brand)

	fmt.Println(queryStr)
	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println(err.Error())
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		fmt.Println("inserted car ID " + strconv.Itoa(carid))
	}

	w.WriteJson(&car)
}

func UpdateCar(w rest.ResponseWriter, r *rest.Request) {

	car := Car{}
	err = r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("UpdateCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("UpdateServer:called")
	queryStr := fmt.Sprintf("update car set ( model, price, year, brand) = ('%s', '%s', '%s', '%s') where id = %s returning id", car.Model, car.Price, car.Year, car.Brand, car.ID)

	fmt.Println(queryStr)
	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println("UpdateCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		fmt.Println("UpdateCar:car updated " + carid)
	}

	var status = "OK"
	w.WriteJson(&status)
}

func DeleteCar(w rest.ResponseWriter, r *rest.Request) {

	car := Car{}
	err = r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("DeleteCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryStr := fmt.Sprintf("delete from car where  id=%s returning id", id)
	fmt.Println(queryStr)

	var carid int
	err := dbConn.QueryRow(queryStr).Scan(&carid)
	switch {
	case err != nil:
		fmt.Println("DeleteCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	default:
		fmt.Println("DeleteCar:car deleted " + carid)
	}
	var status = "OK"
	w.WriteJson(&status)
}
