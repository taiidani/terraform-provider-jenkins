FROM jenkins/jenkins:2.568.3-lts

RUN jenkins-plugin-cli --plugins \
    azure-credentials:426.v5627e6a_4c7a_7 \
    hashicorp-vault-plugin:384.vda_86ec66c537 \
    cloudbees-folder:6.1106.v3a_d9a_6d2465e \
    pipeline-model-definition:2.2293.v6e7193cec599 \
    git:5.10.1 \
    matrix-auth:3.3 \
    aws-credentials:265.v8422a_f384cd9

HEALTHCHECK --interval=4s --start-period=5s --retries=30 CMD [ "curl", "-f", "http://localhost:8080" ]
