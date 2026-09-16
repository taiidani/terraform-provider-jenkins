BINARY=terraform-provider-jenkins

default: build

# Builds the provider and adds it to your GOBIN folder.
build:
	go install
	@echo "Binary has been compiled to $(shell go env GOBIN)/${BINARY}"
	@echo "In order to have Terraform pick this up you will need to add the following to your $$HOME/.terraformrc file:"
	@echo "  provider_installation {"
	@echo "    dev_overrides {"
	@echo "      \"taiidani/jenkins\" = \"$(shell go env GOBIN)\""
	@echo "    }"
	@echo "    direct {}"
	@echo "  }"
	@echo ""
	@echo "This should only be used during development. See https://www.terraform.io/docs/commands/cli-config.html#development-overrides-for-provider-developers for details."

# Cleans up any lingering items in your system created by this provider.
clean:
	rm -f "$(shell go env GOBIN)/$(BINARY)"
