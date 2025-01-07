/*
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Person struct {
	Name      string
	Age       int
	Salary    float64
	Education string
}

/*
type People struct {
	Persons []Person
}

func main() {

	/*
		file, err := os.ReadFile("./people.json")
		if err != nil {
			log.Fatal(err)
		}
		var people []People
		err = json.Unmarshal(file, &people)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(people)
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Person struct {
	Name      string
	Age       int
	Salary    float64
	Education string
}

func main() {

	file, err := os.ReadFile("./people.json")
	if err != nil {
		log.Fatal(err)
	}
	pers := []Person{}
	err = json.Unmarshal(file, &pers)
	if err != nil {
		log.Fatal(err)
	}

	var NameYoungest string = pers[0].Name
	var NameOldest string = pers[0].Name
	var YoungestAge int = pers[0].Age
	var OldestAge int = pers[0].Age
	var NameHighestSalary string
	var HighestSalary float64 = pers[0].Salary
	var NameLowestSalary string
	var LowestSalary float64 = pers[0].Salary
	freq := make(map[string]int)

	totalAge := 0
	var totalSalary float64 = 0
	for i, p := range pers {

		log.Println(i, p.Name, p.Age)
		totalAge += p.Age
		totalSalary += p.Salary
		if p.Age > OldestAge {
			OldestAge = p.Age
			NameOldest = p.Name
		}
		if p.Age < YoungestAge {
			YoungestAge = p.Age
			NameYoungest = p.Name
		}
		if p.Salary > HighestSalary {
			HighestSalary = p.Salary
			NameHighestSalary = p.Name
		}
		if p.Salary < LowestSalary {
			LowestSalary = p.Salary
			NameLowestSalary = p.Name
		}
		freq[p.Education]++
	}

	averageAge := totalAge / len(pers)
	averageSalary := totalSalary / float64(len(pers))
	fmt.Println("Average age: ", averageAge)
	fmt.Println("Average salary: ", averageSalary)
	fmt.Println("Highest salary: ", NameHighestSalary, HighestSalary)
	fmt.Println("Lowest salary: ", NameLowestSalary, LowestSalary)
	fmt.Println("Oldest person: ", NameOldest, OldestAge)
	fmt.Println("Youngest person: ", NameYoungest, YoungestAge)
	fmt.Println("Counts of people per Education:\n")
	fmt.Println(freq)

	data := map[string]interface{}{
		"Average Age":                                      averageAge,
		"Name(s) of the Youngest Person(s)":                NameYoungest,
		"Name(s) of the Oldest Person(s)":                  NameOldest,
		"Average Salary":                                   averageSalary,
		"Name(s) of the Person(s) with the Highest Salary": NameHighestSalary,
		"Name(s) of the Person(s) with the Lowest Salary":  NameLowestSalary,
		"Counts of People by Education Level":              freq,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	fmt.Printf("json data: %s\n", jsonData)
	ioutil.WriteFile("Stats.json", jsonData, os.ModePerm)
}
