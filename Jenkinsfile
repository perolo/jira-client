pipeline {
    agent {
        docker { image 'golang:1.18-stretch' }
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "/home/perolo/Jenkins/workspace/${JOB_NAME}"
        HOME = "/home/perolo/Jenkins/workspace/${JOB_NAME}"
    }
    stages {        
        stage('Pre Test') {
            steps {
                echo 'Installing dependencies'
                sh 'env'
                sh 'pwd'
                sh 'ls -al'
                sh 'ls -al $GOPATH'
                sh 'whoami'
                sh 'go version'
            }
        }
        
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    echo 'Running vetting'
                    sh 'go vet .'
                    echo 'Running linting'
                    sh 'golint .'
                    echo 'Running test'
                    sh 'cd test && go test -v'
                }
            }
        }
        
    }
}
