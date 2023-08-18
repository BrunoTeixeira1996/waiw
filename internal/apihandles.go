package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

// Function that displayes status code and json response to GET an POST methods
// useful in POST requests
func writeJsonResponseToClient(w http.ResponseWriter, statusCode int, status string) {
	var jsonResp []byte

	resp := make(map[string]string)

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	resp["status"] = status
	jsonResp, _ = json.Marshal(resp)
	w.Write(jsonResp)
}

// GET - returns json object
// POST -> receives json object and adds to database new entry
// Handles "/api/ptw"
func PtwApiHandle(db *Db) http.HandlerFunc {
	var (
		sptw     []Ptw
		category Category
	)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Set("Content-Type", "application/json")

			if err := db.GetPlanToWatch(&sptw); err != nil {
				log.Println("Error while querying ptw:", err)
				return
			}
			json.NewEncoder(w).Encode(sptw)

			// This needs to be done or this will have dup information
			// TODO: there's something better here????
			func() {
				sptw = nil
			}()

		case "POST":
			d := json.NewDecoder(r.Body)
			d.DisallowUnknownFields() // error if user sends extra data

			ptwTemp := struct {
				Name         *string `json:"name"`
				CategoryName *string `json:"category_name"`
			}{}

			if err := d.Decode(&ptwTemp); err != nil {
				// bad JSON or unrecognized json field
				log.Println("Error while decoding plan to watch from POST:", err)
				writeJsonResponseToClient(w, http.StatusBadRequest, err.Error())
				return
			}

			if ptwTemp.Name == nil || ptwTemp.CategoryName == nil {
				log.Println("Missing field from JSON object from POST in plan to watch")
				writeJsonResponseToClient(w, http.StatusBadRequest, "Missing field from JSON object")
				return
			}

			// Check if theres more than what we want
			if d.More() {
				log.Println("Extraneous data after JSON object from POST in plan to watch")
				writeJsonResponseToClient(w, http.StatusBadRequest, "Extraneous data after JSON object")
				return
			}

			// Gets categoryid from name
			if err := db.GetCategoryId(*ptwTemp.CategoryName, &category); err != nil {
				log.Println("Error while extracting category name from POST In plan to watch")
				writeJsonResponseToClient(w, http.StatusInternalServerError, "Error while extracting category name")
				return
			}

			if err := db.InsertPlanToWatch("insert into plan_to_watch (name,category_id) VALUES ($1,$2)", *ptwTemp.Name, category.Id); err != nil {
				log.Println("Error while inserting new plan to watch:", err)
				writeJsonResponseToClient(w, http.StatusInternalServerError, "Error while inserting new plan to watch")
				return
			}

			writeJsonResponseToClient(w, http.StatusOK, "Added new plan to watch")
		}
	}
}