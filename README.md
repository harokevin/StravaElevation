# StravaElevation
Adds comments to Strava rides to track custom goals

## How To Run
```
client_id=<> client_secret=<> refresh_token=<> go run StravaElevation.go
```

## How to run on AWS Lambda
1. Setup environment variables on the lambda
2. build to AWS lambdas spec
```
GOARCH=amd64 GOOS=linux go run StravaElevation.go
```
3. zip files
```
zip StravaElevation.zip StravaElevation StravaElevation.go
```
4. Upload zip to AWS

## How to deploy using terraform
1. install terraform CLI
1. install AWS CLI
1. configure AWS CLI `aws configure`
1. Make sure terraform.tf is referring to the correct AWS profile. ("default" by default)
1. `terraform plan` to create a change plan
1. `terraform apply` to apply(deploy) a plan
1. `terraform destroy` to destroy all infrastructure


## How to generate a Terrafom Graph
dot command privided by [Graphviz](https://graphviz.org/download)
```
terraform graph | dot -Tpng > graph.png
```

## Strava Authentication [Docs](https://developers.strava.com/docs/authentication/)

## How to contribute
- Use [Conventional Commits] when writing commit messages (https://www.conventionalcommits.org/en/v1.0.0/)
- Make a pull request
    - [ ] Descriptive title
    - [ ] Descriptive description
    - [ ] Related issues listed in the description
    - [ ] Rebase on top of origin/main
    - [ ] Squash code into a single commit