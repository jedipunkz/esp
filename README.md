# ESP - ECS Stats Plotter

## Description

ESP Retrives AWS ECS Container Stats from Metadata Endpoint and Plots Stats to Cloudwatch Detailed Metrics.
ESP enables faster autoscaling of AWS ECS Tasks by indexing AWS Cloudwatch high-resolution Custom Metrics.

## Requirement

### Add cloudwatch:PutMetricData Permission to Task Role.

add `cloudwatch:PutMetricData` to task role.

```hcl
  statement {
    effect = "Allow"
    actions = [
      "cloudwatch:PutMetricData"
    ]
    resources = ["*"]
  }
```

## Usage

### Build and Push to Repository

Build docker image and push to ECR repository.

```shell
docker build --platform linux/amd64 -t esp .
docker push *******.dkr.ecr.us-east-1.amazonaws.com/esp:latest
```

### Add ESP as a sidecar container

Add ESP as a ECS sidecar container in the container definition below.

```json
[
  {
    "name": "web",
    ...<snip>...
  {
    "name": "esp",
    "image": "********.dkr.ecr.us-east-1.amazonaws.com/esp:latest",
    "essential": true,
    "environment": [
      {
        "name": "CONTAINER_NAME",
        "value": "web"
      },
      {
        "name": "REGION",
        "value": "us-east-1"
      },
      {
        "name": "NAMESPACE",
        "value": "FOO"
      }
    ], 
    "dependsOn": [
      {
        "containerName": "web",
        "condition": "START"
      }
    ]
  }
]
```

| Environment Name | Description |
|---|---|
| CONTAINER_NAME | Container name for which you want to measure impossibility |
| REGION | AWS Region Name |
| NAMESPACE | Namespace name for Cloudwatch high resolution custom metrics |

### Set Autoscalling

Configure Cloudwatch Metric Alarm for Autoscalling. Below is an example of Terraform Code.

```hcl
resource "aws_cloudwatch_metric_alarm" "foo" {
  alarm_name          = "foo-CPU-Utilization-High-30"
  comparison_operator = "GreaterThanOrEqualToThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "FOO"
  period              = "60"
  statistic           = "Average"
  threshold           = "15"
  dimensions = {
    ClusterName = aws_ecs_cluster.foo.name
    ServiceName = aws_ecs_service.foo.name
  }
  alarm_actions = [aws_appautoscaling_policy.scale_out.arn]
}

resource "aws_cloudwatch_metric_alarm" "foo" {
  alarm_name          = "foo-CPU-Utilization-Low-5"
  comparison_operator = "LessThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "CPUUtilization"
  namespace           = "FOO"
  period              = "60"
  statistic           = "Average"
  threshold           = "5"
  dimensions = {
    ClusterName = aws_ecs_cluster.foo.name
    ServiceName = aws_ecs_service.foo.name
  }
  alarm_actions = [aws_appautoscaling_policy.scale_in.arn]
}
```

## Author

[jedipunkz](https://twitter.com/jedipunkz)

## License
The source code is licensed MIT. The website content is licensed CC BY 4.0,see LICENSE.
