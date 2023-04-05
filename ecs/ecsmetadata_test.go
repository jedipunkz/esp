package ecs

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewClientToMetadataEndpoint(t *testing.T) {
	t.Run("valid endpoint", func(t *testing.T) {
		originalEnv := os.Getenv(ecsMetadataUriEnvV4)
		defer os.Setenv(ecsMetadataUriEnvV4, originalEnv)

		expectedEndpoint := "http://example.com"
		os.Setenv(ecsMetadataUriEnvV4, expectedEndpoint)

		client, err := NewClientToMetadataEndpoint()
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if client.endpoint != expectedEndpoint {
			t.Errorf("expected endpoint %s, got: %s", expectedEndpoint, client.endpoint)
		}
	})

	t.Run("missing environment variable", func(t *testing.T) {
		originalEnv := os.Getenv(ecsMetadataUriEnvV4)
		defer os.Setenv(ecsMetadataUriEnvV4, originalEnv)

		os.Unsetenv(ecsMetadataUriEnvV4)

		_, err := NewClientToMetadataEndpoint()
		if err == nil {
			t.Fatal("expected an error, got nil")
		}
	})

	t.Run("invalid endpoint", func(t *testing.T) {
		originalEnv := os.Getenv(ecsMetadataUriEnvV4)
		defer os.Setenv(ecsMetadataUriEnvV4, originalEnv)

		invalidEndpoint := "http://%zz.example.com"
		os.Setenv(ecsMetadataUriEnvV4, invalidEndpoint)

		_, err := NewClientToMetadataEndpoint()
		if err == nil {
			t.Fatal("expected an error, got nil")
		}

		if !strings.Contains(err.Error(), "invalid endpoint") {
			t.Errorf("expected error containing 'invalid endpoint', got: %v", err)
		}
	})
}

func TestClient_RetriveTaskMetadata(t *testing.T) {
	ctx := context.Background()
	expected := TaskMetadata{
		Cluster: "test-cluster",
		TaskARN: "arn:aws:ecs:region:account-id:task/task-id",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/task" {
			http.NotFound(w, r)
			return
		}
		data, _ := json.Marshal(expected)
		if _, err := w.Write(data); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	taskMetadata, err := client.RetriveTaskMetadata(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(taskMetadata, expected) {
		t.Errorf("expected task metadata %v, got: %v", expected, taskMetadata)
	}
}

func TestClient_RetriveContainersMetadata(t *testing.T) {
	ctx := context.Background()
	expected := map[string]ContainerMetadata{
		"container1": {
			CPUStats: struct {
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
			}{
				CPUUsage: struct {
					TotalUsage        int   `json:"total_usage"`
					PerCPUUsage       []int `json:"percpu_usage"`
					UsageInKernelmode int   `json:"usage_in_kernelmode"`
					UsageInUsermode   int   `json:"usage_in_usermode"`
				}{
					TotalUsage:        1000,
					PerCPUUsage:       []int{500, 500},
					UsageInKernelmode: 600,
					UsageInUsermode:   400,
				},
				SystemCPUUsage: 5000,
				OnlineCPUs:     2,
				ThrottlingData: struct {
					Periods          int `json:"periods"`
					ThrottledPeriods int `json:"throttled_periods"`
					ThrottledTime    int `json:"throttled_time"`
				}{
					Periods:          10,
					ThrottledPeriods: 2,
					ThrottledTime:    100,
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/task/stats" {
			http.NotFound(w, r)
			return
		}
		data, _ := json.Marshal(expected)
		if _, err := w.Write(data); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	containersMetadata, err := client.RetriveContainersMetadata(ctx)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !reflect.DeepEqual(containersMetadata, expected) {
		t.Errorf("expected containers metadata %v, got: %v", expected, containersMetadata)
	}
}

func TestClient_request(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"key": "value"}`)); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client := NewClient(server.URL)
		var result map[string]string

		err := client.request(ctx, server.URL+"/success", &result)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		expected := map[string]string{"key": "value"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got: %v", expected, result)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte("invalid-json")); err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		}))
		defer server.Close()

		client := NewClient(server.URL)
		var result map[string]string

		err := client.request(ctx, server.URL+"/invalid-json", &result)
		if err == nil {
			t.Fatal("expected an error, got nil")
		}

		var expectedError *json.SyntaxError
		if !errors.As(err, &expectedError) {
			t.Errorf("expected error of type %T, got: %v", expectedError, err)
		}
	})
}

func TestClient_request_NotFound(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notfound" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if _, err := w.Write([]byte(`{"key": "value"}`)); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL)
	var result map[string]string

	err := client.request(ctx, server.URL+"/notfound", &result)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	var expectedError *ErrorNotFound
	if !errors.As(err, &expectedError) {
		t.Errorf("expected error of type %T, got: %v", expectedError, err)
	}
}
