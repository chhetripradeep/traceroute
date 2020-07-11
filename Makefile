GO ?= go

TRACEROUTE_BIN = traceroute
TRACEROUTE_CMD = cmd/traceroute.go

.PHONY: build
build:
	$(GO) build -o $(TRACEROUTE_BIN) $(TRACEROUTE_CMD)
