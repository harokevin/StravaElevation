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


## [Strava Authentication Docs](https://developers.strava.com/docs/authentication/)
