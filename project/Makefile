FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp


up: 
	@echo "Starting up the docker images ..."
	docker-compose up -d
	@echo "Docker images started !"


## This command will stop all the docker images and build them again
up_build: build_broker build_auth start
	@echo "Stopping the docker images ..."
	docker-compose down
	@echo "Building the docker images when required ... and starting the docker images"
	docker-compose up --build -d

down: stop
	@echo "Stopping the docker images ..."
	docker-compose down
	@echo "Docker images stoppcded !"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd
	@echo "Done!"

build_auth:
	@echo "Building broker binary..."
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./api
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../frontend/web && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY}
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../frontend/web && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"
