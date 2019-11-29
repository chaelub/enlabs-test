prepare-volume:
	mkdir pgvolume

run-compose:
	$(info starting services. http interface available on 8001 port)
	@docker-compose up --build -d
	@cat enlabs_db | docker exec -i enlabs-test_db_1 psql -U postgres -d test
