# go-sample-app
go fiber , go viper,golang zap logger,Datadog , opentelemetry, swagger

go install github.com/swaggo/swag/cmd/swag@latest

swag init -g .\cmd\commonservice\commonservice.go -o docs/ --parseDependency --parseInternal