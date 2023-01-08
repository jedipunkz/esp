package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jedipunkz/esp/cloudwatch"
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

	containersMetadata, err := client.RetriveContainersMetadata(ctx)
	if err != nil {
		log.Printf("error retrieving task stats metadata: %s", err)
	}

	containerName := os.Getenv("CONTAINER_NAME")
	if containerName == "" {
		log.Fatal("error retrieving container name from environment variable")
	}

	for {
		for _, container := range taskMetadata.Containers {
			if container.Name == containerName {
				s := containersMetadata[container.DockerID]
				if &s == nil {
					log.Printf("Could not find stats for container %s", container.DockerID)
					continue
				}

				cpuUsage := (float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)) /
					(float64(s.CPUStats.SystemCPUUsage) - float64(s.PreCPUStats.SystemCPUUsage)) *
					float64(s.CPUStats.OnlineCPUs) * 100

				_, err = cloudwatch.PutMetricData(cpuUsage)
				if err != nil {
					log.Printf("Error putting metric data: %s", err)
				}

				log.Printf("Container Name: %s, CPU Usage: %f", container.Name, cpuUsage)
			}
		}
		time.Sleep(time.Second * 1)
	}
}
