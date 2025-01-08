package main

import "fmt"

type Vehicule struct {
	Make  string
	Model string
	Year  int
	Level int
}

type Insurable interface {
	CalculateInsurance() float64
}

type Printable interface {
	Details()
}

type Car struct {
	Vehicule
	NumberOfDoors int64
}

func (c Car) CalculateInsurance() float64 {
	return float64(c.NumberOfDoors) * float64(c.Level) * 100
}

func (c Car) Details() {
	fmt.Println(c.Make, c.Model, c.Year, c.Level)
}

type Truck struct {
	Vehicule
	PayloadCapacity int
}

func (t Truck) CalculateInsurance() float64 {
	return float64(t.PayloadCapacity) * float64(t.Level) * 100
}

func (t Truck) Details() {
	fmt.Println(t.Make, t.Model, t.Year, t.Level)
}

func PrintAll(p []Printable) {
	for _, i := range p {
		i.Details()
	}
}

func main() {
	c1 := Car{
		Vehicule: Vehicule{
			Make:  "YAMAHA",
			Model: "C20",
			Year:  2004,
			Level: 1,
		},
		NumberOfDoors: 4,
	}
	t1 := Truck{
		Vehicule: Vehicule{
			Make:  "YAMAHA",
			Model: "C20",
			Year:  2004,
			Level: 1,
		},
		PayloadCapacity: 5,
	}
	c2 := Car{
		Vehicule: Vehicule{
			Make:  "BMW",
			Model: "A2",
			Year:  2007,
			Level: 3,
		},
		NumberOfDoors: 4,
	}
	t2 := Truck{
		Vehicule: Vehicule{
			Make:  "YAMAHA",
			Model: "C20",
			Year:  2004,
			Level: 1,
		},
		PayloadCapacity: 4,
	}
	//c.Details()
	//fmt.Println("Insurance: ", c.CalculateInsurance())

	p := []Printable{c1, t1, c2, t2}
	PrintAll(p)

}
