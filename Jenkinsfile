pipeline {
    agent {
        docker { image 'golang:1.18-stretch' }
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "/go"
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
                    //echo 'Running linting'
                    //sh 'golint .'
                    //sh 'staticcheck ./...'
                    echo 'Running test'
                    sh 'go test -v'
                }
            }
        }
        
    }
}
