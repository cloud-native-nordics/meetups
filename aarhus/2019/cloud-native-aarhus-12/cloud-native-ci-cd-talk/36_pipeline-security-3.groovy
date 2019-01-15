// lunar-way-subscription-service/.shuttle/plan/pipeline.groovy
//...
stage('Tests') {
  try {
    parallel 'project tests': {
      try {
        sh './shuttle -v run integration_test report=true'
      } finally {
        step([$class: 'JUnitResultArchiver', testResults: 'reports/*.xml', allowEmptyResults: true])
      }
    },
    'snyk security scan - code': {
      try {
        if(env.BRANCH_NAME == "dev"){
          sh "./shuttle -v run security_scan:code tag=${project.imageTag()} upload=true"
        } else {
          sh "./shuttle -v run security_scan:code tag=${project.imageTag()}"
        }
      } catch(err) {
        echo "Failed snyk security scan - continuing during evaluation phase. Error: $err"
      } finally {
        publishHTML([allowMissing: true, alwaysLinkToLastBuild: false, keepAll: false, reportDir: 'reports', reportFiles: 'snyk-code.html', reportName: 'snyk: code', reportTitles: ''])
      }
    },
    'snyk security scan - docker': {
      try {
        if(env.BRANCH_NAME == "dev"){
          sh "./shuttle -v run security_scan:docker tag=${project.imageTag()} upload=true"
        } else {
          sh "./shuttle -v run security_scan:docker tag=${project.imageTag()}"
        }
      } catch(err) {
        echo "Failed snyk security scan - continuing during evaluation phase. Error: $err"
      } finally {
        publishHTML([allowMissing: true, alwaysLinkToLastBuild: false, keepAll: false, reportDir: 'reports', reportFiles: 'snyk-docker.html', reportName: 'snyk: docker', reportTitles: ''])
      }
    }
  } finally { /* ... */ }
}