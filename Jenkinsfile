pipeline {
    options {
        timeout(time: 1, unit: 'HOURS')
    }
    agent none
    environment {
        registryController = 'teambitflow/bitflow-controller'
        registryProxy = 'teambitflow/bitflow-api-proxy'
        registryCredential = 'dockerhub'
        controllerImage = '' // Empty variables must be declared here to allow passing an object between the stages.
        proxyImage = ''
    }
    stages {
        stage('Build & test') {
            agent {
                docker {
                    image 'teambitflow/golang-build:debian'
                    args '-v /tmp/go-mod-cache/debian:/go'
                }
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
                stage('Build & test api-proxy') {
                    steps {
                        dir ('bitflow-api-proxy') {
                            sh 'rm -f go.sum'
                            sh 'go clean -i -v ./...'
                            sh 'go install -v ./...'
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
                                    ${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=bitflow-controller -Dsonar.branch.name=$BRANCH_NAME \
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
                                    ${scannerHome}/bin/sonar-scanner -Dsonar.projectKey=bitflow-api-proxy -Dsonar.branch.name=$BRANCH_NAME \
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
            }
        }
        stage('Docker bitflow-api-proxy') {
            agent {
                docker {
                    image 'teambitflow/golang-build:alpine'
                    args '-v /tmp/go-mod-cache/alpine:/go -v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
            stages {
                stage('Prepare and build container') {
                    steps {
                        sh 'bitflow-api-proxy/build/native-build.sh'
                        script {
                            proxyImage = docker.build registryProxy + ':$BRANCH_NAME-build-$BUILD_NUMBER', '-f bitflow-api-proxy/build/native-prebuilt.Dockerfile bitflow-api-proxy/build'
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
                                proxyImage.push("build-$BUILD_NUMBER")
                                proxyImage.push("latest")
                            }
                        }
                    }
                }
            }
        }
        stage('Docker bitflow-controller') {
            agent {
                docker {
                    image 'teambitflow/golang-build:alpine'
                    args '-v /tmp/go-mod-cache/alpine:/go -v /var/run/docker.sock:/var/run/docker.sock'
                }
            }
            stages {
                stage('Prepare and build container') {
                    steps {
                        sh 'bitflow-controller/build/native-build.sh'
                        script {
                            controllerImage = docker.build registryController + ':$BRANCH_NAME-build-$BUILD_NUMBER', '-f bitflow-controller/build/native-prebuilt.Dockerfile bitflow-controller/build'
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
                                controllerImage.push("build-$BUILD_NUMBER")
                                controllerImage.push("latest")
                            }
                        }
                    }
                }
            }
        }
    }
    node {
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
}

