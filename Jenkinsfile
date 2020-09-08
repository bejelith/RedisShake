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

    }

    stages {

      stage('Prepare') {
        steps {
            sh 'echo $PATH'
            sh 'echo STARTING BUILD'
        }
      }

      stage('Build') {
        environment {
          BUILD_ID = VersionNumber([versionNumberString : '${BUILD_YEAR}.${BUILD_MONTH}.${BUILD_DAY}.' + env.CHANGE_ID + '.' + env.BUILD_NUMBER])
        }
        steps {
          deleteDir()
          checkout scm

          withCredentials([file(credentialsId: 'DPMBUILD-ARTIF-CREDENTIALS', variable: 'dpmbuildfile')]) {
            sh 'cp $dpmbuildfile ~/.docker/config.json'
          }

		  sh 'docker build -t docker-dev-artifactory.workday.com/dpm/redisshake:${BUILD_ID} .'
        }
      }

    }
}
