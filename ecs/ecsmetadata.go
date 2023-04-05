package ecs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	ecsMetadataUriEnvV4 = "ECS_CONTAINER_METADATA_URI_V4"
)

type Client struct {
	HTTPClient *http.Client
	endpoint   string
}

type TaskMetadata struct {
	Cluster          string `json:"Cluster"`
	TaskARN          string `json:"TaskARN"`
	Family           string `json:"Family"`
	Revision         string `json:"Revision"`
	DesiredStatus    string `json:"DesiredStatus"`
	KnownStatus      string `json:"KnownStatus"`
	AvailabilityZone string `json:"AvailabilityZone"`
	LaunchType       string `json:"LaunchType"`
	Containers       []struct {
		DockerID      string            `json:"DockerId"`
		Name          string            `json:"Name"`
		DockerName    string            `json:"DockerName"`
		Image         string            `json:"Image"`
		ImageID       string            `json:"ImageID"`
		Labels        map[string]string `json:"Labels"`
		DesiredStatus string            `json:"DesiredStatus"`
		KnownStatus   string            `json:"KnownStatus"`
		Type          string            `json:"Type"`
		ContainerARN  string            `json:"ContainerARN"`
	} `json:"Containers"`
}

type ContainerMetadata struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			PerCPUUsage       []int `json:"percpu_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int `json:"system_cpu_usage"`
		OnlineCPUs     int `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"cpu_stats"`
	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage        int   `json:"total_usage"`
			PerCPUUsage       []int `json:"percpu_usage"`
			UsageInKernelmode int   `json:"usage_in_kernelmode"`
			UsageInUsermode   int   `json:"usage_in_usermode"`
		} `json:"cpu_usage"`
		SystemCPUUsage int `json:"system_cpu_usage"`
		OnlineCPUs     int `json:"online_cpus"`
		ThrottlingData struct {
			Periods          int `json:"periods"`
			ThrottledPeriods int `json:"throttled_periods"`
			ThrottledTime    int `json:"throttled_time"`
		} `json:"throttling_data"`
	} `json:"precpu_stats"`
}

type ErrorNotFound struct{}

func (e *ErrorNotFound) Error() string {
	return "not found"
}

// NewClient retrurns a new ECS client and endpoint
func NewClient(endpoint string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		endpoint:   endpoint,
	}
}

// NewClientToMetadataEndpoint returns a new ECS client and endpoint
func NewClientToMetadataEndpoint() (*Client, error) {
	endpoint := os.Getenv(ecsMetadataUriEnvV4)
	if endpoint == "" {
		return nil, fmt.Errorf("environment variable %s not set", ecsMetadataUriEnvV4)
	}

	_, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %s", err)
	}

	return NewClient(endpoint), nil
}

func (c *Client) request(ctx context.Context, uri string, out interface{}) error {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return &ErrorNotFound{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}

func (c *Client) RetriveTaskMetadata(ctx context.Context) (TaskMetadata, error) {
	var output TaskMetadata
	err := c.request(ctx, c.endpoint+"/task", &output)
	return output, err
}

func (c *Client) RetriveContainersMetadata(ctx context.Context) (map[string]ContainerMetadata, error) {
	output := make(map[string]ContainerMetadata)
	err := c.request(ctx, c.endpoint+"/task/stats", &output)
	return output, err
}
