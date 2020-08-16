BINARY=terraform-provider-jenkins
export COMPOSE_FILE=./example/docker-compose.yml

default: build

# Builds the provider and adds it to your GOPATH/bin folder.
build:
	go install

# Registers the built provider against the local Terraform plugins directory, enabling it for use by Terraform.
install: build
	ln -sf "$(shell go env GOPATH)/bin/$(BINARY)" "$$HOME/.terraform.d/plugins/$(BINARY)"

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
	rm -f "$$HOME/.terraform.d/plugins/$(BINARY)"
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"
