.PHONY: serve serve_all import-route-segments headway-run seed-headways

serve:
	@trap 'kill 0' SIGINT; \
	make -C backend serve PIPELINE_ENABLED=false & make -C frontend dev

serve_all:
	@trap 'kill 0' SIGINT; \
	make -C backend serve PIPELINE_ENABLED=true & make -C frontend dev

import-route-segments:
	bash ./scripts/import_route_segments.sh

# Trigger headway rollup against a running local server.
# Usage: make headway-run
#        make headway-run DATE=2026-07-10
#        make headway-run TOKEN=secret DATE=2026-07-10
APP_URL ?= http://localhost:8080
DATE ?=
TOKEN ?=

headway-run:
	@BODY='{}'; \
	if [ -n "$(DATE)" ]; then BODY=$$(printf '{"service_date":"%s"}' "$(DATE)"); fi; \
	if [ -n "$(TOKEN)" ]; then \
	  curl -sf -X POST "$(APP_URL)/api/v1/admin/headways/run" \
	    -H "Authorization: Bearer $(TOKEN)" \
	    -H "Content-Type: application/json" \
	    -d "$$BODY"; \
	else \
	  echo "Set TOKEN=... (HEADWAY_JOB_TOKEN) or use the admin UI at /admin/headways"; \
	  exit 1; \
	fi
	@echo

# Local-only sample arrivals + headway rollups for /headways UI.
# Requires DATABASE_URL (backend/.env). Rollups need a running API + HEADWAY_JOB_TOKEN.
# Usage: make seed-headways
#        make seed-headways TOKEN=local-dev-token
seed-headways:
	TOKEN="$(TOKEN)" APP_URL="$(APP_URL)" bash ./scripts/seed_headways.sh
