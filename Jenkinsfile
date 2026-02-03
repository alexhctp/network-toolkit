pipeline {
    agent any
    
    tools {
        go 'Go setup'
    }
    
    environment {
        CGO_ENABLED = '0'
        GOOS = 'linux'
        GOARCH = 'amd64'
    }
    
    stages {
        stage('Build') {
            steps {
                sh 'go build -o network-toolkit'
            }
        }
        
        stage('Archive') {
            steps {
                archiveArtifacts artifacts: 'network-toolkit', fingerprint: true
            }
        }
    }
}
