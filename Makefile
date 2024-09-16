BUILD_BINARY=bin\f1-results-rod.exe

test:
	@echo start tests
	@go test .\cmd\app\ -v -count=1
	@echo - done

build:
	@echo start build
	@go build -o ${BUILD_BINARY} .\cmd\app
	@echo - done

run: build
	@echo start run
	@${BUILD_BINARY} -localRun=true
	@echo - done
