pipeline {
    agent any

    environment {
        // Variabel lingkungan untuk Datadog
        DD_API_KEY = credentials('datadog-api-key') // Simpan API key di Jenkins Credentials
        GO_VERSION = '1.24.1' // Versi Go yang digunakan
        APP_NAME = 'auth-services'
        APP_DIR = '/opt/auth-services' // Direktori aplikasi
        BINARY_PATH = "${APP_DIR}/${APP_NAME}" // Lokasi binary aplikasi
    }

    stages {
        stage('Checkout Code') {
            steps {
                git branch: 'main', url: 'https://github.com/Malik-VTI/Auth-Services.git'
            }
        }

        stage('Install Dependencies') {
            steps {
                sh '''
                # Install Go jika belum ada
                if ! command -v go &> /dev/null; then
                    wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
                    sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
                    export PATH=$PATH:/usr/local/go/bin
                fi

                # Install Orchestrion
                go install github.com/DataDog/orchestrion/cmd/orchestrion@latest
                '''
            }
        }

        stage('Build Application') {
            steps {
                sh '''
                # Jalankan Orchestrion untuk instrumentasi otomatis
                orchestrion -w .

                # Buat direktori aplikasi jika belum ada
                mkdir -p ${APP_DIR}

                # Build aplikasi Go dan simpan binary ke direktori aplikasi
                go build -o ${BINARY_PATH} .
                '''
            }
        }

        stage('Deploy Application') {
            steps {
                sh '''
                # Hentikan layanan sebelumnya jika ada
                sudo systemctl stop ${APP_NAME} || true

                # Pastikan binary memiliki izin eksekusi
                sudo chmod +x ${BINARY_PATH}

                # Mulai ulang layanan
                sudo systemctl start ${APP_NAME}
                '''
            }
        }

        stage('Verify Deployment') {
            steps {
                sh '''
                # Verifikasi bahwa aplikasi berjalan
                sudo systemctl status ${APP_NAME}
                '''
            }
        }
    }

    post {
        success {
            echo 'Deployment berhasil!'
        }
        failure {
            echo 'Deployment gagal. Silakan periksa log.'
        }
    }
}