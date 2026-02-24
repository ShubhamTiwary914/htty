

.PHONY: dev
dev:
	MODE=debug go run .

.PHONY: build
build:
	go build -o htty .

.PHONY: logwatch
logwatch:
	tail -f ./debug.log
	
.PHONY: help
help:
	@echo "HttY makefile guide (source: Makefile)"
	@echo "Usage: make [Target]"
	@echo ""
	@echo "[Targets]"
	@echo "dev         run htty local in dev/debug mode"
	@echo "build       build htty executable"
	@echo "logwatch    follow debug.log for viewing live logs" 

