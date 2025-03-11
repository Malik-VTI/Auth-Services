pipeline {
    agent any

    environment {
        PATH = "/usr/local/go/bin:$PATH"
        DD_AGENT_HOST = "localhost"
        DD_TRACE_AGENT_PORT = "8126"
        DD_ENV = "poc"
        DD_SERVICE = "auth-services"
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
                    go build -toolexec 'orchestrion toolexec' -o auth-services-app .
                '''
            }
        }

        stage('Test') {
            steps {
                sh '''
                    go test -toolexec 'orchestrion toolexec' -race ./...
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

        stage('Deploy') {
            steps {
                sh '''
                    echo "[Unit]
                    Description=Auth Services

                    [Service]
                    ExecStart=/opt/auth-services/auth-services-app
                    Restart=always
                    User=nobody
                    Group=nogroup
                    Environment=PATH=/usr/bin:/usr/local/bin
                    Environment=GO_ENV=production
                    WorkingDirectory=/opt/auth-services

                    [Install]
                    WantedBy=multi-user.target" | sudo tee /etc/systemd/system/auth-services.service

                    sudo systemctl daemon-reload
                    sudo systemctl enable auth-services
                    sudo systemctl restart auth-services
                    sudo systemctl restart datadog-agent
                '''
            }
        }
    }

    post {
        always {
            cleanWs()
        }
    }
}