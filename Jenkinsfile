#!/usr/bin/env groovy

pipeline {
    agent {
        label 'terraform'
    }
    options {
        skipDefaultCheckout()
        disableConcurrentBuilds()
        ansiColor('xterm')
    }
    parameters {
        string(name: 'RELEASE_VERSION',
            defaultValue: '0.0.1',
            description: 'The version of the terraform provider')
    }
    stages {
        stage("Release") {
            steps {
                script {
                    assert params.RELEASE_VERSION
                    terraformProviderRelease(releaseVersion: params.RELEASE_VERSION)
                }
            }
        }
    }
}
