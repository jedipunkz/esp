# ESP - ECS Stats Plotter

## Description

ESP Retrives AWS ECS Container Stats from Metadata Endpoint and Plots Stats to Cloudwatch high-resolution Custom Metrics.
ESP enables faster autoscaling of AWS ECS Tasks by refering to AWS Cloudwatch.

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
docker build --platform linux/amd64 -t esp:latest .
docker tag esp:latest *******.dkr.ecr.<REGION_NAME>.amazonaws.com/esp:latest
docker push *******.dkr.ecr.<REGION_NAME>.amazonaws.com/esp:latest
```

### Add ESP as a sidecar container

Add ESP as a ECS sidecar container in the container definition below.

```json
[
  {
    "name": "<CONTAINER_TO_MONITOR>",
    ...<snip>...
  },
  {
    "name": "esp",
    "image": "********.dkr.ecr.<REGION_NAME>.amazonaws.com/esp:latest",
    "essential": true,
    "environment": [
      {
        "name": "CONTAINER_NAME",
        "value": "<CONTAINER_TO_MONITOR>"
      },
      {
        "name": "REGION",
        "value": "<REGION_NAME>"
      },
      {
        "name": "NAMESPACE",
        "value": "<CLOUDWATCH_METRICS_NAMESPACE_NAME>"
      }
    ], 
    "dependsOn": [
      {
        "containerName": "<CONTAINER_TO_MONITOR>",
        "condition": "START"
      }
    ]
  }
]
```

| Environment Name | Description |
|---|---|
| CONTAINER_NAME | Container name for which you want to monitor |
| REGION | AWS region name |
| NAMESPACE | Cloudwatch metrics namespace name |

### Set Autoscalling

Configure Cloudwatch Metric Alarm for Autoscalling. Below is an example of Terraform Code.

```hcl
resource "aws_cloudwatch_metric_alarm" "foo" {
  ...snip...
  metric_name         = "CPUUtilization"
  namespace           = "<CLOUDWATCH_METRICS_NAMESPACE_NAME>"
  ...snip...
}

resource "aws_cloudwatch_metric_alarm" "foo" {
  ...snip...
  metric_name         = "CPUUtilization"
  namespace           = "<CLOUDWATCH_METRICS_NAMESPACE_NAME>"
  ...snip...
}
```

## Author

[jedipunkz](https://twitter.com/jedipunkz)

## License
The source code is licensed MIT. The website content is licensed CC BY 4.0,see LICENSE.
