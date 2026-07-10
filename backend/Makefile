.PHONY: serve build tidy

ifneq (,$(wildcard .env))
  include .env
endif

# Parent targets (make serve / make serve_all) override via command line.
PIPELINE_ENABLED ?= false
export

serve:
	air

build:
	go build -o ./tmp/server .

tidy:
	go mod tidy
