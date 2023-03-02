package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type SearchArgs struct {
	year          int
	month         int
	day           int
	end_year      int
	end_month     int
	terms         []string
	codes         []string
	cursor        int
	limit         int
	link_term     string
	link_distance int
	source        string
}

func cast_query_args_to_int(query_args url.Values, key string) (int, error) {
	val, err := strconv.Atoi(query_args.Get(key))
	if err != nil {
		fmt.Println(err)
	}
	return val, err
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var query SearchArgs
		request_args := r.URL.Query()
		query.year, _ = cast_query_args_to_int(request_args, "year")
		query.month, _ = cast_query_args_to_int(request_args, "month")
		query.day, _ = cast_query_args_to_int(request_args, "day")
		query.end_year, _ = cast_query_args_to_int(request_args, "end_year")
		query.end_month, _ = cast_query_args_to_int(request_args, "end_month")
		query.cursor, _ = cast_query_args_to_int(request_args, "cursor")
		query.limit, _ = cast_query_args_to_int(request_args, "limit")
		query.link_distance, _ = cast_query_args_to_int(request_args, "link_distance")
		query.terms = request_args["terms"]
		query.codes = request_args["codes"]
		query.link_term = request_args.Get("link_term")
		query.source = request_args.Get("source")
		fmt.Println(query)
		fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}

func get_row_context(args SearchArgs) {
	total_records := 0
	original_url := ""

}

type GallicaRecord struct {
	data string
}

type GallicaWrapper interface {
	get() []GallicaRecord
	parse() []GallicaRecord
}
