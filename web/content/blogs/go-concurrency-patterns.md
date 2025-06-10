---
title: "Go Concurrency Patterns: Channels, Goroutines, and Real-World Applications"
excerpt: "Explore advanced Go concurrency patterns including worker pools, fan-in/fan-out, and pipeline patterns with practical code examples."
author: "Full Stack Developer"
date: "2024-01-20"
tags: ["golang", "concurrency", "goroutines", "channels", "backend"]
category: "Backend"
read_time: 12
published: true
has_mermaid: false
has_code_blocks: true
---

# Go Concurrency Patterns: Channels, Goroutines, and Real-World Applications

Go's concurrency model, built around goroutines and channels, is one of its most powerful features. Unlike traditional threading models, Go's approach makes concurrent programming more accessible and less error-prone. In this post, we'll explore essential concurrency patterns that every Go developer should master.

## The Foundation: Goroutines and Channels

Before diving into patterns, let's review the basics. Goroutines are lightweight threads managed by the Go runtime, and channels are the pipes that connect them.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Simple goroutine example
    go func() {
        fmt.Println("Hello from goroutine!")
    }()

    // Channel communication
    ch := make(chan string)
    go func() {
        ch <- "Hello from channel!"
    }()

    message := <-ch
    fmt.Println(message)

    time.Sleep(100 * time.Millisecond) // Let goroutine finish
}
```

## Pattern 1: Worker Pool

The worker pool pattern is essential for controlling resource usage and processing jobs concurrently. It's particularly useful when you have many tasks but want to limit concurrent execution.

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID     int
    Data   string
    Result chan string
}

type WorkerPool struct {
    workerCount int
    jobQueue    chan Job
    wg          sync.WaitGroup
}

func NewWorkerPool(workerCount int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobQueue:    make(chan Job, 100), // Buffered channel
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()

    for job := range wp.jobQueue {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)

        // Simulate work
        time.Sleep(100 * time.Millisecond)
        result := fmt.Sprintf("Processed: %s by worker %d", job.Data, id)

        job.Result <- result
        close(job.Result)
    }
}

func (wp *WorkerPool) Submit(job Job) {
    wp.jobQueue <- job
}

func (wp *WorkerPool) Stop() {
    close(wp.jobQueue)
    wp.wg.Wait()
}

func main() {
    pool := NewWorkerPool(3)
    pool.Start()

    // Submit jobs
    for i := 0; i < 10; i++ {
        job := Job{
            ID:     i,
            Data:   fmt.Sprintf("data-%d", i),
            Result: make(chan string, 1),
        }

        pool.Submit(job)

        // Get result
        go func(j Job) {
            result := <-j.Result
            fmt.Printf("Job %d result: %s\n", j.ID, result)
        }(job)
    }

    time.Sleep(2 * time.Second)
    pool.Stop()
}
```

## Pattern 2: Fan-In/Fan-Out

Fan-out distributes work across multiple goroutines, while fan-in collects results from multiple sources. This pattern is excellent for parallel processing and result aggregation.

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Fan-out: distribute work to multiple workers
func fanOut(input <-chan int, workerCount int) []<-chan int {
    workers := make([]<-chan int, workerCount)

    for i := 0; i < workerCount; i++ {
        worker := make(chan int)
        workers[i] = worker

        go func(w chan<- int) {
            defer close(w)
            for data := range input {
                // Simulate processing time
                time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
                w <- data * data // Square the number
            }
        }(worker)
    }

    return workers
}

// Fan-in: collect results from multiple workers
func fanIn(workers []<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup

    for _, worker := range workers {
        wg.Add(1)
        go func(w <-chan int) {
            defer wg.Done()
            for result := range w {
                output <- result
            }
        }(worker)
    }

    go func() {
        wg.Wait()
        close(output)
    }()

    return output
}

func main() {
    // Create input channel
    input := make(chan int)

    // Start the pipeline
    workers := fanOut(input, 3)
    results := fanIn(workers)

    // Send data
    go func() {
        defer close(input)
        for i := 1; i <= 10; i++ {
            input <- i
        }
    }()

    // Collect results
    for result := range results {
        fmt.Printf("Result: %d\n", result)
    }
}
```

## Pattern 3: Pipeline

Pipelines allow you to chain operations where each stage processes data and passes it to the next stage. This pattern promotes clean separation of concerns and efficient data processing.

```go
package main

import (
    "fmt"
    "strings"
)

// Stage 1: Generate data
func generator(data []string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for _, item := range data {
            out <- item
        }
    }()
    return out
}

// Stage 2: Transform data (uppercase)
func transformer(input <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for item := range input {
            out <- strings.ToUpper(item)
        }
    }()
    return out
}

// Stage 3: Filter data (only items with length > 3)
func filter(input <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for item := range input {
            if len(item) > 3 {
                out <- item
            }
        }
    }()
    return out
}

// Stage 4: Add prefix
func prefixer(input <-chan string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for item := range input {
            out <- fmt.Sprintf("PROCESSED: %s", item)
        }
    }()
    return out
}

