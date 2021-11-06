package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest/data"
	"strconv"

	"github.com/gorilla/mux"
)

// ErrInvalidProductPath is an error message when the product path is not valid
var ErrInvalidProductPath = fmt.Errorf("Invalid Path, path should be /products/[id]")

type GenericError struct {
	Message string `json:"message"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

// KeyProduct is a key used for the Product object in the context
type KeyProduct struct{}

// Products handler for getting and updating products
type Products struct {
	l *log.Logger
}

// New Products returns a new products handler with the given logger
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// swagger:route GET /products products
// Returns a list of products
// responses:
//   200: productsResponse

// GetProducts returns the products from the data store
func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	if err := lp.ToJSON(rw); err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

// swagger:route POST /products products createProduct
// Create a new product

// responses:
// 	201: noContentResponse
// 422: errorValidation
// 501: errorResponse

// AddProduct will write a product to the date store
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST product")
	prod := r.Context().Value(KeyProduct{}).(*data.Product)
	data.AddProduct(prod)
	rw.WriteHeader(http.StatusCreated)
}

// swagger:route PUT /products products updateProduct
// Update a products details

// responses:
// 	200: noContentResponse
// 404: errorResponse
// 422: errorValidation

// UpdateProducts will find the product with the id given in the uri params and
// update the product in the data store
func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}
}

// MiddlewardProductValidation will take a request body and validate that the body
// is a valid product.
func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}

		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusInternalServerError)
			return
		}

		// validate the product
		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Error validating product: %s", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}
