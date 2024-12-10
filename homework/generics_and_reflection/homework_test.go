package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(person Person) string {
	t := reflect.TypeOf(person)
	v := reflect.ValueOf(person)
	res := make([]string, 0, t.NumField())
	for i := range t.NumField() {
		tag, exists := t.Field(i).Tag.Lookup("properties")
		if !exists {
			continue
		}
		if v.Field(i).IsZero() && strings.Contains(tag, "omitempty") {
			continue
		}
		tagName := strings.Replace(tag, ",omitempty", "", -1)
		res = append(res, fmt.Sprintf("%s=%v", tagName, v.Field(i)))
	}
	fmt.Println(strings.Join(res, "\n"))
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
