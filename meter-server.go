package main

import (
	"fmt"
)

func main() {

	c := Collection{1, "Electra (laag)", ""};
	m := Measurement{Collection: c, Id:1, Value:200};

	fmt.Print(m);
}

