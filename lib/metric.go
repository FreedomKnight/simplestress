package simplestress

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"time"
)

type Metric struct {
    MaxLatency time.Duration `json:"max_latency"`
    MinLatency time.Duration `json:"min_latency"`
    AverageLatency time.Duration `json:"average_latency"`
    Requests int64 `json:"requests"`
    Errors int64 `json:"errors"`
    Success int64 `json:"success"`
    TotalLatency time.Duration `json:"total_latency"`
}

func NewMetric() *Metric {
    return &Metric{
        MaxLatency: 0,
        MinLatency: math.MaxInt64 * time.Nanosecond,
        AverageLatency: 0,
        Requests: 0,
        Errors: 0,
        Success: 0,
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

func (m *Metric) Save(path string) {
    jsonString, _ := json.MarshalIndent(m, "", "  ")
    ioutil.WriteFile(path, jsonString, 0644)
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

                m.Success = m.Requests - m.Errors

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
