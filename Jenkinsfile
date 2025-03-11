pipeline {
    agent any

    environment {
    PATH = "/usr/local/go/bin:$PATH"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/Malik-VTI/Auth-Services.git'
            }
        }

        stage('Install Dependencies') {
            steps {
                sh 'go install github.com/DataDog/orchestrion@latest'
            }
        }

        stage('Build') {
            steps {
                sh '''
                    go build -toolexec 'orchestrion toolexec' .
                '''
            }
        }

        stage('Test') {
            steps {
                sh '''
                    go test -toolexec 'orchestrion toolexec' -race .
                '''
            }
        }

        stage('Package') {
            steps {
                sh '''
                    sudo mkdir -p /opt/auth-services
                    sudo cp auth-services-app /opt/auth-services/
                    sudo chmod +x /opt/auth-services/auth-services-app
                '''
            }
        }

        stage('Restart Application') {
            steps {
                sh '''
                    echo "P@ssw0rd" | sudo -S systemctl restart auth-services.service
                '''
            }
        }
    }
}
