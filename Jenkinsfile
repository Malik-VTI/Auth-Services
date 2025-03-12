pipeline {
    agent any

    environment {
        GO_VERSION = '1.24.1' // Versi Go yang digunakan
        APP_NAME = 'auth-services'
        APP_DIR = '/opt/auth-services' // Direktori aplikasi
        BINARY_PATH = "${APP_DIR}/${APP_NAME}" // Lokasi binary aplikasi
        GOPATH = "/usr/local/go/bin:/home/malik/go/bin:${env.PATH}"
    }

    stages {
        stage('Checkout Code') {
            steps {
                git branch: 'main', url: 'https://github.com/Malik-VTI/Auth-Services.git'
            }
        }

        // stage('Install Orchestrion') {
        //     steps {
        //         sh '''
        //         # Install Orchestrion jika belum ada
        //         if ! command -v orchestrion &> /dev/null; then
        //             go install github.com/DataDog/orchestrion/cmd/orchestrion@latest
        //         fi
        //         '''
        //     }
        // }

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