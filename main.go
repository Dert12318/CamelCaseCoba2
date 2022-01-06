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
	Ages     string `json:"Ages"`
}
type Person2 struct {
	Name     Name2  `json:"Nameku"`
	FullName string `json:"FullName"`
	Ages     string `json:"Ages"`
}

func main() {
	r := gin.Default()
	r.POST("/ping", func(c *gin.Context) {
		dataku, errlala := ioutil.ReadAll(c.Request.Body)
		var req2 map[string]interface{}
		fmt.Println(errlala)
		finalData := string(dataku)
		errJson := json.Unmarshal([]byte(string(finalData)), &req2)
		fmt.Println(errJson)
		err := coba(req2, Person{})
		fmt.Println(err, "error guys")
		c.JSON(200, err.Error())
	})
	r.Run()

}

func coba(data interface{}, model interface{}) error {
	var modelMap map[string]interface{}
	var dataMap map[string]interface{}
	boolku := true
	jsonModel, _ := json.Marshal(model)
	json.Unmarshal(jsonModel, &modelMap)
	fmt.Println("Model Map :", modelMap)
	jsonModelMap, _ := json.Marshal(data)
	json.Unmarshal(jsonModelMap, &dataMap)
	fmt.Println("Data Map :", dataMap)
	dt := reflect.ValueOf(modelMap)
	dt2 := reflect.ValueOf(data)
	for _, key := range dt.MapKeys() {
		strct := dt.MapIndex(key)
		res := fmt.Sprintf("%s", strct.Interface())
		for _, key2 := range dt2.MapKeys() {
			strct2 := dt2.MapIndex(key2)
			res2 := fmt.Sprintf("%s", strct2.Interface())
			fmt.Println(res, res2, "lop")
			if key.String() != key2.String() {
				fmt.Println(key.String(), key2.String())
				boolku = false
			} else if key.String() == key2.String() {
				boolku = true
				fmt.Println("masuk sini guys")
				break
			}
		}

		if !boolku {
			return errors.New(key.String())
		} else if strings.Contains(res, "map[") {
			fmt.Println("masuk sini satu")
			err := coba(dataMap[key.String()], modelMap[key.String()])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
