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
