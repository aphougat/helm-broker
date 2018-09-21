#!/usr/bin/env groovy

def label = "kyma-${UUID.randomUUID().toString()}"
def application = "helm-broker"
def dockerPushRoot = "${env.DOCKER_REGISTRY}${params.PUSH_DIR}"
def dockerImageTag = params.APP_VERSION

echo """
********************************
Job started with the following parameters:
DOCKER_REGISTRY=${env.DOCKER_REGISTRY}
PUSH_DIR=${params.PUSH_DIR}
DOCKER_CREDENTIALS=${env.DOCKER_CREDENTIALS}
GIT_REVISION=${params.GIT_REVISION}
GIT_BRANCH=${params.GIT_BRANCH}
APP_VERSION=${params.APP_VERSION}
APP_FOLDER=${env.APP_FOLDER}
FULL_BUILD=${params.FULL_BUILD}
********************************
"""

podTemplate(label: label) {
    node(label) {
        try {
            timestamps {
                ansiColor('xterm') {
                    timeout(time: 20, unit: "MINUTES") {
                        stage("setup") {
                            checkout scm

                            if(dockerImageTag == ""){
                                error("No version for docker tag defined, please set APP_VERSION parameter for master branch or GIT_BRANCH parameter for any branch")
                            }

                            withCredentials([usernamePassword(credentialsId: env.DOCKER_CREDENTIALS, passwordVariable: 'pwd', usernameVariable: 'uname')]) {
                                sh "docker login -u $uname -p '$pwd' $env.DOCKER_REGISTRY"
                            }
                        }

                        stage("build and test $application") {
                            execute "./before-commit.sh ci"
                        }

                        stage("build image $application") {
                            dir(env.APP_FOLDER){
                                sh "cp ./broker deploy/broker/helm-broker"
                                sh "cp ./reposerver deploy/reposerver/reposerver"
                                sh "cp ./indexbuilder deploy/tools/indexbuilder"
                                sh "cp ./targz deploy/tools/targz"

                                sh "docker build -t ${dockerPushRoot}helm-broker:latest deploy/broker --label version=${dockerImageTag} --label component=${application}"
                                sh "docker build -t ${dockerPushRoot}helm-broker-reposerver:latest deploy/reposerver --label version=${dockerImageTag} --label component=${application}"
                                sh "docker build -t ${dockerPushRoot}helm-broker-tools:latest deploy/tools --label version=${dockerImageTag} --label component=${application}"
                            }
                        }

                        stage("push image $application") {
                            sh "docker tag ${dockerPushRoot}helm-broker:latest ${dockerPushRoot}helm-broker:${dockerImageTag}"
                            sh "docker push ${dockerPushRoot}helm-broker:${dockerImageTag}"

                            sh "docker tag ${dockerPushRoot}helm-broker-reposerver:latest ${dockerPushRoot}helm-broker-reposerver:${dockerImageTag}"
                            sh "docker push ${dockerPushRoot}helm-broker-reposerver:${dockerImageTag}"

                            sh "docker tag ${dockerPushRoot}helm-broker-tools:latest ${dockerPushRoot}helm-broker-tools:${dockerImageTag}"
                            sh "docker push ${dockerPushRoot}helm-broker-tools:${dockerImageTag}"

                            if (params.GIT_BRANCH == 'master') {
                                sh "docker push ${dockerPushRoot}helm-broker:latest"
                                sh "docker push ${dockerPushRoot}helm-broker-reposerver:latest"
                                sh "docker push ${dockerPushRoot}helm-broker-tools:latest"
                            }

                        }
                    }
                }
            }
        } catch (ex) {
            echo "Got exception: ${ex}"
            currentBuild.result = "FAILURE"
            def body = "${currentBuild.currentResult} ${env.JOB_NAME}${env.BUILD_DISPLAY_NAME}: on branch: ${params.GIT_BRANCH}. See details: ${env.BUILD_URL}"
            emailext body: body, recipientProviders: [[$class: 'DevelopersRecipientProvider'], [$class: 'CulpritsRecipientProvider'], [$class: 'RequesterRecipientProvider']], subject: "${currentBuild.currentResult}: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]'"
        }
    }
}

def execute(command, envs = '') {
    def buildpack = 'golang-buildpack:0.0.8'
    def envText = envs=='' ? '' : "--env $envs"
    workDir = pwd()
    sh "docker run --rm -v $workDir:/go/src/github.com/kyma-project/kyma/ -w /go/src/github.com/kyma-project/kyma/$env.APP_FOLDER $envText ${env.DOCKER_REGISTRY}$buildpack /bin/bash -c '$command'"
}
