package main

import (
	"flag"
	"time"

	lib "github.com/FreedomKnight/simplestress/lib"
)

var (
    addr = flag.String("addr", "localhost:50051", "the address to connect to")
    frequency = flag.Int64("frequency", 1, "the frequency of requests per second")
    concurrent = flag.Int64("concurrent", 2, "the number of concurrent requests")
    reportPath = flag.String("report-path", "report.json", "the path to the report file")
    runtime = flag.Int64("runtime", 10, "the number of seconds to run")
)

func main() {
    flag.Parse()

    r := lib.NewRunner(
        lib.Address(*addr),
        lib.Frequency(*frequency),
        lib.Concurrent(*concurrent),
        lib.Runtime(time.Duration(*runtime) * time.Second),
    )
    results := r.Run()

    stopMetric := make(chan int)
    metric := lib.NewMetric()
    metric.Watch(results, stopMetric)
    metric.Print()
    metric.Save(*reportPath)
}

