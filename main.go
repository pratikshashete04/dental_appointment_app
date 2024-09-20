package main

import (
	"fmt"
	"log"
	"net/http"

	util "dental_appointment_app/utilities"

	"github.com/gorilla/mux"
)

func main() {

	// Using gorilla/mux router for flexible routing
	router := mux.NewRouter()

	// Setting up routes
	router.HandleFunc("/api/patient", util.AddNewPatient).Methods("POST")        // POST to create patient
	router.HandleFunc("/api/patients", util.GetAllPatients).Methods("GET")       // GET to retrieve all patients
	router.HandleFunc("/api/patient/{id}", util.GetPatientByID).Methods("GET")   // GET to retrieve a single patient by ID
	router.HandleFunc("/api/patient/{id}", util.UpdatePatient).Methods("PUT")    // PUT to update a patient by ID
	router.HandleFunc("/api/patient/{id}", util.DeletePatient).Methods("DELETE") // DELETE to remove a patient by ID

	// Start server
	fmt.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
