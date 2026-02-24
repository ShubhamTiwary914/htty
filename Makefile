LOGFILE=debug.log


.PHONY: dev
dev:
	MODE=debug go run .

.PHONY: build
build:
	go build -o htty .

.PHONY: logwatch
logwatch:
	tail -f $(LOGFILE)

.PHONY: logflush
logflush: 
	echo "" > $(LOGFILE)


# yes this is not "recommended" pattern for tests in go, but simple to manage
# recommended way is file.go should have file_test.go in same dir (but that is stupid)
.PHONY: test
test: 
	go test ./tests -v
	
.PHONY: help
help:
	@echo "HttY makefile guide (source: Makefile)"
	@echo "Usage: make [Target]"
	@echo ""
	@echo "[Targets]"
	@echo "dev         run htty local in dev/debug mode"
	@echo "build       build htty executable"
	@echo "test        run test from ./tests folder"
	@echo "logwatch    follow $(LOGFILE) for viewing live logs" 

