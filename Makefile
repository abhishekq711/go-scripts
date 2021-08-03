# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))


.PHONY: all
all: build

.PHONY: build
build:
	CGO_ENABLED=0 go build -o go-code-server

# DOCKER TASKS
# Build the container
.PHONY: image
image: ## Build the container
	docker build -t $(APP_NAME) .

# Run go program from the entrypoint
.PHONY: run
run:
	go run main.go


# Gofmt is a tool that automatically formats Go source code.
.PHONY: fmt
fmt:
	gofmt -s -w .

# Vet examines Go source code and reports suspicious constructs, 
# such as Printf calls whose arguments do not align with the format string
.PHONY: vet
vet:
	go vet .

.PHONY: clean
clean:
	rm -f go-code-server