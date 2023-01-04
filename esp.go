package main

import (
	"context"
	"log"

	"github.com/jedipunkz/esp/ecs"
)

func main() {
	ctx := context.Background()
	client, err := ecs.NewClientToMetadataEndpoint()
	if err != nil {
		log.Printf("error creating client: %s", err)
	}

	taskMetadata, err := client.RetriveTaskMetadata(ctx)
	if err != nil {
		log.Printf("error retrieving task metadata: %s", err)
	}

	statsMetadata, err := client.RetriveStatsMetadata(ctx)
	if err != nil {
		log.Printf("error retrieving task stats metadata: %s", err)
	}

	for _, container := range taskMetadata.Containers {
		s := statsMetadata[container.DockerID]
		if &s == nil {
			log.Printf("Could not find stats for container %s", container.DockerID)
			continue
		}

		log.Printf("Total CPU Usage: %d", s.CPUStats.CPUUsage.TotalUsage)
		log.Printf("CPU Usage: %f", (float64(s.CPUStats.CPUUsage.TotalUsage)-float64(s.PreCPUStats.CPUUsage.TotalUsage))/
			(float64(s.CPUStats.SystemCPUUsage)-float64(s.PreCPUStats.SystemCPUUsage))*
			float64(s.CPUStats.OnlineCPUs)*100)
	}
}
