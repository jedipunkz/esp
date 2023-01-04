package cloudwatch

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const (
	metricName = "CPUUsage"
)

type Cloudwatch struct {
	sess   *session.Session
	svc    *cloudwatch.CloudWatch
	params *cloudwatch.PutMetricDataInput
}

func NewCloudwatch(params *cloudwatch.PutMetricDataInput) *Cloudwatch {
	// ToDo: get current region
	region := os.Getenv("REGION")
	if region == "" {
		log.Printf("error retriving region from environment variable")
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	svc := cloudwatch.New(sess)
	return &Cloudwatch{sess, svc, params}
}

func PutMetricData(value float64) error {
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		log.Printf("error retriving namespace from environment variable")
	}

	params := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			{
				MetricName:        aws.String(metricName),
				Unit:              aws.String(cloudwatch.StandardUnitPercent),
				Value:             aws.Float64(value),
				StorageResolution: aws.Int64(1),
			},
		},
		Namespace: aws.String(namespace),
	}

	client := NewCloudwatch(params)

	_, err := client.svc.PutMetricData(params)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	// log.Println(resp)

	return nil
}
