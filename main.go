package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type Name struct {
	FirstName  string `json:"FirstName"`
	LastName   string `json:"LastName"`
	MiddleName string `json:"MiddleName"`
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
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		data, err := ShouldBindJsonWithCamelCase(c, Person{})
		if err != nil {
			c.JSON(400, "error")
			fmt.Println("error gys", err)
			return
		}
		fmt.Println(data, "ini data guys")
		c.JSON(200, "Masuk Gan")
	})
	r.Run()

}

// type errorKu struct {
// }

func ShouldBindJsonWithCamelCase(c *gin.Context, model interface{}) (interface{}, error) {
	//Declare Local Variable
	var dataMap map[string]interface{} // for mapping data

	//Get data from web
	data, errReadGin := ioutil.ReadAll(c.Request.Body)
	if errReadGin != nil {
		fmt.Println("1")
		return nil, errReadGin
	}
	// Marshaling to MapString
	errUnmarshalDataToMap := json.Unmarshal([]byte(string(data)), &dataMap)
	if errUnmarshalDataToMap != nil {
		fmt.Println("2")
		return nil, errUnmarshalDataToMap
	}
	errCamelCase := CamelCaseChecker(dataMap, model)
	if errCamelCase != nil {
		fmt.Println("3")
		return nil, errCamelCase
	}
	return nil, nil

}
func CamelCaseChecker(data interface{}, model interface{}) error {
	var tempData map[string]interface{}
	jsonModel, _ := json.Marshal(model)
	json.Unmarshal(jsonModel, &tempData)
	boolku := true
	fmt.Println("Model Map :", model)
	fmt.Println("Data Map :", tempData)
	rt := reflect.TypeOf(model)
	dt2 := reflect.ValueOf(tempData)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get("json"), ",")[0] // use split to ignore tag "options"
		for _, key2 := range dt2.MapKeys() {
			strct2 := dt2.MapIndex(key2)
			res2 := fmt.Sprintf("%s", strct2.Interface())
			fmt.Println(f, res2, "lop")
			if v != key2.String() {
				fmt.Println("gak cocok")
				boolku = false
			} else if v == key2.String() {
				fmt.Println("cocok")
				boolku = true
				break
			}
		}
		if !boolku {
			isOmitEmpty := strings.Contains(f.Tag.Get("json"), "omitempty")
			if !isOmitEmpty {
				return errors.New(v)
			}
		} else if strings.Contains(f.Type.String(), ".") {
			fmt.Println(tempData[v], Name{}, "rekursif")
			err := CamelCaseChecker(tempData[v], Name{})
			fmt.Println(err, "errKUalat")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
