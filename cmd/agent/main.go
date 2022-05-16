package main

import (
	"fmt"

	"github.com/MustCo/Mon_go/internal/utils/my_metrics"
)

func main() {
	m := new(my_metrics.Metrics)
	m.Init()
	m.Poll()
	fmt.Println(m)
}
