BINARY=terraform-provider-jenkins
ACC_DOCKER_URL=localhost
ACC_USE_DOCKER=true
export COMPOSE_FILE=./jenkins/test-fixtures/docker-compose.yml

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
	@if [ "${ACC_USE_DOCKER}" == "true" ]; then docker-compose build; fi
	@if [ "${ACC_USE_DOCKER}" == "true" ]; then docker-compose up -d --force-recreate jenkins; fi
	@if [ "${ACC_USE_DOCKER}" == "true" ]; then while [ "$$(docker inspect jenkins-provider-acc --format '{{ .State.Health.Status }}')" != "healthy" ]; do echo "Waiting for Jenkins to start..."; sleep 3; done; fi
	TF_ACC=1 JENKINS_URL="http://${ACC_DOCKER_URL}:8080" JENKINS_USERNAME="admin" JENKINS_PASSWORD="admin" go test -v -cover ./...
	@if [ "${ACC_USE_DOCKER}" == "true" ]; then docker-compose down; fi

# Cleans up any lingering items in your system created by this provider.
clean:
	rm -f "$$HOME/.terraform.d/plugins/$(BINARY)"
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"
