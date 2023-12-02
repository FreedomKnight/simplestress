main: server/server client/client

*.pb.go: proto/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto

server/server: server/main.go lib/*.go *.pb.go
	GOOS=linux GOARCH=amd64 go build -o server/server server/main.go

client/client: client/main.go lib/*.go *.pb.go
	GOOS=linux GOARCH=amd64 go build -o client/client client/main.go

clean:
	rm -f server/server client/client

build-client-docker: client
	docker build --platform linux/amd64 . -f deployments/docker/client.dockerfile -t freedomknight/simplestress-client

build-server-docker: client
	docker build --platform linux/amd64 . -f deployments/docker/server.dockerfile -t freedomknight/simplestress-server

build-docker: build-client-docker build-server-docker

docker-push: build-docker
	docker push freedomknight/simplestress-client
	docker push freedomknight/simplestress-server

deploy:
	kubectl apply -f deployments/k8s/report.yaml
	kubectl apply -f deployments/k8s/server.yaml
	kubectl apply -f deployments/k8s/client.yaml

.phony: deploy clean
