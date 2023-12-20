package main

import (
	"fmt"
	"log"
	"net/http"
)

func SearchAddressHandler(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("query")

	exists, err := repository.CheckAddressInDatabase(searchQuery)
	if err != nil {
		log.Printf("Error checking address in database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if exists {
		respondWithMessage(w, fmt.Sprintf("Address found in the database: %s", searchQuery), http.StatusOK)
		return
	}

	if err := processNewAddress(searchQuery); err != nil {
		log.Printf("Error processing new address: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	respondWithMessage(w, fmt.Sprintf("Address found from Dadata.ru: %s", searchQuery), http.StatusOK)
}

func processNewAddress(searchQuery string) error {
	addressID, err := repository.SaveAddress(searchQuery)
	if err != nil {
		return err
	}

	if err := repository.SaveSearchHistory(searchQuery); err != nil {
		return err
	}

	searchHistoryID := 1 // Это значение должно быть получено корректно, а не быть захардкоженным
	return repository.LinkAddressToSearchHistory(searchHistoryID, addressID)
}

func respondWithMessage(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(message))
}
