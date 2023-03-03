package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
		gallica_records := get_row_context(query)
		c.JSON(http.StatusOK, gallica_records)
	})
	r.Run()
}

type UserResponse struct {
	records     []GallicaRecord
	num_results int
	origin_urls []string
}

func get_row_context(args SearchArgs) UserResponse {
	gallica_params := rest_args_to_gallica_params(args)
	fmt.Println(gallica_params)
	return UserResponse{
		records:     []GallicaRecord{},
		num_results: 0,
		origin_urls: []string{},
	}

}

func rest_args_to_gallica_params(args SearchArgs) GallicaQueryParams {
	var gram_cql string
	var date_cql string
	var paper_cql string

	if args.link_term != "" && args.link_distance != 0 && len(args.terms) == 1 {
		gram_cql = fmt.Sprintf(`text adj "%s" prox/unit=word/distance=%d "%s"`, args.terms[0], args.link_distance, args.link_term)
	} else if len(args.terms) > 0 {
		gram_cql = `text adj "` + strings.Join(args.terms, `" or text adj "`) + `"`
	}

	if args.start_date != "" && args.end_date != "" {
		date_cql = fmt.Sprintf("gallicapublication_date >= \"%s\" and gallicapublication_date < \"%s\"", args.start_date, args.end_date)
	} else if args.start_date != "" {
		date_cql = fmt.Sprintf("gallicapublication_date >= \"%s\"", args.start_date)
	} else if args.end_date != "" {
		date_cql = fmt.Sprintf("gallicapublication_date < \"%s\"", args.end_date)
	}

	if len(args.codes) > 0 {
		formatted_codes := make([]string, len(args.codes))
		for i, code := range args.codes {
			formatted_codes[i] = code + "_date"
		}
		paper_cql = `arkPress adj "` + strings.Join(formatted_codes, `" or arkPress adj "`) + `"`
	} else if args.source == "periodical" {
		paper_cql = `dc.type all "fascicule"`
	} else if args.source == "book" {
		paper_cql = `dc.type all "monographie"`
	} else {
		paper_cql = `dc.type all "fascicule" or dc.type all "monographie"`
	}

	cql_components := []string{gram_cql, date_cql, paper_cql}
	cql := strings.Join(cql_components, " and ")

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
