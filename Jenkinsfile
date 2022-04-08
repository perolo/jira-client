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
                sh 'go version'
                sh 'go install honnef.co/go/tools/cmd/staticcheck@latest'
                sh 'go install github.com/jstemmer/go-junit-report@latest'
                sh 'go install github.com/axw/gocov/gocov@latest'
                sh 'go install github.com/AlekSi/gocov-xml@latest'
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
                    //echo 'Running staticcheck'
                    sh 'staticcheck ./...'
                    echo 'Running test'
                    sh 'go test -v 2>&1 | go-junit-report > report.xml'
                    echo 'Running coverage'
                    sh 'gocov test ./... | gocov-xml > coverage.xml'
                }
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'report.xml', fingerprint: true
            archiveArtifacts artifacts: 'coverage.xml', fingerprint: true
            junit 'report.xml'
            cobertura coberturaReportFile: 'coverage.xml'
        }
    }        
}
