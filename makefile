.PHONY: test innertest fasttest

innertest:
	docker-compose -f docker-compose-test.yml -p test up --force-recreate --remove-orphans --exit-code-from tester --abort-on-container-exit
	docker-compose -f docker-compose-test.yml -p test kill
	docker-compose -f docker-compose-test.yml -p test rm --force -v

test:
	docker-compose -f docker-compose-test.yml -p test build --no-cache
	make innertest

fasttest:
	docker-compose -f docker-compose-test.yml -p test build
	make innertest