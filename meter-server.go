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

func db() (db *sql.DB){
	db, err := sql.Open("mysql", "root:root@/go_meter_server");

	checkErr(err);

	return db
}

func collectionListHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := db().Query("SELECT id, name, icon FROM collection");

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

