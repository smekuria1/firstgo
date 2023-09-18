package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/smekuria1/firstgo/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("GET Products")

	lp := data.GetProducts()

	err := lp.ToJson(rw)

	if err != nil {
		http.Error(rw, "Unable to Marshall JSon", http.StatusInternalServerError)

	}
}

func (p *Products) AddProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("POST Products")

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Printf("Prod: %#v", prod)
	data.AddProducts(&prod)

}

func (p *Products) UpdateProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("PUT Products")
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
	}

	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProdNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			p.l.Println("[Error] deserializing product", err)
			http.Error(rw, "Unable to Unmarshal Json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("[Error] Validating prodect:  ", err)
			http.Error(rw, fmt.Sprintf("Error Validating product: %s", err), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)

		next.ServeHTTP(rw, req)
	})
}
