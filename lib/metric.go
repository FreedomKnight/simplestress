package simplestress

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

func (m *Metric) Plot(path string) {
    bar := charts.NewBar()
    bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
        Title:    "Success/Errors",
    }))

    items := make([]opts.BarData, 0)
    items = append(items, opts.BarData{Value: m.Success})
    items = append(items, opts.BarData{Value: m.Errors})

    bar.SetXAxis([]string{"Success", "Errors"}).
        AddSeries("Success/Errors", items)
    f, _ := os.Create(path)
    err := bar.Render(f)
    if err != nil {
        log.Printf("Failed to save histogram file: %v\n", err)
    }
}

func (m *Metric) Save(path string) {
    jsonString, _ := json.MarshalIndent(m, "", "  ")
    err := ioutil.WriteFile(path, jsonString, 0644)
    if err != nil {
        log.Printf("Failed to save report file: %v\n", err)
    }
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
