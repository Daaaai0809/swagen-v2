build-example:
	@echo "[INFO] Build example develop environment"
	@echo "[INFO] Copy .env"
	
	cp ./example/.env.example ./example/.env
	
	@echo "[INFO] Success copy .env"
	@echo "[INFO] Build swagen-v2.local"
	
	go build -o ./example/swagen-v2.local
	
	@echo "[INFO] Success build swagen-v2.local"

build-binary:
	@echo "[INFO] Build swagen-v2.local"

	go build -o ./example/swagen-v2.local

	@echo "[INFO] Success build swagen-v2.local"

check-golangci:
	@echo "[INFO] Check golangci-lint"
	golangci-lint run
	@echo "[INFO] Success check golangci-lint"

fix-golangci:
	@echo "[INFO] Fix golangci-lint"
	golangci-lint run --fix
	@echo "[INFO] Success fix golangci-lint"