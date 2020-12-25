name = zero

run:
	go run ./cmd/zerod/main.go

migrate_up:
	go run ./cmd/migrate/main.go up

migrate_down:
	go run ./cmd/migrate/main.go down



# Database stuff

db_docker_name = ${name}
db_database_name = ${name}

db_create:
	sudo docker run --name ${db_docker_name} -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

db_delete:
	sudo docker rm -f ${db_docker_name}

db_start:
	sudo docker start ${db_docker_name}

db_stop:
	sudo docker stop ${db_docker_name}

db_psql:
	sudo docker run -it --rm --network host postgres psql -h localhost -U postgres ${db_database_name}

