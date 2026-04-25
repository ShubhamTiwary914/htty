#INFO: default logfile path, is configurable
LOGFILE=.logs/htty.log
CONFIG_FILE="$(PWD)/config.json"
CACHE_PREFIX="$(PWD)/.cache"
TMP_DIR="/tmp"

# build & run --------------

runbuild:  
	rm -f ./htty
	go build -o htty 
	./htty

.PHONY: dev
dev: 
	LOGLEVEL=all CONFIG_FILE=$(CONFIG_FILE) \
		TMP_DIR=$(TMP_DIR) \
		CACHE_PREFIX=$(CACHE_PREFIX) go run .

.PHONY: debug 
debug:
	LOGLEVEL=debug \
		CONFIG_FILE=$(CONFIG_FILE) \
		TMP_DIR=$(TMP_DIR) \
		CACHE_PREFIX=$(CACHE_PREFIX) go run .


.PHONY: build
build: 
	rm -f ./htty
	go build -o htty


.PHONY: clean
clean:
	rm -f ./htty
	$(MAKE) logflush


.PHONY: setup
setup:
	bash -c ./setup/setup.sh


# logs & debugging --------------
.PHONY: logwatch
logwatch:
	if [[ -z "$$(command -v lnav)" ]]; then \
		tail -f "$(LOGFILE)"
	else
		lnav "$(LOGFILE)"	
	fi

.PHONY: logflush
logflush: 
	echo "" > $(LOGFILE)

.PHONY: logwatch_new
logwatch_new:
	$(MAKE) logflush 
	$(MAKE) logwatch

# tests & docs --------------
.PHONY: test
test: 
	go test ./tests -v

.PHONY: rough
roughtest:
	CACHE_PREFIX="$(PWD)" go test ./tests/rough_test.go -v

.PHONY: docslive
docslive:
	go doc -http


.PHONY: stats
.ONESHELL:
stats: 
	$(MAKE) build
	@exec_size=$$(du -sh ./htty | awk '{print $$1}')
	@echo ""
	@echo "$$(go version)"
	@echo "Binary(htty) size: $$exec_size"
	@echo "Lines of Code: $$(find . -type f -name '*.go' -exec cat {} + | awk 'END{print NR}')"
	@echo "Current Branch: $$(git branch --show-current)"
	@echo "Last Updated by $$(git log -1 --pretty=format:'%an at %ad' --date=default)"
	@rm -f ./htty

.PHONY: help
help:
	@echo "HttY makefile guide (source: Makefile)"
	@echo "Usage: make [Target]"
	@echo ""
	@echo "[Targets]"
	@echo "dev         run htty local in dev mode (all logs - info, warn, error, debug are shown) with $(PWD)/config.json"
	@echo "debug       run htty local in debug mode"
	@echo "setup	   initial setup script & default completions to load once"
	@echo "build       build htty executable"
	@echo "runbuild    build htty(if not already) + run executable"
	@echo "clean       remove build executable & clean log file($(LOGFILE))"
	@echo "test        run test from ./tests folder"
	@echo "roughtest   run test for only the rough ./tests/rough_test.go"
	@echo "logwatch    follow $(LOGFILE) for viewing live logs" 
	@echo "docslive    gopls docs for htty project in live webapp"
	@echo "stats       stats like bin size, other utilities, ..."
