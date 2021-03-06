# get-ecs-services
A script to get versions of services in a cluster and the difference between two clusters
The version is pulled from the SSM VERSION parameter with the following syntax "/<cluster>/<service-name>/VERSION"


    Usage:
      get-ecs-services [flags]

    Flags:
      --cluster string           cluster to get list services and versions (default "dev")
      --profile string           override default AWS profile (default will discover default AWS profile)
      --diffCluster string       cluster to compare difference
      --diffProfile string       override AWS profile for '-diffCluster' only if a cluster exists in a different AWS account
      --help                     help for get-ecs-services
      --version                  show version for get-ecs-services


## Examples

Get service versions for every service running in a cluster

    get-ecs-services -cluster dev

Diff service versions for each service running in either cluster

    get-ecs-services -cluster dev -cluster tst

Diff service versions for each service in two clusters existing on different AWS accounts

    get-ecs-services -cluster tst -cluster prd -diffProfile my-prd-aws-profile
