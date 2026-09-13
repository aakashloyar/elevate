SHELL := /bin/bash

ROOT_DIR := $(CURDIR)
RUN_DIR := $(ROOT_DIR)/.run
PID_DIR := $(RUN_DIR)/pids
LOG_DIR := $(RUN_DIR)/logs

SERVICES := user assessment problem submission evaluation assessment_runner problem_generation
PORTS := 8081 8082 8083 8084 8085 8086 8087

.PHONY: run stop

run:
	@mkdir -p "$(PID_DIR)" "$(LOG_DIR)"
	@for service in $(SERVICES); do \
		pid_file="$(PID_DIR)/$$service.pid"; \
		if [ -f "$$pid_file" ] && kill -0 "$$(cat "$$pid_file")" 2>/dev/null; then \
			echo "$$service already running with PID $$(cat "$$pid_file")"; \
		else \
			echo "Starting $$service..."; \
			( cd "$(ROOT_DIR)/$$service" && go run ./cmd/server/main.go > "$(LOG_DIR)/$$service.log" 2>&1 & echo $$! > "$$pid_file" ); \
			echo "$$service started with PID $$(cat "$$pid_file")"; \
		fi; \
	done

stop:
	@mkdir -p "$(PID_DIR)"
	@for service in $(SERVICES); do \
		pid_file="$(PID_DIR)/$$service.pid"; \
		if [ -f "$$pid_file" ]; then \
			pid="$$(cat "$$pid_file")"; \
			if kill -0 "$$pid" 2>/dev/null; then \
				echo "Stopping $$service with PID $$pid..."; \
				pkill -TERM -P "$$pid" 2>/dev/null || true; \
				kill "$$pid" 2>/dev/null || true; \
			else \
				echo "$$service is not running"; \
			fi; \
			rm -f "$$pid_file"; \
		else \
			echo "$$service is not running"; \
		fi; \
	done
	@sleep 1
	@for port in $(PORTS); do \
		pids="$$(lsof -ti tcp:$$port -sTCP:LISTEN 2>/dev/null || true)"; \
		if [ -n "$$pids" ]; then \
			echo "Stopping process on port $$port: $$pids"; \
			kill $$pids 2>/dev/null || true; \
		fi; \
	done
