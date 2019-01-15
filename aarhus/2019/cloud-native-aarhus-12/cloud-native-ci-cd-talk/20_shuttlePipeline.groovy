def call(body) {
  properties([[$class: 'BuildDiscarderProperty', strategy: [$class: 'LogRotator', numToKeepStr: '25']], pipelineTriggers([[$class: 'BitBucketTrigger']])])

  node {
    // Checkout the project
    checkout scm

    // Init shuttle
    sh "curl -LO https://github.com/lunarway/shuttle/releases/download/\$(curl -s https://api.github.com/repos/lunarway/shuttle/releases/latest | grep tag_name | cut -d '\"' -f 4)/shuttle-linux-amd64"
    sh "chmod +x shuttle-linux-amd64"
    sh "mv shuttle-linux-amd64 shuttle"
    sh "./shuttle prepare -vc"

    // Load Pipeline defined in the plan
    def pipeline = load pwd() + '/.shuttle/plan/pipeline.groovy'
    pipeline.init()
  }
}
