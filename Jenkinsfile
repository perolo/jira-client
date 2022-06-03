pipeline {
    agent {
        docker { image 'yocreo/go-docker:latest' }
    }
    environment {
        GO114MODULE = 'on'
        CGO_ENABLED = 0 
        GOPATH = "/go"
        HOME = "/home/perolo/Jenkins/workspace/${JOB_NAME}"
    }
    options { 
        buildDiscarder(logRotator(numToKeepStr: '10')) 
    }    
    stages {        
        
        stage('Build') {
            steps {
                echo 'Compiling and building'
                sh 'go build'
            }
        }

        stage('Test') {
            steps {
                withEnv(["PATH+GO=${GOPATH}/bin"]){
                    catchError(buildResult: 'SUCCESS', stageResult: 'UNSTABLE', message: 'Static codecheck errors!') {
                        echo 'Running vetting'
                        sh 'go vet .'
                        echo 'Running staticcheck'
                        sh 'staticcheck ./...'
                        echo 'Running golangci-lint'
                        sh 'golangci-lint run --out-format junit-xml --config .golangci.yml > golangci-lint.xml'
                    }
                    echo 'Running test'
                    sh 'go test -v 2>&1 | go-junit-report > report.xml'
                    echo 'Running coverage'
                    sh 'gocov test ./... | gocov-xml > coverage.xml'
                }
            }
        }
        stage('Vulnerabilities') {
            steps {
                echo 'Vulnerabilities'
                sh 'env'
                sh 'pwd'             
                sh 'go version'
                sh 'nancy -V'
                sh 'git --version'
                sh '/usr/local/bin/nancy -V'
                sh 'go list -json -m all | nancy sleuth'
            }
        }
        stage('Artifacts') {
            steps { 
                script {
                    if (fileExists('report.xml')) {
                        archiveArtifacts artifacts: 'report.xml', fingerprint: true
                        try {
                            junit 'report.xml'
                        } catch (err) {
                            echo err.getMessage()
                            echo "Error detected, but we will continue."
                            echo "No lint errors found is not an error."
                        }
                    }
                    if (fileExists('coverage.xml')) {
                        try {
                            archiveArtifacts artifacts: 'coverage.xml', fingerprint: true
                            cobertura coberturaReportFile: 'coverage.xml'
                        } catch (err) {
                            echo err.getMessage()
                            echo "Error detected, but we will continue."
                            echo "No lint errors found is not an error."
                        }
                    }
                    if (fileExists('golangci-lint.xml')) {
                        archiveArtifacts artifacts: 'golangci-lint.xml'            
                        try {
                            junit 'golangci-lint.xml'
                        } catch (err) {
                            echo err.getMessage()
                            echo "Error detected, but we will continue."
                            echo "No lint errors found is not an error."
                        }
                    }
                }
            }   
        }
    }
}
