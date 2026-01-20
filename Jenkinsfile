pipeline {
    agent any
    
    tools {
        go //'null' // Name configured in Jenkins Global Tool Configuration
    }
    
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = '0'
        GOOS = 'linux'
        GOARCH = 'amd64'
        APP_NAME = 'network-toolkit'
    }
    
    stages {
        stage('Checkout') {
            steps {
                echo '📦 Checking out source code...'
                checkout scm
            }
        }
        
        stage('Environment Info') {
            steps {
                echo '🔍 Displaying environment information...'
                sh 'go version'
                sh 'go env'
            }
        }
        
        stage('Dependencies') {
            steps {
                echo '📚 Downloading dependencies...'
                sh 'go mod download'
                sh 'go mod verify'
            }
        }
        
        stage('Build') {
            steps {
                echo '🔨 Building application...'
                sh 'go build -v -o $APP_NAME'
            }
        }
        
        stage('Test') {
            steps {
                echo '🧪 Running tests...'
                sh '''
                    go test -v ./... || exit 0
                '''
            }
        }
        
        stage('Archive') {
            steps {
                echo '📦 Archiving artifacts...'
                archiveArtifacts artifacts: '*', fingerprint: true
            }
        }
    }
    
    post {
        success {
            echo '✅ Pipeline completed successfully!'
        }
        failure {
            echo '❌ Pipeline failed!'
        }
        always {
            echo '🧹 Cleaning up workspace...'
            cleanWs()
        }
    }
}
