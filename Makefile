#INFO: default logfile path, is configurable
LOGFILE=htty.log


.PHONY: dev
dev:
	LOGLEVEL=all CONFIG_FILE="$(PWD)/config.json"  go run .

.PHONY: debug 
debug:
	LOGLEVEL=debug go run .

.PHONY: build
build:
	go build -o htty .

.PHONY: logwatch
logwatch:
	tail -f $(LOGFILE)

.PHONY: logflush
logflush: 
	echo "" > $(LOGFILE)

.PHONY: logwatch_new
logwatch_new:
	$(MAKE) logflush 
	$(MAKE) logwatch


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
	@echo "dev         run htty local in dev mode (all logs - info, warn, error, debug are shown) with ./config.json"
	@echo "debug       run htty local in debug mode"
	@echo "build       build htty executable"
	@echo "test        run test from ./tests folder"
	@echo "logwatch    follow $(LOGFILE) for viewing live logs" 

