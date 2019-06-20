package main

import "fmt"

type Foo struct {
	aNum int
}

func main() {
	a := make([][]Foo, 5)

	for r := range a {
		a[r] = make([]Foo, 2)
		for i := range a[r] {
			a[r][i] = Foo{r + i}
		}
	}

	for r := range a {
		for i := range a[r] {
			fmt.Printf("%d/%d: %t\n", r, i, a[r][i])
		}
	}
}
