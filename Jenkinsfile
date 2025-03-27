pipeline {
    agent any

    environment {
        REGISTRY = "docker.io/malikvti"
        IMAGE_NAME = "auth-service"
        IMAGE_TAG = "1.0"
        KUBECONFIG_CREDENTIAL = 'kubeconfig' // ID credential di Jenkins
        DOCKERHUB_CREDENTIAL = 'dockerhub-credential' // ID credential DockerHub
    }

    stages {
        stage('Checkout Code') {
            steps {
                git branch: 'main', url: 'https://github.com/Malik-VTI/Auth-Services.git'
            }
        }

        stage('Build with Orchestrator') {
            steps {
                script {
                    sh '''
                    orchestrion go build . --output=auth-service
                    '''
                }
            }
        }

        stage('Build & Push Docker Image') {
            steps {
                withCredentials([usernamePassword(credentialsId: "${DOCKERHUB_CREDENTIAL}", usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                    sh '''
                    echo "$DOCKER_PASS" | docker login -u "$DOCKER_USER" --password-stdin
                    docker build -t $REGISTRY/$IMAGE_NAME:$IMAGE_TAG .
                    docker push $REGISTRY/$IMAGE_NAME:$IMAGE_TAG
                    '''
                }
            }
        }

        stage('Deploy to Kubernetes') {
            steps {
                withCredentials([file(credentialsId: 'kubeconfig', variable: 'KUBECONFIG')]) {
                    sh '''
                    kubectl apply -f k8s/configmap.yaml
                    kubectl apply -f k8s/secret.yaml
                    kubectl apply -f k8s/deployment.yaml
                    kubectl apply -f k8s/service.yaml
                    '''
                }
            }
        }
    }
}
