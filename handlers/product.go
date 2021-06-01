package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

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
		p.l.Println("PUT")
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(rw, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Bad URI ", http.StatusBadRequest)
		}
		p.l.Println("got id", id)

		p.updateProducts(id, rw, r)
		return
	}

	//catch all exception
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

//Get all product List
func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {

	lp := data.GetProducts()
	err := lp.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marhsal json ", http.StatusInternalServerError)
	}

}

//POSt add a product
func (p *Products) addProducts(rw http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal Json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

//PUT update a product
func (p *Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request) {

	prod := &data.Product{}
	err := prod.FromJSON(r.Body)

	if err != nil {
		http.Error(rw, "Unable to unmarshal Json", http.StatusBadRequest)
	}

	err = data.UpdateProduct(id, prod)

	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found ", http.StatusNotFound)
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
	}

}

//you are getting a reader is because it hasnt read everything .... if you get a huge amount of data it buffers some
