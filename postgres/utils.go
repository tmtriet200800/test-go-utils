package pkgPostgres

import (
	"fmt"
	"reflect"
	"strings"
)

func GetKeysString(data interface{}, tag string, filterNull bool) string {
	val := reflect.ValueOf(data)
	keynames := []string{}

	startIdx := 1
	
	for i:=0; i<val.NumField();i++{
		valid := val.Field(i).FieldByName("Valid").Bool()

		if filterNull && !valid {continue}

		if tag == "name" {
			keynames = append(keynames, val.Type().Field(i).Name)
		} else if tag == "db-map" {
			keynames = append(keynames, fmt.Sprintf("$%v", startIdx))
			startIdx += 1
		} else {
			keynames = append(keynames, val.Type().Field(i).Tag.Get(tag))
		}
	}

	return strings.Join(keynames[:], ",")
}

func GetValues(data interface{}, filterNull bool) []interface{}{
	val := reflect.ValueOf(data)

	values := []interface{}{}
	
	for i:=0; i<val.NumField();i++{
		valid := val.Field(i).FieldByName("Valid").Bool()

		if filterNull && !valid {continue}

		values = append(values, val.Field(i).Interface())
	}

	return values
}

func GetDBSetStatmentString(data interface{}, filterNull bool) string{
	val := reflect.ValueOf(data)
	setStatements := []string{}

	startIdx := 1;
	
	for i:=0; i<val.NumField();i++{
		valid := val.Field(i).FieldByName("Valid").Bool()
		
		if filterNull && !valid {continue}

		setStatements = append(
			setStatements, 
			fmt.Sprintf(
				"%v=%v",
				val.Type().Field(i).Tag.Get("db"),
				fmt.Sprintf("$%v",startIdx),
			),
		)
		startIdx += 1
	}

	return strings.Join(setStatements[:], ",")
}
