.PHONY: serve

serve:
	@trap 'kill 0' SIGINT; make -C backend serve & make -C frontend dev
