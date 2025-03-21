pipeline {
    agent any

    environment {
        APP_NAME = 'services-auth'
        BINARY_NAME = 'auth-services'
        SYSTEMD_SERVICE_NAME = 'services-auth.service'
        DEPLOY_PATH = '/opt/auth-services'
        PATH = '/usr/local/go/bin'
    }

    stages {
        stage('Go Version') {
            steps {
                script {
                    sh 'go version || (echo "Go is not installed" && exit 1)'
                }
            }
        }
        stage('Build') {
            steps {
                script {
                    sh 'orchestrion go build -o ${BINARY_NAME} .'
                }
                archiveArtifacts artifacts: "${BINARY_NAME}", fingerprint: true
            }
        }

        stage('Test') {
            steps {
                script {
                    sh 'orchestrion go test ./...'
                }
            }
        }

        stage('Deploy') {
            agent any
            steps {
                script {
                    // Pindahkan binary ke lokasi deploy
                    sh """
                        mkdir -p ${DEPLOY_PATH}
                        cp ${WORKSPACE}/${BINARY_NAME} ${DEPLOY_PATH}/${APP_NAME}
                        chmod +x ${DEPLOY_PATH}/${APP_NAME}

                        // Konfigurasi dan restart systemd service
                        sudo systemctl stop ${SYSTEMD_SERVICE_NAME} || true
                        sudo systemctl daemon-reload
                        sudo systemctl start ${SYSTEMD_SERVICE_NAME}
                        sudo systemctl enable ${SYSTEMD_SERVICE_NAME}
                    """
                }
            }
        }
    }
}