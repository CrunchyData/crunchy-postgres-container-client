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
	"errors"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/crunchydata/crunchy-postgres-container-client/backend"
	"log"
	"net/http"
	"os"
)

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
		&rest.Route{"POST", "/car/add", AddCar},
		&rest.Route{"POST", "/car/update", UpdateCar},
		&rest.Route{"GET", "/car/list", GetAllCars},
		&rest.Route{"GET", "/car/:ID.", GetCar},
		&rest.Route{"POST", "/car/delete", DeleteCar},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":13001", &handler))
}

func GetCar(w rest.ResponseWriter, r *rest.Request) {
	ID := r.PathParam("ID")
	if ID == "" {
		fmt.Println("GetCar: ID not found in request")
		rest.Error(w, "ID not passed", http.StatusBadRequest)
		return
	}

	fmt.Println("GetCar called with ID=" + ID)

	dbConn, err := getConnection()
	if err != nil {
		fmt.Println("GetCar: error in connection" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var car backend.Car
	car, err = backend.GetCar(dbConn, ID)
	if err != nil {
		fmt.Println("GetCar: error" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&car)
}

func GetAllCars(w rest.ResponseWriter, r *rest.Request) {
	fmt.Println("GetAllCars:called")

	dbConn, err := getConnection()
	if err != nil {
		fmt.Println("GetAllCars: error in connection" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cars, err := backend.GetAllCars(dbConn)
	if err != nil {
		fmt.Println("GetAllCars: error " + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&cars)
}

func AddCar(w rest.ResponseWriter, r *rest.Request) {

	car := backend.Car{}
	err := r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("AddCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = getConnection()
	if err != nil {
		fmt.Println("AddCar: error in connection" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var carid string
	carid, err = backend.AddCar(dbConn, car)

	if err != nil {
		fmt.Println(err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	car.ID = carid

	w.WriteJson(&car)
}

func UpdateCar(w rest.ResponseWriter, r *rest.Request) {

	car := backend.Car{}
	err := r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("UpdateCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = getConnection()
	if err != nil {
		fmt.Println("UpdateCar: error in connection" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = backend.UpdateCar(dbConn, car)
	if err != nil {
		fmt.Println("UpdateCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status = "OK"
	w.WriteJson(&status)
}

func DeleteCar(w rest.ResponseWriter, r *rest.Request) {

	car := backend.Car{}
	err := r.DecodeJsonPayload(&car)
	if err != nil {
		fmt.Println("DeleteCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dbConn *sql.DB
	dbConn, err = getConnection()
	if err != nil {
		fmt.Println("DeleteCar: error in connection" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = backend.DeleteCar(dbConn, car.ID)
	if err != nil {
		fmt.Println("DeleteCar: error in decode" + err.Error())
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var status = "OK"
	w.WriteJson(&status)
}

func getConnection() (*sql.DB, error) {
	var dbConn *sql.DB
	var user = os.Getenv("PG_USER")
	if user == "" {
		return dbConn, errors.New("PG_USER env var not set")
	}
	var host = os.Getenv("PG_HOST")
	if host == "" {
		return dbConn, errors.New("PG_HOST env var not set")
	}
	var password = os.Getenv("PG_PASSWORD")
	if password == "" {
		return dbConn, errors.New("PG_PASSWORD env var not set")
	}
	var database = os.Getenv("PG_DATABASE")
	if database == "" {
		return dbConn, errors.New("PG_DATABASE env var not set")
	}

	var port = "5432"

	dbConn, err := sql.Open("postgres", "sslmode=disable user="+user+" host="+host+" port="+port+" dbname="+database+" password="+password)
	return dbConn, err

}
