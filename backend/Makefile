.PHONY: serve build tidy

ifneq (,$(wildcard .env))
  include .env
  export
endif

serve:
	air

build:
	go build -o ./tmp/server .

tidy:
	go mod tidy
