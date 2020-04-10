get-deps:
	@echo Installing dependencies
	go get github.com/denisenkom/go-mssqldb
	go install github.com/denisenkom/go-mssqldb