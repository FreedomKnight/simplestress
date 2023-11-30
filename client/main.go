package main

import (
	"flag"

	lib "github.com/FreedomKnight/simplestress/lib"
)

var (
    addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
    flag.Parse()
    r := lib.NewRunner(
        lib.Address(*addr),
        lib.Frequency(2),
        lib.Concurrent(1),
    )
    results := r.Run()

    stopMetric := make(chan int)
    metric := lib.NewMetric()
    metric.Watch(results, stopMetric)
    metric.Print()
}

