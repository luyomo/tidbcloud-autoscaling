package main

import (
    "fmt"
    "time"
    "os"
    "context"

    "github.com/prometheus/client_golang/api"
    "github.com/prometheus/common/model"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

func main() {
    client, err := api.NewClient(api.Config{
		Address: "http://localhost:9090",
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}


    ticker := time.NewTicker(time.Minute)
    for {
       select {
           case t := <-ticker.C:
               fmt.Println("Tick at", t)

	           v1api := v1.NewAPI(client)
	           ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	           defer cancel()
	           r := v1.Range{
	           	   Start: time.Now().Add(-time.Minute * 5),
	           	   End:   time.Now(),
	           	   Step:  time.Minute,
	           }
	           result, warnings, err := v1api.QueryRange(ctx, "rate(tidbcloud_node_cpu_seconds_total{cluster_name=\"scalingtest\", component=\"tidb\"}[5m])", r, v1.WithTimeout(10*time.Second))
	           if err != nil {
	               fmt.Printf("Error querying Prometheus: %v\n", err)
	               os.Exit(1)
	           }
	           if len(warnings) > 0 {
	           	fmt.Printf("Warnings: %v\n", warnings)
	           }
	           // fmt.Printf("Result:\n%v\n", len(result.(Matrix)))
               for _, val := range result.(model.Matrix)[0].Values {
                   fmt.Printf("Result: %s:%s\n", val.Timestamp, val.Value)
                   if val.Value < 1.6 {
                       continue
                   }
               }
       }
    }
}
