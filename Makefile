# Simple Makefile to build/run examples with the right build tags.
#
# Usage:
#   make test
#   make run-httprouter
#   make run-httprouter-security
#
#   make run-gin
#   make run-gin-security
#
#   make run-echo
#   make run-echo-security
#
#   make run-fiber
#   make run-fiber-security

GO ?= go

.PHONY: help test tidy fmt lint \
	run-httprouter run-httprouter-security \
	run-gin run-gin-security \
	run-echo run-echo-security \
	run-fiber run-fiber-security

help:
	@echo "Targets:"
	@echo "  test                              - run unit tests"
	@echo "  tidy                              - go mod tidy"
	@echo "  fmt                               - gofmt all go files"
	@echo "  run-httprouter                     - run net/http example"
	@echo "  run-httprouter-security            - run net/http security example (-tags security)"
	@echo "  run-gin                            - run gin basic example (-tags gin)"
	@echo "  run-gin-security                   - run gin basic security example (-tags gin,security)"
	@echo "  run-echo                           - run echo basic example (-tags echo)"
	@echo "  run-echo-security                  - run echo basic security example (-tags echo,security)"
	@echo "  run-fiber                          - run fiber basic example (-tags fiber)"
	@echo "  run-fiber-security                 - run fiber basic security example (-tags fiber,security)"


test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...
	# gofmt on files that may not be in a package list in some setups
	gofmt -w $$(find . -name '*.go' -not -path './.git/*')

tidy:
	$(GO) mod tidy

# --- net/http examples ---
run-httprouter:
	$(GO) run ./example/httprouter

run-httprouter-security:
	$(GO) run -tags security ./example/httprouter

# --- Gin examples ---
run-gin:
	$(GO) run -tags gin ./example/gin

run-gin-security:
	$(GO) run -tags "gin,security" ./example/gin

# --- Echo examples ---
run-echo:
	$(GO) run -tags echo ./example/echo

run-echo-security:
	$(GO) run -tags "echo,security" ./example/echo

# --- Fiber examples ---
run-fiber:
	$(GO) run -tags fiber ./example/fiber

run-fiber-security:
	$(GO) run -tags "fiber,security" ./example/fiber
