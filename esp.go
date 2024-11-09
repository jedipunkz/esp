package main

import (
	"context"
	"os"
	"time"

	"github.com/jedipunkz/esp/cloudwatch"
	"github.com/jedipunkz/esp/ecs"
	log "github.com/sirupsen/logrus"
)

func main() {
	for {
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

		var totalCPUUsage float64
		var containerCount int

		for _, container := range taskMetadata.Containers {
			if container.Name == containerName {
				s := containersMetadata[container.DockerID]

				cpuUsage := ((float64(s.CPUStats.CPUUsage.TotalUsage) - float64(s.PreCPUStats.CPUUsage.TotalUsage)) /
					(float64(s.CPUStats.SystemCPUUsage) - float64(s.PreCPUStats.SystemCPUUsage))) *
					float64(s.CPUStats.OnlineCPUs) * 100

				totalCPUUsage += cpuUsage
				containerCount++

				log.SetFormatter(&log.JSONFormatter{})
				log.WithFields(log.Fields{
					"CPUStats.CPUStats.TotalUsage":    s.CPUStats.CPUUsage.TotalUsage,
					"PreCPUStats.CPUUsage.TotalUsage": s.PreCPUStats.CPUUsage.TotalUsage,
					"CPUStats.SystemCPUUsage":         s.CPUStats.SystemCPUUsage,
					"PreCPUStats.SystemCPUUsage":      s.PreCPUStats.SystemCPUUsage,
					"CPUStats.OnlineCPUs":             s.CPUStats.OnlineCPUs,
					"CPUUsage":                        cpuUsage,
					"ContainerName":                   container.Name,
					"Namespace":                       os.Getenv("NAMESPACE"),
					"REGION":                          os.Getenv("REGION"),
				}).Info("ESP Stats")
			}
		}

		if containerCount > 0 {
			averageCPUUsage := totalCPUUsage / float64(containerCount)

			_, err = cloudwatch.PutMetricData(averageCPUUsage)
			if err != nil {
				log.Printf("Error putting metric data: %s", err)
			}
		}

		time.Sleep(time.Second * 1)
	}
}
