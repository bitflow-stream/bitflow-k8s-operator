pipeline {
    options {
        timeout(time: 1, unit: 'HOURS')
    }
    agent {
        docker {
            image 'teambitflow/golang-build'
            args '-v /root/.goroot:/go -v /var/run/docker.sock:/var/run/docker.sock'
        }
    }
    environment {
        registryController = 'teambitflow/bitflow-controller'
        registryProxy = 'teambitflow/bitflow-api-proxy'
        registryCredential = 'dockerhub'
        controllerImage = '' // Empty variables must be declared here to allow passing an object between the stages.
        proxyImage = ''
    }
    stages {
        stage('Git') {
            steps {
                script {
                    env.GIT_COMMITTER_EMAIL = sh(
                        script: "git --no-pager show -s --format='%ae'",
                        returnStdout: true
                        ).trim()
                }
            }
        }
        stage('Build & test controller') {
            steps {
                dir ('bitflow-controller') {
                    sh 'rm -f go.sum'
                    sh 'go clean -i -v ./...'
                    sh 'go install -v ./...'
                    sh 'go build -v -o ./build/_output/bin/bitflow-controller ./cmd/manager'
                    sh 'rm -rf reports && mkdir -p reports'
                    sh 'go test -v ./... -coverprofile=reports/test-coverage.txt 2>&1 | go-junit-report > reports/test.xml'
                    sh 'go vet ./... &> reports/vet.txt || true'
                    sh 'golint $(go list -f "{{.Dir}}" ./...) &> reports/lint.txt'
                }
            }
            post {
                always {
                    archiveArtifacts 'bitflow-controller/reports/*'
                    junit 'bitflow-controller/reports/test.xml'
                }
            }
        }
        stage('Build & test REST API proxy') {
            steps {
                dir ('bitflow-api-proxy') {
                    sh 'rm -f go.sum'
                    sh 'go clean -i -v ./...'
                    sh 'go install -v ./...'
                    sh 'go build -v -o ./build/_output/bin/bitflow-api-proxy .'
                    sh 'rm -rf reports && mkdir -p reports'
                    sh 'go test -v ./... -coverprofile=reports/test-coverage.txt 2>&1 | go-junit-report > reports/test.xml'
                    sh 'go vet ./... &> reports/vet.txt || true'
                    sh 'golint $(go list -f "{{.Dir}}" ./...) &> reports/lint.txt'
                }
            }
            post {
                always {
                    archiveArtifacts 'bitflow-api-proxy/reports/*'
                    junit 'bitflow-api-proxy/reports/test.xml'
                }
            }
        }
        stage('SonarQube') {
            steps {
                script {
                    // sonar-scanner which don't rely on JVM
                    def scannerHome = tool 'sonar-scanner-linux'
                    withSonarQubeEnv('CIT SonarQube') {
                        sh """
                            ${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=bitflow-controller \
                                -Dsonar.sources=bitflow-controller -Dsonar.tests=bitflow-controller \
                                -Dsonar.inclusions="**/*.go" -Dsonar.test.inclusions="**/*_test.go" \
                                -Dsonar.go.golint.reportPath=bitflow-controller/reports/lint.txt \
                                -Dsonar.go.govet.reportPaths=bitflow-controller/reports/vet.txt \
                                -Dsonar.go.coverage.reportPaths=bitflow-controller/reports/test-coverage.txt \
                                -Dsonar.test.reportPath=bitflow-controller/reports/test.xml
                        """
                    }
                    withSonarQubeEnv('CIT SonarQube') {
                        sh """
                            ${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=bitflow-api-proxy \
                                -Dsonar.sources=bitflow-api-proxy -Dsonar.tests=bitflow-api-proxy \
                                -Dsonar.inclusions="**/*.go" -Dsonar.test.inclusions="**/*_test.go" \
                                -Dsonar.go.golint.reportPath=bitflow-api-proxy/reports/lint.txt \
                                -Dsonar.go.govet.reportPaths=bitflow-api-proxy/reports/vet.txt \
                                -Dsonar.go.coverage.reportPaths=bitflow-api-proxy/reports/test-coverage.txt \
                                -Dsonar.test.reportPath=bitflow-api-proxy/reports/test.xml
                        """
                    }
                }
                timeout(time: 10, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true
                }
            }
        }
        stage('Docker build') {
            steps {
                script {
                    controllerImage = docker.build registryController + ':$BRANCH_NAME-build-$BUILD_NUMBER', '-f bitflow-controller/build/Dockerfile bitflow-controller'
                    proxyImage = docker.build registryProxy + ':$BRANCH_NAME-build-$BUILD_NUMBER', '-f bitflow-api-proxy/build/cached.Dockerfile bitflow-api-proxy'
                }
            }
        }
        stage('Docker push') {
            when {
                branch 'master'
            }
            steps {
                script {
                    docker.withRegistry('', registryCredential) {
                        aggregatorImage.push("build-$BUILD_NUMBER")
                        aggregatorImage.push("latest")
                        controllerImage.push("build-$BUILD_NUMBER")
                        controllerImage.push("latest")
                        collectorImage.push("build-$BUILD_NUMBER")
                        collectorImage.push("latest-amd64")
                        collectorImageARM32.push("build-$BUILD_NUMBER-arm32v7")
                        collectorImageARM32.push("latest-arm32v7")
                        collectorImageARM64.push("build-$BUILD_NUMBER-arm64v8")
                        collectorImageARM64.push("latest-arm64v8")
                        proxyImage.push("build-$BUILD_NUMBER")
                        proxyImage.push("latest")
                    }
                }
                withCredentials([
                  [   
                    $class: 'UsernamePasswordMultiBinding',
                    credentialsId: 'dockerhub',
                    usernameVariable: 'DOCKERUSER',
                    passwordVariable: 'DOCKERPASS'
                  ]   
                ]) {
                    // Dockerhub Login
                    sh '''#! /bin/bash
                    echo $DOCKERPASS | docker login -u $DOCKERUSER --password-stdin
                    ''' 
                    
                    sh "docker manifest create ${registryCollector}:latest ${registryCollector}:latest-amd64 ${registryCollector}:latest-arm32v7 ${registryCollector}:latest-arm64v8"
                    sh "docker manifest annotate ${registryCollector}:latest ${registryCollector}:latest-arm32v7 --os=linux --arch=arm --variant=v7"                                                                                 
                    sh "docker manifest annotate ${registryCollector}:latest ${registryCollector}:latest-arm64v8 --os=linux --arch=arm64 --variant=v8"
                    sh "docker manifest push --purge ${registryCollector}:latest"
                }
            }
        }
    }
    post {
        success {
            withSonarQubeEnv('CIT SonarQube') {
                slackSend channel: '#jenkins-builds-all', color: 'good',
                    message: "Build ${env.JOB_NAME} ${env.BUILD_NUMBER} was successful (<${env.BUILD_URL}|Open Jenkins>) (<${env.SONAR_HOST_URL}|Open SonarQube>)"
            }
        }
        failure {
            slackSend channel: '#jenkins-builds-all', color: 'danger',
                message: "Build ${env.JOB_NAME} ${env.BUILD_NUMBER} failed (<${env.BUILD_URL}|Open Jenkins>)"
        }
        fixed {
            withSonarQubeEnv('CIT SonarQube') {
                slackSend channel: '#jenkins-builds', color: 'good',
                    message: "Thanks to ${env.GIT_COMMITTER_EMAIL}, build ${env.JOB_NAME} ${env.BUILD_NUMBER} was successful (<${env.BUILD_URL}|Open Jenkins>) (<${env.SONAR_HOST_URL}|Open SonarQube>)"
            }
        }
        regression {
            slackSend channel: '#jenkins-builds', color: 'danger',
                message: "What have you done ${env.GIT_COMMITTER_EMAIL}? Build ${env.JOB_NAME} ${env.BUILD_NUMBER} failed (<${env.BUILD_URL}|Open Jenkins>)"
        }
    }
}

