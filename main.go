package main

import (
	"encoding/json"
	//"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Story struct {
	Title       string
	Description string
	Img_url     string
	Social_id   int
	Click_url   string
	Platform_id string
	Event_id    string
}

type StoryConfig struct {
	app_id     int
	game_event string
	story_id   int
	properties []string
	operators  []string
	values     []string
}

type Filter struct {
	filter_id      int
	parent_id      string
	sort_order     string
	story_id       string
	property_name  string
	property_class string
	operator       string
	value          string
	created_at     string
	updated_at     string
}

var storyConfigTable []StoryConfig

func main() {
	println("Hello Dave...")
	http.HandleFunc("/", hello)
	http.HandleFunc("/get_stories", StoryHandler)
	http.ListenAndServe(":8080", nil)
	//storyConfigTable = buildFilterTable()
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!"))
	//fmt.Fprint(w, "Test")
	//fmt.Println("Connect Succesful")
}

func StoryHandler(w http.ResponseWriter, r *http.Request) {
	params := getParams(r.URL)
	js, err := json.Marshal(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func getParams(a_url *url.URL) map[string]string {
	params := make(map[string]string)
	m, _ := url.ParseQuery(a_url.RawQuery)
	for key, value := range m {
		params[key] = value[len(value)-1]
	}
	return params
}

func retrieveStoryTree(app int) []interface{} {
	var storyTree []interface{}
	var bytes []byte = returnJSonForApp(app)
	json.Unmarshal(bytes, &storyTree)
	return storyTree
}

func returnJSonForApp(app int) []byte {
	b := []byte(`{"game_event":"5 in a row", "age":25, "gender":"male", "story":"Click now and save!" /}`)
	return b
}

func buildFilterTable(filters []Filter, app_id int, game_event string) []StoryConfig {
	var storyTable []StoryConfig
	for _, filter := range filters {
		if filter.story_id != "" {
			storyId, _ := strconv.Atoi(filter.story_id)
			parsedProperties, parsedOperators, parsedValues := parseFilter(filters, filter.parent_id)
			storyTable = append(storyTable, StoryConfig{
				app_id:     app_id,
				game_event: game_event,
				story_id:   storyId,
				properties: parsedProperties,
				operators:  parsedOperators,
				values:     parsedValues,
			})
		}
	}
	return storyTable
}

func parseFilter(filters []Filter, parentId string) (parsedProperties, parsedOperators, parsedValues []string) {
	parentID, _ := strconv.Atoi(parentId)
	for _, filter := range filters {

		if filter.filter_id == parentID && filter.parent_id != "" {
			parsedProperties = append(parsedProperties, filter.property_name)
			parsedOperators = append(parsedOperators, filter.operator)
			parsedValues = append(parsedValues, filter.value)
			properties, operators, values := parseFilter(filters, filter.parent_id)
			parsedProperties = append(parsedProperties, properties...)
			parsedOperators = append(parsedOperators, operators...)
			parsedValues = append(parsedValues, values...)
		}

	}

	return parsedProperties, parsedOperators, parsedValues
}

func getStory(scope map[string]string, table []StoryConfig) string {
	story_id := ""
	for _, storyConfig := range table {
		propertiesPass := 0
		for key, value := range scope {
			for i := 0; i < len(storyConfig.properties); i++ {
				if storyConfig.properties[i] == key {
					if evalProperty(storyConfig.operators[i], storyConfig.values[i], value) {
						propertiesPass++
					}
				}
			}
		}
		if propertiesPass == len(storyConfig.properties) {
			story_id = strconv.Itoa(storyConfig.story_id)
		}
	}
	return story_id
}

func evalProperty(operator string, targetValue string, actualValue string) bool {

	actualValueNumber, err := strconv.Atoi(actualValue)
	targetValueNumber, err := strconv.Atoi(targetValue)
	if err != nil && operator == "==" {
		return actualValue == targetValue
	}
	switch {
	case operator == "==":
		return actualValueNumber == targetValueNumber
	case operator == ">=":
		return actualValueNumber >= targetValueNumber
	case operator == "<=":
		return actualValueNumber <= targetValueNumber
	case operator == "<":
		return actualValueNumber < targetValueNumber
	case operator == ">":
		return actualValueNumber > targetValueNumber
	}

	return false
}
