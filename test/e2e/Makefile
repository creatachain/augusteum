all: docker generator runner

docker:
	docker build --tag augusteum/e2e-node -f docker/Dockerfile ../..

# We need to build support for database backends into the app in
# order to build a binary with a Augusteum node in it (for built-in
# MSM testing).
app:
	go build -o build/app -tags badgerdb,boltdb,cleveldb,rocksdb ./app
	
# To be used primarily by the e2e docker instance. If you want to produce this binary
# elsewhere, then run go build in the maverick directory. 
maverick:
	go build -o build/maverick -tags badgerdb,boltdb,cleveldb,rocksdb ../maverick

generator:
	go build -o build/generator ./generator

runner:
	go build -o build/runner ./runner

.PHONY: all app docker generator maverick runner 
