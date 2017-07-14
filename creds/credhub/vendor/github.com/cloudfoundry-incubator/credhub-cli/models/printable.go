package models

import (
	"fmt"
)

type Printable interface {
	ToYaml() string
	ToJson() string
}

func Println(p Printable, asJson bool) {
	if asJson {
		fmt.Println(p.ToJson())
	} else {
		fmt.Println(p.ToYaml())
	}
}
