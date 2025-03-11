pipeline {
    agent any
    stages {
        stage('Checkout') {
            steps {
                git branch: 'main', url: 'https://github.com/Malik-VTI/Auth-Services.git'
            }
        }

        stage('Install Dependencies') {
            steps {
                sh '''
                    go install github.com/DataDog/orchestrion@latest
                    orchestrion pin
                '''
            }
        }

        stage('Build') {
            steps {
                sh 'orchestrion go build -o auth-services-app .'
            }
        }

        stage('Test') {
            steps {
                sh 'orchestrion go test ./...'
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
<<<<<<< HEAD:jenkinsfile
=======

    post {
        always {
            cleanWs()
        }
    }
>>>>>>> f6e952a3a7823521985c400c138c77f41bddf5b9:Jenkinsfile
}
