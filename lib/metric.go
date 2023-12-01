package simplestress

import (
	"log"
	"math"
	"time"
)

type Metric struct {
    MaxLatency time.Duration
    MinLatency time.Duration
    AverageLatency time.Duration
    Requests int64
    Errors int64
    TotalLatency time.Duration
}

func NewMetric() *Metric {
    return &Metric{
        MaxLatency: 0,
        MinLatency: math.MaxInt64 * time.Nanosecond,
        AverageLatency: 0,
        Requests: 0,
        Errors: 0,
        TotalLatency: 0,
    }
}

func (m *Metric) Print() {
    log.Printf("Requests: %v\n", m.Requests)
    log.Printf("Errors: %v\n", m.Errors)
    log.Printf("Success: %v\n", m.Requests - m.Errors)
    log.Printf("MinLatency: %v\n", m.MinLatency)
    log.Printf("MaxLatency: %v\n", m.MaxLatency)
    log.Printf("TotalLatency: %v\n", m.TotalLatency)
    log.Printf("AverageLatency: %v\n", m.AverageLatency)
}

func (m *Metric) Watch(results chan *Result, stop chan int) {
    for {
        select {
            case <-stop:
                return
            case result, ok := <-results:
                if !ok {
                    return
                }

                m.Requests++
                if result.err != nil{
                    m.Errors++
                }
                if m.MaxLatency < result.latency {
                    m.MaxLatency = result.latency
                }

                if  m.MinLatency > result.latency {
                    m.MinLatency = result.latency
                }

                m.TotalLatency += result.latency
                m.AverageLatency = m.TotalLatency / time.Duration(m.Requests)
                
        }
    }
}
