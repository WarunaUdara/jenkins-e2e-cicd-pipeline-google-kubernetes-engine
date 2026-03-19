pipeline {
    
    agent any 
    
    environment {
        // DockerHub configuration
        DOCKERHUB_USERNAME = 'warunaudara'
        DOCKERHUB_REPO = "${DOCKERHUB_USERNAME}/todo-app-go"
        IMAGE_TAG = "${BUILD_NUMBER}"
        
        // GitHub credentials ID (to be created in Jenkins)
        // NOTE: Create this credential in Jenkins: Manage Jenkins → Credentials
        // Type: Username with password
        // Username: Your GitHub username
        // Password: GitHub Personal Access Token with 'repo' scope
        GITHUB_CREDENTIALS_ID = 'github-credentials'
        
        // Repository URLs
        APP_REPO = 'https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-google-kubernetes-engine.git'
        MANIFEST_REPO = 'https://github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-manifests.git'
    }
    
    stages {
        
        stage('Checkout Application Code'){
           steps {
                echo '==================== Stage 1: Checkout Application Code ===================='
                git credentialsId: "${GITHUB_CREDENTIALS_ID}", 
                    url: "${APP_REPO}",
                    branch: 'main'
                echo 'Application code checked out successfully'
           }
        }

        stage('Build Docker Image'){
            steps{
                script{
                    echo '==================== Stage 2: Build Docker Image ===================='
                    sh """
                    echo 'Building Docker image: ${DOCKERHUB_REPO}:${IMAGE_TAG}'
                    docker build -t ${DOCKERHUB_REPO}:${IMAGE_TAG} .
                    docker tag ${DOCKERHUB_REPO}:${IMAGE_TAG} ${DOCKERHUB_REPO}:latest
                    echo 'Docker image built successfully'
                    """
                }
            }
        }

        stage('Push Docker Image to DockerHub'){
           steps{
                script{
                    echo '==================== Stage 3: Push Docker Image ===================='
                    withCredentials([usernamePassword(
                        credentialsId: 'dockerhub-credentials',
                        usernameVariable: 'DOCKER_USER',
                        passwordVariable: 'DOCKER_PASS'
                    )]) {
                        sh """
                        echo 'Logging into DockerHub'
                        echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin
                        echo 'Pushing Docker image'
                        docker push ${DOCKERHUB_REPO}:${IMAGE_TAG}
                        docker push ${DOCKERHUB_REPO}:latest
                        docker logout
                        echo 'Image pushed successfully: ${DOCKERHUB_REPO}:${IMAGE_TAG}'
                        """
                    }
                }
            }
        }
        
        stage('Checkout Kubernetes Manifests'){
            steps {
                echo '==================== Stage 4: Checkout K8s Manifests ===================='
                // Clean workspace and checkout manifest repo
                cleanWs()
                git credentialsId: "${GITHUB_CREDENTIALS_ID}", 
                    url: "${MANIFEST_REPO}",
                    branch: 'main'
                echo 'Kubernetes manifests checked out successfully'
            }
        }
        
        stage('Update Manifest and Push to Git'){
            steps {
                script{
                    echo '==================== Stage 5: Update Manifest & Push ===================='
                    withCredentials([usernamePassword(
                        credentialsId: "${GITHUB_CREDENTIALS_ID}", 
                        passwordVariable: 'GIT_PASSWORD', 
                        usernameVariable: 'GIT_USERNAME'
                    )]) {
                        sh """
                        # Configure git
                        git config user.email "jenkins@cicd.local"
                        git config user.name "Jenkins CI/CD"
                        
                        # Show current deployment manifest
                        echo "Current deployment.yaml:"
                        cat deployment.yaml
                        
                        # Update image tag in deployment.yaml
                        # Replace the image line with new tag
                        sed -i "s|image: ${DOCKERHUB_REPO}:.*|image: ${DOCKERHUB_REPO}:${IMAGE_TAG}|g" deployment.yaml
                        
                        # Show updated manifest
                        echo "Updated deployment.yaml:"
                        cat deployment.yaml
                        
                        # Commit and push changes
                        git add deployment.yaml
                        git commit -m "Update image tag to ${IMAGE_TAG} | Jenkins Build #${BUILD_NUMBER}"
                        
                        # Push using credentials
                        git push https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/WarunaUdara/jenkins-e2e-cicd-pipeline-manifests.git HEAD:main
                        
                        echo 'Manifest updated and pushed to Git successfully'
                        """
                    }
                }
            }
        }
    }
    
    post {
        success {
            echo '==================== Pipeline Completed Successfully ===================='
            echo "Docker Image: ${DOCKERHUB_REPO}:${IMAGE_TAG}"
            echo 'ArgoCD will now sync and deploy to Kubernetes'
        }
        failure {
            echo '==================== Pipeline Failed ===================='
            echo 'Check the logs above for error details'
        }
        always {
            echo '==================== Cleaning Up ===================='
            // Clean up Docker images to save space
            sh """
            docker rmi ${DOCKERHUB_REPO}:${IMAGE_TAG} || true
            docker rmi ${DOCKERHUB_REPO}:latest || true
            """
        }
    }
}
