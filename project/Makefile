FRONT_END_BINARY=frontApp
BROKER_BINARY=brokerApp


up: 
	@echo "Starting up the docker images ..."
	docker-compose up -d
	@echo "Docker images started !"


## This command will stop all the docker images and build them again
up_build: build_broker
	@echo "Stopping the docker images ..."
	docker-compose down
	@echo "Building the docker images when required ... and starting the docker images"
	docker-compose up --build -d

down:
	@echo "Stopping the docker images ..."
	docker-compose down
	@echo "Docker images stopped !"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../frontend && env CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./web
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../frontend && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"
