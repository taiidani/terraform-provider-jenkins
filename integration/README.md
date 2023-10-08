# Integration

This folder contains an example Dockerized instance of Jenkins for performing manual and automated testing against. It can be used during development to validate against a real Jenkins installation.

## Automated Tests

Use the [terraform test](https://developer.hashicorp.com/terraform/language/tests) framework to automatically spin a copy of Jenkins up, apply the resources defined in [main.tftest.hcl](./main.tftest.hcl) against it, then tear the installation down.

```sh
terraform init
terraform test
```

## Manual Testing

Start Jenkins with:

```sh
docker compose up --detach
```

And then open it in your web browser at http://localhost:8080.

You can now run `terraform init`, `plan`, `apply`, etc. commands within this directory test the provider against the Jenkins instance with your change. If you are testing a version of the provider that you are developing locally, ensure that you've run the `make` command in the repository root and followed the provided instructions to configure your machine to use the built binary.

When done with testing, clean up the instance with:

```sh
docker compose down --volumes
```
