.PHONY: serve serve_all import-route-segments

serve:
	@trap 'kill 0' SIGINT; \
	make -C backend serve PIPELINE_ENABLED=false & make -C frontend dev

serve_all:
	@trap 'kill 0' SIGINT; \
	make -C backend serve PIPELINE_ENABLED=true & make -C frontend dev

import-route-segments:
	bash ./scripts/import_route_segments.sh
