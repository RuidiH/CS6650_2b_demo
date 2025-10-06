package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

// find Product and errorResponses schema from api.yaml line 419 - 485
type Product struct {
	ProductID    int    `json:"product_id"`
	SKU          string `json:"sku"`
	Manufacturer string `json:"manufacturer"`
	CategoryID   int    `json:"category_id"`
	Weight       int    `json:"weight"`
	SomeOtherID  int    `json:"some_other_id"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// helper to send JSON error responses
func writeError(w http.ResponseWriter, status int, errType, message, details string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   errType,
		Message: message,
		Details: details,
	})
}

// in-memory store
var products = make(map[int]Product)

// POST with 404 product not found error implies that the product must already exist in the system before we can add or update its details.
// other words, I write the code as the endpoint should update details for existing products, and return a 404 error if the specified product ID is not found in the store.

func createOrUpdateProduct(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			writeError(w, http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				"An unexpected error occurred",
				fmt.Sprintf("%v", rec))
		}
	}() // this is to add the 500 error response

	//We use mux because it allows us to easily define REST-style endpoints with path variables like /products/{productId}.
	//The call mux.Vars(r) extracts those variables into a map so we can access and use them in our handler.
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT",
			"Invalid input data",
			"Product ID must be a positive integer")
		return
	}

	var p Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT",
			"Invalid input data",
			"Body must be valid JSON")
		return
	}

	// Basic validation
	if p.SKU == "" || p.Manufacturer == "" || p.CategoryID <= 0 || p.ProductID <= 0 {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT",
			"Invalid input data",
			"Missing or invalid product fields")
		return
	}

	// Ensure path ID matches body ID
	if p.ProductID != id {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT",
			"Invalid input data",
			"Product ID in path and body must match")
		return
	}

	// Check if product exists before updating
	if _, exists := products[id]; !exists {
		writeError(w, http.StatusNotFound,
			"NOT_FOUND",
			"Product not found",
			"Cannot add details for a product that does not exist")
		return
	}

	// Update existing product
	products[id] = p
	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			writeError(w, http.StatusInternalServerError,
				"INTERNAL_SERVER_ERROR",
				"An unexpected error occurred",
				fmt.Sprintf("%v", rec))
		}
	}() // this is to add the 500 error response

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest,
			"INVALID_INPUT",
			"Invalid input data",
			"Product ID must be a positive integer")
		return
	}

	p, exists := products[id]
	if !exists {
		writeError(w, http.StatusNotFound,
			"NOT_FOUND",
			"Product not found",
			"No product found with the given ID")
		return
	}

	// Return 200 with JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

func main() {
	// products[1] = Product{
	// 	ProductID: 1,
	// 	SKU: "OLD123",
	// 	Manufacturer: "Acme",
	// 	CategoryID: 10,
	// 	Weight: 100,
	// 	SomeOtherID: 5,
	// }

	for i := 1; i <= 50; i++ {
		products[i] = Product{
			ProductID:    i,
			SKU:          fmt.Sprintf("SKU-%04d", i),
			Manufacturer: "Acme",
			CategoryID:   10,
			Weight:       100,
			SomeOtherID:  5,
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}", getProduct).Methods("GET")
	router.HandleFunc("/products/{id}/details", createOrUpdateProduct).Methods("POST")
	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
}