func main() {
    data := []string{"go", "rust", "python", "java", "c++", "js"}

    // Build pipeline
    pipeline := prefixer(filter(transformer(generator(data))))

    // Process results
    for result := range pipeline {
        fmt.Println(result)
    }
}
```

## Pattern 4: Timeout and Cancellation

Handling timeouts and cancellation is crucial for building robust applications. Go's `context` package provides excellent tools for this.

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func longRunningTask(ctx context.Context, id int) <-chan string {
    result := make(chan string, 1)

    go func() {
        defer close(result)

        // Simulate long-running work
        select {
        case <-time.After(2 * time.Second):
            result <- fmt.Sprintf("Task %d completed", id)
        case <-ctx.Done():
            result <- fmt.Sprintf("Task %d cancelled: %v", id, ctx.Err())
        }
    }()

    return result
}

func main() {
    // Example 1: Timeout
    fmt.Println("=== Timeout Example ===")
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()

    result := longRunningTask(ctx, 1)
    fmt.Println(<-result)

    // Example 2: Manual cancellation
    fmt.Println("\n=== Manual Cancellation Example ===")
    ctx2, cancel2 := context.WithCancel(context.Background())

    result2 := longRunningTask(ctx2, 2)

    // Cancel after 500ms
    go func() {
        time.Sleep(500 * time.Millisecond)
        cancel2()
    }()

    fmt.Println(<-result2)

    // Example 3: Multiple tasks with shared context
    fmt.Println("\n=== Multiple Tasks Example ===")
    ctx3, cancel3 := context.WithTimeout(context.Background(), 1500*time.Millisecond)
    defer cancel3()

    tasks := make([]<-chan string, 3)
    for i := 0; i < 3; i++ {
        tasks[i] = longRunningTask(ctx3, i+3)
    }

    for i, task := range tasks {
        fmt.Printf("Task %d: %s\n", i+3, <-task)
    }
}
```

## Pattern 5: Rate Limiting

Rate limiting is essential for controlling resource usage and preventing system overload. Here's a token bucket implementation:

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type RateLimiter struct {
    tokens chan struct{}
    ticker *time.Ticker
    mu     sync.Mutex
}

func NewRateLimiter(rps int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rps),
        ticker: time.NewTicker(time.Second / time.Duration(rps)),
    }

    // Fill initial tokens
    for i := 0; i < rps; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill tokens
    go func() {
        for range rl.ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
                // Channel full, skip
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

func (rl *RateLimiter) Stop() {
    rl.ticker.Stop()
}

func simulateAPICall(id int) {
    fmt.Printf("API call %d at %s\n", id, time.Now().Format("15:04:05.000"))
    time.Sleep(50 * time.Millisecond) // Simulate processing
}

func main() {
    // Allow 2 requests per second
    limiter := NewRateLimiter(2)
    defer limiter.Stop()

    var wg sync.WaitGroup

    // Make 10 API calls
    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            limiter.Wait() // This will block if rate limit exceeded
            simulateAPICall(id)
        }(i)
    }

    wg.Wait()
}
```

## Real-World Example: Web Scraper

Let's combine these patterns in a practical web scraper example:

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "sync"
    "time"
)

type ScrapeJob struct {
    URL    string
    Result chan ScrapeResult
}

type ScrapeResult struct {
    URL        string
    StatusCode int
    Error      error
    Content    string
}

type WebScraper struct {
    client      *http.Client
    rateLimiter *RateLimiter
    workerPool  *WorkerPool
}

func NewWebScraper(workers int, rps int) *WebScraper {
    return &WebScraper{
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
        rateLimiter: NewRateLimiter(rps),
        workerPool:  NewWorkerPool(workers),
    }
}

func (ws *WebScraper) scrapeURL(ctx context.Context, url string) ScrapeResult {
    // Rate limiting
    ws.rateLimiter.Wait()

    // Create request with context
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return ScrapeResult{URL: url, Error: err}
    }

    resp, err := ws.client.Do(req)
    if err != nil {
        return ScrapeResult{URL: url, Error: err}
    }
    defer resp.Body.Close()

    return ScrapeResult{
        URL:        url,
        StatusCode: resp.StatusCode,
        Content:    fmt.Sprintf("Content length: %d", resp.ContentLength),
    }
}

func (ws *WebScraper) ScrapeURLs(ctx context.Context, urls []string) <-chan ScrapeResult {
    results := make(chan ScrapeResult, len(urls))
    var wg sync.WaitGroup

    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            result := ws.scrapeURL(ctx, u)
            results <- result
        }(url)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}

func main() {
    urls := []string{
        "https://httpbin.org/delay/1",
        "https://httpbin.org/delay/2",
        "https://httpbin.org/status/200",
        "https://httpbin.org/status/404",
        "https://httpbin.org/json",
    }

    scraper := NewWebScraper(3, 2) // 3 workers, 2 requests per second
    defer scraper.rateLimiter.Stop()

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    fmt.Println("Starting web scraper...")
    start := time.Now()

    results := scraper.ScrapeURLs(ctx, urls)

    for result := range results {
        if result.Error != nil {
            fmt.Printf("Error scraping %s: %v\n", result.URL, result.Error)
        } else {
            fmt.Printf("Success: %s (Status: %d) - %s\n",
                result.URL, result.StatusCode, result.Content)
        }
    }

    fmt.Printf("Completed in %v\n", time.Since(start))
}
```

## Best Practices

1. **Always close channels**: Use `defer close(ch)` when appropriate
2. **Use buffered channels wisely**: They can prevent goroutine leaks but use memory
3. **Handle context cancellation**: Always check `ctx.Done()` in long-running operations
4. **Avoid sharing memory**: Communicate through channels instead of shared variables
5. **Use sync.WaitGroup for coordination**: When you need to wait for multiple goroutines
6. **Set timeouts**: Never let operations run indefinitely

## Conclusion

Go's concurrency patterns provide powerful tools for building efficient, scalable applications. By mastering these patterns - worker pools, fan-in/fan-out, pipelines, timeouts, and rate limiting - you can write robust concurrent code that handles real-world challenges effectively.

Remember: "Don't communicate by sharing memory; share memory by communicating." This philosophy, embodied in Go's channel-based concurrency model, leads to cleaner, more maintainable code.
