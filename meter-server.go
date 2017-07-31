package main

import (

	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"flag"
)

type Meter struct {
	db *sql.DB
}

func (meter *Meter) collectionListHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := meter.db.Query("SELECT id, name, icon FROM collection");

	checkErr(err);

	defer rows.Close();

	collections := make([]*Collection, 0);

	for rows.Next() {
		collection := new(Collection);
		err := rows.Scan(&collection.Id, &collection.Name, &collection.Icon);

		checkErr(err);

		collections = append(collections, collection);
	}

	w.Header().Add("Content-type", "application/json");
	json.NewEncoder(w).Encode(collections);
}

func (meter *Meter) collectionDetailHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Handler string
	}{
		"collectionDetailHandler",
	}

	json.NewEncoder(w).Encode(data);
}

func (meter *Meter) measurementDetailHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Handler string
	}{
		"measurementDetailHandlers",
	}

	json.NewEncoder(w).Encode(data);
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	dsnPtr := flag.String("db", "", "-db=user:password@/database");

	flag.Parse();

	db, err := sql.Open("mysql", *dsnPtr);

	checkErr(err);

	meter := Meter{db: db};


	h := RegexpHandler{};

	h.HandleFunc(regexp.MustCompile("^/([a-z0-9]+)/([0-9]{4})/([0-9]{2})/([0-9]{2})$"), meter.measurementDetailHandler)
	h.HandleFunc(regexp.MustCompile("^/([a-z0-9]+)$"), meter.collectionDetailHandler)
	h.HandleFunc(regexp.MustCompile("/"), meter.collectionListHandler)

	http.Handle("/", &h)

	http.ListenAndServe(":8080", nil)
}

