#!groovy

pipeline {
    agent any

    options {
        buildDiscarder(logRotator(numToKeepStr:'10'))
        disableConcurrentBuilds()
        timeout(time: 30)
    }

    environment {
        BERKSHELF_PATH="/tmp/berkshelf-chef-parcel-omsdpm"
        JAVA_HOME='/usr/java/default/'
        BUILD_ID = VersionNumber([versionNumberString : '${BRANCH_NAME}.' + env.BUILD_NUMBER])
        BASE_IMAGE = "docker-dev-artifactory.workday.com/dpm/redisshake"
        IMAGE = "${BASE_IMAGE}:$env.BUILD_ID"
    }

    stages {

        stage('Prepare') {
            steps {
                sh 'echo $PATH'
                sh 'echo STARTING BUILD'
            }
        }

        stage('Build') {
            steps {
                VersionNumber([versionNumberString : '${env.BUILD_ID}'])
                deleteDir()
                checkout scm

                withCredentials([file(credentialsId: 'DPMBUILD-ARTIF-CREDENTIALS', variable: 'dpmbuildfile')]) {
                    sh 'cp $dpmbuildfile ~/.docker/config.json'
                }

                withCredentials([usernameColonPassword(credentialsId: 'DPMBUILD_ARTIF', variable: 'USERPASS')]) {
                    sh "docker build --build-arg goproxy=\"https://${USERPASS}@artifactory.workday.com/artifactory/api/go/go\" -t ${IMAGE} ."

                }
            }
        }
        stage('Publish'){
            steps{
                sh "docker push ${env.IMAGE}"
                sh "docker tag ${env.IMAGE} ${env.BASE_IMAGE}:latest"
                sh "docker push ${env.BASE_IMAGE}:latest"
            }
        }

    }
}