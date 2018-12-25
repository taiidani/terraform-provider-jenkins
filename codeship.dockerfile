FROM circleci/golang:1.11

RUN wget https://releases.hashicorp.com/terraform/0.11.8/terraform_0.11.8_linux_amd64.zip \
    && unzip terraform_0.11.8_linux_amd64.zip \
    && sudo mv terraform /usr/local/bin/

WORKDIR /app
