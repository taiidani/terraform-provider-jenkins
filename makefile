BINARY=terraform-provider-jenkins

default: build

# Builds the provider and adds it to your GOPATH/bin folder.
build:
	go install

# Registers the built provider against the local Terraform plugins directory, enabling it for use by Terraform.
install: build
	ln -sf "$(shell go env GOPATH)/bin/$(BINARY)" "$$HOME/.terraform.d/plugins/$(BINARY)"

# Executes all tests for the provider
test:
	go test -cover ./...

# Cleans up any lingering items in your system created by this provider.
clean:
	rm -f "$$HOME/.terraform.d/plugins/$(BINARY)"
	rm -f "$(shell go env GOPATH)/bin/$(BINARY)"
