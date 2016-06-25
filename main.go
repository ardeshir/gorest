package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Link struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon"`
}

var linkStore = make(map[string]Link)

var id int = 0

func PostLinkHandler(w http.ResponseWriter, r *http.Request) {
	var link Link
	err := json.NewDecoder(r.Body).Decode(&link)
	if err != nil {
		panic(err)
	}

	link.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	linkStore[k] = link

	j, err := json.Marshal(link)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func GetLinksHandler(w http.ResponseWriter, r *http.Request) {
	var links []Link
	for _, v := range linkStore {
		links = append(links, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(links)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func GetALinkHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	k := vars["id"]

	if link, ok := linkStore[k]; ok {
		j, err := json.Marshal(link)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(j)

	} else {

		log.Printf("Could not find key of Link %s to show", k)
		w.WriteHeader(http.StatusNoContent)
	}

}

func PutLinkHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	k := vars["id"]
	var linkToUpd Link

	err = json.NewDecoder(r.Body).Decode(&linkToUpd)
	if err != nil {
		panic(err)
	}

	if link, ok := linkStore[k]; ok {
		linkToUpd.CreatedOn = link.CreatedOn
		delete(linkStore, k)
		linkStore[k] = linkToUpd
	} else {
		log.Printf("Could not find key of Link %s to update", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteLinkHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if _, ok := linkStore[k]; ok {
		delete(linkStore, k)
	} else {
		log.Printf("Could not find key of Link %s to detele", k)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/links", GetLinksHandler).Methods("GET")
	r.HandleFunc("/api/links", PostLinkHandler).Methods("POST")
	r.HandleFunc("/api/links/{id}", GetALinkHandler).Methods("GET")
	r.HandleFunc("/api/links/{id}", PutLinkHandler).Methods("PUT")
	r.HandleFunc("/api/links/{id}", DeleteLinkHandler).Methods("DELETE")

	server := &http.Server{
		Addr:    ":9090",
		Handler: r,
	}
	log.Println("We're up on port 9090...")
	server.ListenAndServe()
}
