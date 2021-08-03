# import deploy config
# You can change the default deploy config with `make cnf="deploy_special.env" release`
dpl ?= deploy.env
include $(dpl)
export $(shell sed 's/=.*//' $(dpl))


# DOCKER TASKS
# Build the container
build: ## Build the container
	docker build -t $(APP_NAME) .


.PHONY: run
run:
	go run main.go

# .PHONY: run
# run: ## Run container on port configured in `deploy.env`
# 	docker run -i -t -d -p=$(PORT):$(PORT) --name="$(APP_NAME)" $(APP_NAME)

# up: build run ## Run container on port configured in `config.env` (Alias to run)