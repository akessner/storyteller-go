package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//TEST DATA

var filters []Filter = []Filter{
	{filter_id: 1, parent_id: "", sort_order: "", story_id: "", property_name: "", property_class: "", operator: "", value: "", created_at: "2015-01-26 16:30:15", updated_at: "2015-01-26 16:30:15"},
	{filter_id: 12, parent_id: "1", sort_order: "", story_id: "", property_name: "age", property_class: "Integer", operator: "<", value: "35", created_at: "2015-01-28 14:37:20", updated_at: "2015-01-28 14:37:20"},
	{filter_id: 17, parent_id: "1", sort_order: "", story_id: "", property_name: "gender", property_class: "String", operator: "==", value: "m", created_at: "2015-02-03 15:28:24", updated_at: "2015-02-03 15:28:24"},
	{filter_id: 18, parent_id: "17", sort_order: "", story_id: "", property_name: "age", property_class: "Integer", operator: ">=", value: "25", created_at: "2015-02-03 15:28:37", updated_at: "2015-02-03 15:28:37"},
	{filter_id: 19, parent_id: "18", sort_order: "", story_id: "1", property_name: "", property_class: "", operator: "", value: "", created_at: "2015-02-03 15:28:43", updated_at: "2015-02-03 15:28:43"},
	{filter_id: 20, parent_id: "12", sort_order: "", story_id: "", property_name: "gender", property_class: "String", operator: "==", value: "f", created_at: "2015-02-03 15:28:52", updated_at: "2015-02-03 15:28:52"},
	{filter_id: 21, parent_id: "20", sort_order: "", story_id: "2", property_name: "", property_class: "", operator: "", value: "", created_at: "2015-02-08 09:25:15", updated_at: "2015-02-08 09:25:15"},
}

func sampleProperties() ([]string, []string, []string) {
	expectedProperties := []string{"age", "gender"}
	expectedOperators := []string{">=", "=="}
	expectedValues := []string{"25", "m"}
	//expectedProperties, expectedOperators, expectedValues
	return expectedProperties, expectedOperators, expectedValues
}

func sampleURL() *url.URL {

	my_url, _ := url.Parse(sampleURLString())
	return my_url
}

func sampleURLString() string {
	params := "age=36&gender=m"
	return fmt.Sprintf("http://api.url/get_stories?%s", params)
}

//TESTS

func TestStoryHandler(t *testing.T) {
	expectedBody := "{\"age\":\"36\",\"gender\":\"m\"}"

	fakeResponseWriter := httptest.NewRecorder()
	req, err := http.NewRequest("GET", sampleURLString(), nil)
	assert.Nil(t, err)

	StoryHandler(fakeResponseWriter, req)
	assert.Equal(t, expectedBody, fakeResponseWriter.Body.String())
}

func TestGetParams(t *testing.T) {
	expectedParams := map[string]string{
		"gender": "m",
		"age":    "36",
	}
	actualParams := getParams(sampleURL())
	assert.Equal(t, expectedParams, actualParams)
}

func TestConstructTableFromFilter(t *testing.T) {
	properties, property_operators, property_values := sampleProperties()
	expectedRow1 := StoryConfig{1, "bonus game", 1, properties, property_operators, property_values}
	expectedRow2 := StoryConfig{1, "bonus game", 2,
		[]string{"gender", "age"},
		[]string{"==", "<"},
		[]string{"f", "35"}}
	actualTable := buildFilterTable(filters, 1, "bonus game")
	assert.Equal(t, expectedRow1, actualTable[0])
	assert.Equal(t, expectedRow2, actualTable[1])

}

func TestConstructTableForMultipleGameEvents(t *testing.T) {
	var actualTable []StoryConfig
	actualTable = append(actualTable, buildFilterTable(filters, 1, "bonus game")...)
	actualTable = append(actualTable, buildFilterTable(filters, 1, "5 in a row")...)
	actualTable = append(actualTable, buildFilterTable(filters, 2, "5 in a row")...)

	assert.Equal(t, 6, len(actualTable))
}

func TestParseFilters(t *testing.T) {
	expectedProperties, expectedOperators, expectedValues := sampleProperties()
	actualProperties, actualOperators, actualValues := parseFilter(filters, "18")

	assert.Equal(t, expectedProperties, actualProperties)
	assert.Equal(t, expectedOperators, actualOperators)
	assert.Equal(t, expectedValues, actualValues)

}

func TestGetStoryFromFilterTable_story1(t *testing.T) {
	expectedStory := "1"
	scope := getParams(sampleURL())
	table := buildFilterTable(filters, 1, "bonus game")
	actualStory := getStory(scope, table)
	assert.Equal(t, expectedStory, actualStory)

}
func TestGetStoryFromFilterTable_story2(t *testing.T) {
	expectedStory := "2"
	testUrl, _ := url.Parse("http://api/get_stories?age=20&gender=f")
	scope := getParams(testUrl)
	table := buildFilterTable(filters, 1, "bonus game")
	actualStory := getStory(scope, table)
	assert.Equal(t, expectedStory, actualStory)

}
func TestGetStoryFromFilterTable_storyNone(t *testing.T) {
	expectedStory := ""
	testUrl, _ := url.Parse("http://api/get_stories?age=50&gender=h")
	scope := getParams(testUrl)
	table := buildFilterTable(filters, 1, "bonus game")
	actualStory := getStory(scope, table)
	assert.Equal(t, expectedStory, actualStory)

}

func TestEvalPropsGreaterEqualThan(t *testing.T) {
	expectedResult := true
	actualResult := evalProperty(">=", "25", "36")
	assert.Equal(t, expectedResult, actualResult)
}
func TestEvalPropsGreaterThan(t *testing.T) {
	expectedResult := false
	actualResult := evalProperty(">", "10", "5")
	assert.Equal(t, expectedResult, actualResult)
}
func TestEvalPropsLessThan(t *testing.T) {
	expectedResult := false
	actualResult := evalProperty("<", "105", "105")
	assert.Equal(t, expectedResult, actualResult)
}
func TestEvalPropsEqualString(t *testing.T) {
	expectedResult := true
	actualResult := evalProperty("==", "105", "105")
	assert.Equal(t, expectedResult, actualResult)
}
func TestEvalPropsEqualStringAndInt(t *testing.T) {
	expectedResult := false
	actualResult := evalProperty("==", "105", "s105")
	assert.Equal(t, expectedResult, actualResult)
}
