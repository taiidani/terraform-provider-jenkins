FROM jenkins/jenkins:lts

RUN /usr/local/bin/install-plugins.sh hashicorp-vault-plugin cloudbees-folder pipeline-model-definition git matrix-auth

HEALTHCHECK --interval=4s --start-period=5s --retries=30 CMD [ "curl", "-f", "http://localhost:8080" ]
