package main

import (

	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"flag"
	"fmt"
	"strings"
	"time"
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

	parts := strings.Split(r.URL.Path, "/");
	id := parts[1];

	rows, err := meter.db.Query("SELECT id, name, icon FROM collection WHERE id = ? LIMIT 1", id);

	checkErr(err);

	defer rows.Close();

	if rows.Next() == false {
		http.NotFound(w, r);
		return;
	}
	collection := new(Collection);
	rows.Scan(&collection.Id, &collection.Name, &collection.Icon);

	mRows, err := meter.db.Query("SELECT id, measured_value, measured_at FROM measurement WHERE collection_id = ? ORDER BY measured_at DESC", collection.Id);

	checkErr(err);

	defer mRows.Close();

	collection.Measurements = make([]Measurement, 0);

	var measureTime []byte

	for mRows.Next() {
		measurement := new(Measurement);
		mRows.Scan(&measurement.Id, &measurement.Value, &measureTime);
		measurement.Date, _ = time.Parse("2006-01-02 15:04:00", string(measureTime));

		collection.Measurements = append(collection.Measurements, *measurement);
	}

	w.Header().Add("Content-type", "application/json");
	json.NewEncoder(w).Encode(collection);
}

func (meter *Meter) measurementDetailHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "URL: %s\n", r.URL.Path[1:]);
	fmt.Fprintf(w, "URL: %s\n", r.URL.Path[1:1]);
	fmt.Fprintf(w, "URL: %s\n", r.URL.Path[1:2]);
	fmt.Fprintf(w, "URL: %s\n", r.URL.Path[3:7]);
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	dsnPtr := flag.String("db", "", "-db=user:password@/database?parseTime=true");

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

