.PHONY: docker-build
docker-build:
	$(MAKE) -C ./ext_proc
	$(MAKE) -C ./http_server

.PHONY: run
run:
	docker-compose up -d

.PHONY: stop
stop:
	docker-compose down -v