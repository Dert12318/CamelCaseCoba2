package main

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Name struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	MiddleName string `json:"MiddleName,omitempty"`
}
type Name2 struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	MiddleName string `json:"MiddleName"`
}
type Person struct {
	Name     Name   `json:"Nameku"`
	FullName Name   `json:"FullName"`
	Ages     string `json:"Ages,omitempty"`
}
type Person2 struct {
	Name     Name2  `json:"Nameku"`
	FullName string `json:"FullName"`
	Ages     string `json:"Ages"`
}

func main() {
	errKosong := error{}

	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		var aku []interface{}
		aku = append(aku, Person{})
		aku = append(aku, Name{})
		data, err := ShouldBindJsonWithCamelCase(c, aku)
		if err != errKosong {
			c.JSON(400, err)
			return
		}
		c.JSON(200, data)
	})
	r.Run()

}

type error struct {
	Type     string
	Location string
	Parent   string
}

func ShouldBindJsonWithCamelCase(c *gin.Context, model []interface{}) (interface{}, error) {
	errKosong := error{}
	//Declare Local Variable
	var dataMap map[string]interface{} // for mapping data

	//Get data from web
	data, errReadGin := ioutil.ReadAll(c.Request.Body)
	if errReadGin != nil {
		errNewReadGin := error{Type: "Error to Gin", Location: "ioutil.ReadAll(c.Request.Body)", Parent: "none"}
		return nil, errNewReadGin
	}
	// Marshaling to MapString
	errUnmarshalDataToMap := json.Unmarshal([]byte(string(data)), &dataMap)
	if errUnmarshalDataToMap != nil {
		errNewUnmarshalDataToMap := error{Type: "Error to Map String", Location: "json.Unmarshal([]byte(string(data)), &dataMap)", Parent: "none"}
		return nil, errNewUnmarshalDataToMap
	}
	errCamelCase := CamelCaseChecker(dataMap, model)
	if errCamelCase != errKosong {
		return nil, errCamelCase
	}
	errSouldBind := json.Unmarshal([]byte(string(data)), &model[0])
	if errSouldBind != nil {
		errNewSouldBind := error{Type: "Error to Bind Data", Location: "json.Unmarshal([]byte(string(data)), &model[0])", Parent: "none"}
		return nil, errNewSouldBind
	}
	return model[0], errKosong
}

var parent []string

func CamelCaseChecker(data interface{}, model []interface{}) error {
	errKosong := error{}
	var tempData map[string]interface{}
	jsonData, _ := json.Marshal(data)
	json.Unmarshal(jsonData, &tempData)
	boolku := true
	rt := reflect.TypeOf(model[0])
	dt2 := reflect.ValueOf(tempData)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("json"), ",")[0]
		for _, key2 := range dt2.MapKeys() {
			if v != key2.String() {
				if strings.EqualFold(v, key2.String()) {
					if len(parent) < 1 {
						err := error{Type: "Error Camel Case", Location: v, Parent: "main"}
						return err
					} else {
						err := error{Type: "Error Camel Case", Location: v, Parent: parent[len(parent)-1]}
						parent = nil
						return err
					}
				}
				boolku = false
			} else if v == key2.String() {
				boolku = true
				break
			}
		}
		if !boolku {
			isOmitEmpty := strings.Contains(f.Tag.Get("json"), "omitempty")
			if !isOmitEmpty {
				if len(parent) < 1 {
					err := error{Type: "Error Camel Case", Location: v, Parent: "main"}
					return err
				} else {
					err := error{Type: "Error Camel Case", Location: v, Parent: parent[len(parent)-1]}
					parent = nil
					return err
				}
			}
		} else if strings.Contains(f.Type.String(), ".") {
			parent = append(parent, v)
			for z := 1; z < len(model); z++ {
				temp := model[z:]
				err := CamelCaseChecker(tempData[v], temp)
				if err != errKosong {
					return err
				}
			}
			parent = RemoveIndex(parent, v)
		}
	}
	return error{}
}
func RemoveIndex(s []string, index string) []string {
	if len(s) >= 1 {
		for i := 0; i < len(s); i++ {
			if s[i] == index {
				return append(s[:i], s[i+1:]...)
			}
		}
	}
	return []string{}
}
