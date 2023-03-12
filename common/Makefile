

build:
	@go mod tidy -go="${GOLANG_VERSION}"
	@echo "building common"
	@go build -o build/commonservice -a cmd/commonservice/commonservice.go
	@echo "build common"