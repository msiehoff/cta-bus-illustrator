.PHONY: serve import-route-segments

serve:
	@trap 'kill 0' SIGINT; make -C backend serve & make -C frontend dev

import-route-segments:
	bash ./scripts/import_route_segments.sh
