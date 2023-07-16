BINARY=terraform-provider-jenkins
export COMPOSE_FILE=./integration/docker-compose.yml

default: build

# Builds the provider and adds it to your GOPATH/bin folder.
build:
	go install
	@echo "Binary has been compiled to $(shell go env GOPATH)/bin/${BINARY}"
	@echo "In order to have Terraform pick this up you will need to add the following to your $$HOME/.terraformrc file:"
	@echo "  provider_installation {"
	@echo "    dev_overrides {"
	@echo "      \"taiidani/jenkins\" = \"$(shell go env GOPATH)/bin\""
	@echo "    }"
	@echo "    direct {}"
	@echo "  }"
	@echo ""
	@echo "This should only be used during development. See https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers for details."

# Formats TF files and generates documentation
generate:
	cd tools; go generate ./...

# Executes all unit tests for the provider
test:
	go test -cover ./...

# Executes all acceptance tests for the provider
testacc:
	@docker-compose build
	@docker-compose up -d --force-recreate jenkins
	@while [ "$$(docker inspect jenkins-provider-acc --format '{{ .State.Health.Status }}')" != "healthy" ]; do echo "Waiting for Jenkins to start..."; sleep 3; done
	TF_ACC=1 JENKINS_URL="http://localhost:8080" JENKINS_USERNAME="admin" JENKINS_PASSWORD="admin" go test -v -cover ./...
	@docker-compose down

# Cleans up any lingering items in your system created by this provider.
clean:
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"
