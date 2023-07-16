# Integration

This folder contains an example Dockerized instance of Jenkins for performing manual testing against. It can be used during development to validate against a real Jenkins installation.

## Usage

Start Jenkins with:

```sh
docker-compose up --detach
```

And then open it in your web browser at http://localhost:8080.

You can now run `terraform init`, `plan`, `apply`, etc. commands within this directory to test the provider against the Jenkins instance. If you are testing a version of the provider that you are developing locally, ensure that you've run the `make` command in the repository root and followed the provided instructions to configure your machine to use the built binary.

When done with testing, clean up the instance with:

```sh
docker-compose down --volumes
```
