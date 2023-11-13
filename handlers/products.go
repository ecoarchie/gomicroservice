// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/ecoarchie/gomicroservice/data"
	"github.com/gorilla/mux"
)

// Products is a http.Handler
type Products struct {
	l *log.Logger
}

// A list of products returns in a response
// swagger:response productsResponse
type productsResponse struct {
	// All products in the system
	// in: body
	Body []data.Product
}

// swagger:parameters deleteProduct
type ProductIDParameterWrapper struct {
	// the ID of the product to delete from database
	// in: path
	// required: true
	ID int `json:"id"`
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET / products GetProducts
// Returns a list of products
// responses:
//	200: productsResponse

// GetProducts returns the products from data store
func (p Products) GetProducts(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshall json", http.StatusInternalServerError)
	}
}

func (p Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)
	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(&prod)
}

func (p Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) 
	if err != nil {
		http.Error(w, "unable to read id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle PUT product with id ", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)


	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Product not found", http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /{id} products deleteProduct
// responses:
//	201: noContent

// DeleteProduct deletes product with specified ID
func (p Products) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"]) 
	if err != nil {
		http.Error(w, "unable to read id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle DELETE product with id ", id)

	data.DeleteProduct(id)
}

type KeyProduct struct {}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(w, "Unable to unmarshall json", http.StatusBadRequest)
			return
		}

		// validate product
		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(w, "Failed input product fields validation", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)

	})
}