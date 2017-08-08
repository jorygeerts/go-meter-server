package main

import (

	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"flag"
	"strings"
	"time"
	"errors"
	"strconv"
)

type Meter struct {
	db *sql.DB
}

func (meter *Meter) getCollection(id int) (*Collection, error)  {
	rows, err := meter.db.Query("SELECT id, name, icon FROM collection WHERE id = ? LIMIT 1", id);

	checkErr(err);

	defer rows.Close();

	if rows.Next() == false {
		return nil, errors.New("No collection with id " + string(id));
	}
	collection := new(Collection);
	rows.Scan(&collection.Id, &collection.Name, &collection.Icon);

	return collection, nil;
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
	id, err := strconv.Atoi(parts[1]);

	if err != nil {
		http.NotFound(w, r);
		return;
	}

	collection, err := meter.getCollection(id);

	if err != nil {
		http.NotFound(w, r);
		return;
	}

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

	collectionId, err := strconv.Atoi(r.URL.Path[1:2]);
	year := r.URL.Path[3:7];
	month := r.URL.Path[8:10];
	day := r.URL.Path[11:13];

	if err != nil {
		http.NotFound(w, r);
		return;
	}

	collection, err := meter.getCollection(collectionId);

	if err != nil {
		http.NotFound(w, r);
		return;
	}

	dash := string('-');
	percent := string('%');
	date := year + dash + month + dash + day;

	rows, err := meter.db.Query("SELECT id, measured_value, measured_at FROM measurement WHERE collection_id = ? AND measured_at LIKE ?", collection.Id, percent + date + percent);

	checkErr(err);

	defer rows.Close();

	if rows.Next() == false {
		http.NotFound(w, r);
		return;
	}

	var measureTime []byte
	measurement := new(Measurement);
	measurement.Collection = *collection;

	rows.Scan(&measurement.Id, &measurement.Value, &measureTime);
	measurement.Date, _ = time.Parse("2006-01-02 15:04:00", string(measureTime));

	w.Header().Add("Content-type", "application/json");
	json.NewEncoder(w).Encode(measurement);
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

