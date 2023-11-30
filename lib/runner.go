package simplestress

import (
	"context"
	"log"
	"sync"
	"time"

	pb "github.com/FreedomKnight/simplestress/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var DefaultTimeout = 1 * time.Second
var DefaultFrequency = 1
var DefaultAddress = "localhost:50051"
var DefaultRuntime = 10 * time.Second

type Runner struct {
    timeout time.Duration
    runtime time.Duration
    startAt time.Time
    stop chan int
    once sync.Once
    concurrent int64
    frequency int64
    address string
}

type Result struct {
    latency time.Duration
    err error
}

func NewRunner(options ...func(*Runner)) *Runner {
    r := &Runner{
        timeout: DefaultTimeout,
        address: DefaultAddress,
        runtime: DefaultRuntime,
        stop: make(chan int),
    }

    for _, option := range options {
        option(r)
    }

    return r
}

func Frequency(frequency int64) func(*Runner) {
    return func (r *Runner) { r.frequency = frequency }
}

func Concurrent(concurrent int64) func(*Runner) {
    return func (r *Runner) { r.concurrent = concurrent }
}

func Address(address string) func(*Runner) {
    return func (r *Runner) { r.address = address }
}

func Runtime(runtime time.Duration) func(*Runner) {
    return func (r *Runner) { r.runtime = runtime }
}

func (r *Runner) ping() *Result {
    now := time.Now()
    conn, err := grpc.Dial(r.address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    defer conn.Close()
    if ( err != nil ) {
        log.Println(err.Error())
    }

    c := pb.NewPaddleClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
    defer cancel()
    pong, err := c.Serve(ctx, &pb.Ping{Message: "Hello World!"})
    if (err != nil) {
        log.Println(err.Error())
    }

    log.Printf("Received Pong %s\n", pong.GetMessage())
    return &Result{time.Since(now), err}
}

func (r *Runner) work(wg *sync.WaitGroup, ticks chan int, results chan *Result) {
    defer wg.Done()
    // ping if get tick signal
    for range ticks {
        results <- r.ping()
    }
}

func (r *Runner) calculateWaitingTime(elapsed time.Duration, hits int64) (time.Duration) {
    expectedHits := r.frequency * int64(elapsed / time.Second)

    // run immediately if hits is less than expected hits
    if hits < expectedHits {
        return 0
    }

    // nanoseconds per hit
    interval := int64(time.Second.Nanoseconds()) / r.frequency

    // evaluate next hit should be wait
    return time.Duration(interval * (hits + 1)) - elapsed
}

func (r *Runner) Stop() {
    select {
        case <-r.stop:
            return
        default:
            r.once.Do(func() {
                close(r.stop)
            })
    }
}

func (r *Runner) Run() chan *Result {
    wg := &sync.WaitGroup{}

    r.startAt = time.Now()

    // build worker and let them wait for tick signal
    ticks := make(chan int)
    results := make(chan *Result, 10)
    for i := int64(0); i < r.concurrent; i++ {
        wg.Add(1)
        go r.work(wg, ticks, results)
    }

    // dont't block main thread
    // it will return return result channel to metrics
    go func() {
        defer func() {
            defer close(ticks)
            defer wg.Wait()
            defer close(results)
            defer r.Stop()
        }()

        hits := int64(0)
        for {
            elaspse := time.Since(r.startAt)
            // check if time is up
            if time.Since(r.startAt) > r.runtime {
                break
            }

            waitingNanoSeconds := r.calculateWaitingTime(elaspse, hits)
            time.Sleep(waitingNanoSeconds)

            // send tick signal to worker
            select  {
                case <-r.stop:
                    return
                case ticks <- 1:
                    hits++
            }
        }
    }()
    return results
}
