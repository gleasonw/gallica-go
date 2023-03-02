package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
)

type SearchArgs struct {
	terms         []string
	codes         []string
	cursor        int
	limit         int
	link_term     string
	link_distance int
	source        string
	sort          string
	start_date    string
	end_date      string
}

func cast_query_args_to_int(query_args url.Values, key string) (int, error) {
	val, err := strconv.Atoi(query_args.Get(key))
	if err != nil {
		fmt.Println(err)
	}
	return val, err
}

func main() {
	r := gin.Default()
	r.GET("/search", func(c *gin.Context) {
		var query SearchArgs
		get_date_string := func(year int, month int) string {
			if month != 0 {
				return fmt.Sprintf("%d-%d", year, month)
			} else {
				return fmt.Sprintf("%d", year)
			}
		}
		request_args := c.Request.URL.Query()
		year, _ := cast_query_args_to_int(request_args, "year")
		month, _ := cast_query_args_to_int(request_args, "month")
		end_year, _ := cast_query_args_to_int(request_args, "end_year")
		end_month, _ := cast_query_args_to_int(request_args, "end_month")
		query.start_date = get_date_string(year, month)
		query.end_date = get_date_string(end_year, end_month)
		query.cursor, _ = cast_query_args_to_int(request_args, "cursor")
		query.limit, _ = cast_query_args_to_int(request_args, "limit")
		query.link_distance, _ = cast_query_args_to_int(request_args, "link_distance")
		query.terms = request_args["terms"]
		query.codes = request_args["codes"]
		query.link_term = request_args.Get("link_term")
		query.source = request_args.Get("source")
		query.sort = request_args.Get("sort")
		c.JSON(http.StatusOK, gin.H{"start": start, "end": end})
	})
	r.Run()
}

func get_row_context(args SearchArgs) {
}

func rest_args_to_gallica_params(args SearchArgs) GallicaQueryParams {
	cql := "make string"
	
	return GallicaQueryParams{
		operation:      "searchRetrieve",
		exactSearch:    true,
		version:        1.2,
		startRecord:    args.cursor,
		maximumRecords: args.limit,
		query:          cql,
		collapsing:     false,
	}

}

type GallicaQueryParams struct {
	operation      string
	exactSearch    bool
	version        float32
	startRecord    int
	maximumRecords int
	query          string
	collapsing     bool
}

type GallicaRecord struct {
}

type GallicaWrapper interface {
	get() []GallicaRecord
	parse() []GallicaRecord
}
