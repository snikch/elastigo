package elastigo

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// Test all aggregate types and nested aggregations
func TestAggregateDsl(t *testing.T) {

	min := Aggregate("min_price").Min("price")
	max := Aggregate("max_price").Max("price")
	sum := Aggregate("sum_price").Sum("price")
	avg := Aggregate("avg_price").Avg("price")
	stats := Aggregate("stats_price").Stats("price")
	extendedStats := Aggregate("extended_stats_price").ExtendedStats("price")
	valueCount := Aggregate("value_count_price").ValueCount("price")
	percentiles := Aggregate("percentiles_price").Percentiles("price")
	cardinality := Aggregate("cardinality_price").Cardinality("price", true, 50)
	global := Aggregate("global").Global()
	missing := Aggregate("missing_price").Missing("price")
	terms := Aggregate("terms_price").Terms("price")
	termsSize := Aggregate("terms_price_size").TermsWithSize("price", 0)
	significantTerms := Aggregate("significant_terms_price").SignificantTerms("price")
	histogram := Aggregate("histogram_price").Histogram("price", 50).MinDocCount(float64(0)).ExtendedBounds(nil, float64(0))
	histogramOther := Aggregate("histogram_other").Histogram("price", 50)

	time2013, _ := time.Parse(time.RFC3339, "2013-01-01T00:00:00+00:00")
	time2014, _ := time.Parse(time.RFC3339, "2014-01-01T00:00:00+00:00")
	dateAgg := Aggregate(
		"articles_over_time",
	).DateHistogram(
		"date",
		"month",
	).MinDocCount(
		float64(0),
	).ExtendedBounds(
		time2013,
		time2014,
	)
	dateAgg.Aggregates(
		min,
		max,
		sum,
		avg,
		stats,
		extendedStats,
		valueCount,
		percentiles,
		cardinality,
		global,
		missing,
		terms,
		termsSize,
		significantTerms,
		histogram,
		histogramOther,
	)

	qry := Search("github").Aggregates(dateAgg)

	marshaled, err := json.MarshalIndent(qry.AggregatesVal, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal AggregatesVal: %s", err.Error())
		return
	}

	assertJsonMatch(
		t,
		marshaled,
		[]byte(`
			{
				"articles_over_time": {
					"date_histogram" : {
						"field" : "date",
						"interval" : "month",
						"min_doc_count": 0,
						"extended_bounds": {
							"min": "2013-01-01T00:00:00Z",
							"max": "2014-01-01T00:00:00Z"
						}
					},
					"aggregations": {
						"min_price":{
							"min": { "field": "price" }
						},
						"max_price":{
							"max": { "field": "price" }
						},
						"sum_price":{
							"sum": { "field": "price" }
						},
						"avg_price": {
							"avg": { "field": "price" }
						},
						"stats_price":{
							"stats": { "field": "price" }
						},
						"extended_stats_price":{
							"extended_stats": { "field": "price" }
						},
						"value_count_price":{
							"value_count": { "field": "price" }
						},
						"percentiles_price":{
							"percentiles": { "field": "price" }
						},
						"cardinality_price":{
							"cardinality": { "field": "price", "precision_threshold": 50 }
						},
						"global":{
							"global": {}
						},
						"missing_price":{
							"missing": { "field": "price" }
						},
						"terms_price":{
							"terms": { "field": "price" }
						},
						"terms_price_size":{
							"terms": { "field": "price", "size": 0 }
						},
						"significant_terms_price":{
							"significant_terms": { "field": "price" }
						},
						"histogram_other":{
							"histogram": {
								"field": "price",
								"min_doc_count": 1,
								"interval": 50
							}
						},
						"histogram_price":{
							"histogram": {
								"field": "price",
								"interval": 50,
								"min_doc_count": 0,
								"extended_bounds": {
									"max": 0
								}
							}
						}
					}
				}
			}
	`),
	)

}

func TestAggregateFilter(t *testing.T) {

	avg := Aggregate("avg_price").Avg("price")

	dateAgg := Aggregate("in_stock_products").Filter(
		Range().Field("stock").Gt(0),
	)

	dateAgg.Aggregates(
		avg,
	)

	qry := Search("github").Aggregates(dateAgg)

	marshaled, err := json.MarshalIndent(qry.AggregatesVal, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal AggregatesVal: %s", err.Error())
		return
	}

	assertJsonMatch(
		t,
		marshaled,
		[]byte(`
	{
		"in_stock_products" : {
			"filter" : {
				"range" : { "stock" : { "gt" : 0 } }
			},
			"aggregations" : {
				"avg_price" : { "avg" : { "field" : "price" } }
			}
		}
	}
	`),
	)
}

func assertJsonMatch(t *testing.T, match, expected []byte) {
	var m interface{}
	var e interface{}

	err := json.Unmarshal(expected, &e)
	if err != nil {
		t.Errorf("Failed to unmarshal expectation: %s", err.Error())
		return
	}
	err = json.Unmarshal(match, &m)
	if err != nil {
		t.Errorf("Failed to unmarshal match: %s", err.Error())
		return
	}

	if !reflect.DeepEqual(m, e) {
		t.Errorf("Expected %s but got %s", string(expected), string(match))
		return
	}

}
