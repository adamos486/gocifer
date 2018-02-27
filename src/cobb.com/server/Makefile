GOSRC = $(shell find . -name "*.go" ! -name "*test.go" ! -name "*fake*")
out/app : $(GOSRC)
	go build -o ../../../out/app
	cp ../../../.env ../../../out/.env
build : out/app
format :
	go fmt ./...
test : units build
units : format vet
	ginkgo -r
	@echo "\nSWEET SUCCESS"
vet :
	@echo "Vetting packages for potential issues..."
	go tool vet -all -shadow=true ./database ./events ./fakes ./test ./utils