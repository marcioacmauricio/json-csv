package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"
	"strings"
	"reflect"
	"bytes"
)
// definições json

var fieldDelim = ","
var lineDelim = "\n"
var dataMap = map[string]interface{}{}
func main() {

	// all requests will be converted to map to be subsequently converted to any format
	// get request csv and convert to map format
	getData("http://localhost:5000/file.csv")
	// convert map to json string
	getCsv2JsonString := getJsonString()
	// print
	fmt.Println(getCsv2JsonString)
	// get request csv and convert to map
	getData("http://localhost:5000/file.json")
	// convert map to csv string
	getJson2JsonString := getCsvString()
	// imprime
	fmt.Println(getJson2JsonString)

}

func getJsonString() string {
	ret , err := json.Marshal(dataMap)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.NewBuffer(ret).String()
}
func getCsvString() string {
	var ret bytes.Buffer
	var colString bytes.Buffer
	var valString bytes.Buffer
	i := 0;
	for key, val := range dataMap {
		if i > 0 {
			colString.WriteString(fieldDelim)
			valString.WriteString(fieldDelim)			
		}
		if reflect.TypeOf(val).Kind() == reflect.String {
			colString.WriteString(fmt.Sprint(key))
			valString.WriteString(fmt.Sprint(val))
		}
		
		i++
	}
	ret.WriteString(colString.String())
	ret.WriteString(lineDelim)
	ret.WriteString(valString.String())
	return ret.String()
}
// adquire data através da 
func getData(url string) {
	var jsonString bytes.Buffer
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	content, err := ioutil.ReadAll(res.Body)
	data := bytes.NewBuffer(content).String()
	parts := strings.Split(res.Header["Content-Type"][0],"; ")
	tp := parts[0]
	chkType := false
	switch tps := tp
	tps {
		case "text/csv":
			lines := strings.Split(data,lineDelim)
			columns := strings.Split(lines[0],fieldDelim)
			values := strings.Split(lines[1],fieldDelim)
			jsonString.WriteString("{")
			for i := 0; i < len(columns); i++ {
				if i > 0 {
					jsonString.WriteString(",")
				}
				jsonString.WriteString(`"`)
				jsonString.WriteString(columns[i])
				jsonString.WriteString(`":`)
				jsonString.WriteString(`"`)
				jsonString.WriteString(values[i])
				jsonString.WriteString(`"`)
			}
			jsonString.WriteString("}")
		case "application/json":
			jsonString.WriteString(data)
		default:
			jsonString.WriteString("{}")
			chkType = true
	}
	if chkType {
		log.Fatal("Tipo não suportado")
	} 
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte(jsonString.String()), &dataMap)
	if err != nil {
		log.Fatal(err)
	}

}