BINARY=terraform-provider-jenkins
DOCKER_URL=localhost

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
	# @docker build -t jenkins-provider-acc jenkins/test-fixtures/
	# @cd ./jenkins/test-fixtures && terraform init
	# @cd ./jenkins/test-fixtures && terraform taint docker_container.jenkins
	# @cd ./jenkins/test-fixtures && terraform apply -auto-approve
	# while [ "$$(docker inspect jenkins-provider-acc --format '{{ .State.Health.Status }}')" != "healthy" ]; do echo "Waiting for Jenkins to start..."; sleep 3; done
	TF_ACC=1 JENKINS_URL="http://${DOCKER_URL}:8080" JENKINS_USERNAME="admin" JENKINS_PASSWORD="admin" go test -v -cover ./...
	# @cd ./jenkins/test-fixtures && terraform destroy -auto-approve

# Cleans up any lingering items in your system created by this provider.
clean:
	rm -f "$$HOME/.terraform.d/plugins/$(BINARY)"
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"
