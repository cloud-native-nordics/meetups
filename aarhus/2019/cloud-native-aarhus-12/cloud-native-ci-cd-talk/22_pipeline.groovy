def init() {
    try {
      stage('Generate Build Info'){
        project.buildInfo()
      }

      withCredentials([string(credentialsId: 'jenkins_ssh_key', variable: 'SSH_KEY')]) {
        stage('Build Image') {
          sh "./shuttle run build tag=${project.imageTag()}"
        }

        stage('Tests') {
          try {
            sh './shuttle -v run integration_test report=true'
          }
          finally {
            step([$class: 'JUnitResultArchiver', testResults: 'reports/*.xml', allowEmptyResults: true])
          }
        }
      }

      if (project.runCD()) {

        stage('Push Image') {
          sh "./shuttle run push tag=${project.imageTag()} env=$env.BRANCH_NAME"
        }

        stage('Deploy to Kubernetes') {
          sh "./shuttle run deploy env=$env.BRANCH_NAME tag=${project.imageTag()}"
        }
      }

      stage('Notify slack') {
        slackslim.sendSuccessful("${shuttle.service()}", "${project.repo()}", "${shuttle.squad()}", "${git.author()}", "${git.hash()}", "${git.message()}", "${project.imageTag()}")
      }
      stage('Clean up workspace') {
        step([$class: 'WsCleanup'])
      }
    }
    catch(exc){
      slackslim.sendFailure("${shuttle.service()}", "${project.repo()}", "${shuttle.squad()}", "${git.author()}", "${git.hash()}", "${git.message()}", "${project.imageTag()}", "$exc")
      step([$class: 'WsCleanup'])
      throw exc
    }
}

return this