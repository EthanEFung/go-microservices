.DEFAULT_GOAL := swagger

check_install:
	which swagger || go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger: check_install
	@echo Ensure you have the swagger CLI or this command will fail.
	swagger generate spec -o ./swagger.yaml --scan-models