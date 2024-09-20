package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.etcd.io/bbolt"
)

// AddNewPatient handles the HTTP requests for storing patient data
func AddNewPatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Store patient data in the BoltDB
	if dbHandle := handle(); dbHandle != nil {
		err = db.Update(func(tx *bbolt.Tx) error {
			b := tx.Bucket([]byte("Patients"))

			// Generate a unique ID for the patient (using current time)
			id := time.Now().Format("20060102150405")

			// Serialize the patient data to JSON
			patientData, err := json.Marshal(patient)
			if err != nil {
				return err
			}

			// Store the patient data in the bucket using the generated ID as the key
			return b.Put([]byte(id), patientData)
		})
	}

	if err != nil {
		http.Error(w, "Failed to store patient data", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Patient data stored successfully"))
}

// GetAllPatients handles the HTTP requests for retrieving all patients
func GetAllPatients(w http.ResponseWriter, r *http.Request) {
	var patients []Patient

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Patients"))

		return b.ForEach(func(k, v []byte) error {
			var patient Patient

			// Deserialize the stored JSON data back into a Patient struct
			err := json.Unmarshal(v, &patient)
			if err != nil {
				return err
			}

			// Append the patient to the list
			patients = append(patients, patient)
			return nil
		})
	})

	if err != nil {
		http.Error(w, "Failed to retrieve patient data", http.StatusInternalServerError)
		return
	}

	// Respond with all patients in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patients)
}

// GetPatientByID handles the HTTP requests for retrieving a single patient by ID
func GetPatientByID(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	vars := mux.Vars(r)
	id := vars["id"]

	var patient Patient

	// Retrieve the patient data by ID from BoltDB
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Patients"))

		// Retrieve the patient data by ID
		patientData := b.Get([]byte(id))
		if patientData == nil {
			return fmt.Errorf("Patient not found")
		}

		// Deserialize the JSON data back into a Patient struct
		return json.Unmarshal(patientData, &patient)
	})

	if err != nil {
		http.Error(w, "Patient not found", http.StatusNotFound)
		return
	}

	// Respond with the patient data in JSON format
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(patient)
}

// UpdatePatient handles the HTTP requests for updating patient information
func UpdatePatient(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedPatient Patient

	// Decode the request body to get the updated patient information
	err := json.NewDecoder(r.Body).Decode(&updatedPatient)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update the patient data in the BoltDB
	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Patients"))

		// Retrieve the existing patient data by ID
		patientData := b.Get([]byte(id))
		if patientData == nil {
			return fmt.Errorf("Patient not found")
		}

		// Serialize the updated patient data to JSON
		updatedData, err := json.Marshal(updatedPatient)
		if err != nil {
			return err
		}

		// Store the updated patient data in the bucket using the same ID
		return b.Put([]byte(id), updatedData)
	})

	if err != nil {
		if err.Error() == "Patient not found" {
			http.Error(w, "Patient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update patient data", http.StatusInternalServerError)
		}
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Patient data updated successfully"))
}

// DeletePatient handles the HTTP requests for deleting a patient by ID
func DeletePatient(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL path
	vars := mux.Vars(r)
	id := vars["id"]

	// Delete the patient data from the BoltDB
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("Patients"))

		// Check if the patient exists
		patientData := b.Get([]byte(id))
		if patientData == nil {
			return fmt.Errorf("Patient not found")
		}

		// Delete the patient data by ID
		return b.Delete([]byte(id))
	})

	if err != nil {
		if err.Error() == "Patient not found" {
			http.Error(w, "Patient not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete patient data", http.StatusInternalServerError)
		}
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Patient data deleted successfully"))
}
