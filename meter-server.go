package main

import (

	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
)

func handler(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "root:root@/go_meter_server");

	checkErr(err);

	_, err2 := db.Exec("INSERT INTO collection (name, icon) VALUES ('Electra (laag)', ''), ('Electra (hoog)', ''	), ('Gas', '')");
	checkErr(err2);

	c := Collection{1, "Electra (laag)", ""};
	m := Measurement{Collection: c, Id:1, Value:200};

	w.Header().Add("Content-type", "application/json");
	json.NewEncoder(w).Encode(m);
}

func collectionListHandler(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Handler string
	}{
		"collectionListHandler",
	}

	json.NewEncoder(w).Encode(data);
}

func collectionDetailHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Handler string
	}{
		"collectionDetailHandler",
	}

	json.NewEncoder(w).Encode(data);
}

func measurementDetailHandler(w http.ResponseWriter, r *http.Request) {
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
	h := RegexpHandler{};

	h.HandleFunc(regexp.MustCompile("^/([a-z0-9]+)/([0-9]{4})/([0-9]{2})/([0-9]{2})$"), measurementDetailHandler)
	h.HandleFunc(regexp.MustCompile("^/([a-z0-9]+)$"), collectionDetailHandler)
	h.HandleFunc(regexp.MustCompile("/"), collectionListHandler)

	http.Handle("/", &h)

	http.ListenAndServe(":8080", nil)
}

