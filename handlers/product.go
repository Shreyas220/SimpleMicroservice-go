package handlers

import (
	"log"
	"net/http"

	"github.com/Shreyas220/SimpleMicroservice-go/data"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

//Main handler
func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProducts(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		//except the ID in the URI
		p := r.URL.Path
		return
	}

	//catch all exception
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal Json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

//Get all product List
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marhsal json ", http.StatusInternalServerError)
	}

}

//you are getting a reader is because it hasnt read everything .... if you get a huge amount of data it buffers some
