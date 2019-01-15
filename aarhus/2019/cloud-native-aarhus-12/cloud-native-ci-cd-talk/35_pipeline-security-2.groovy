// lunar-way-subscription-service/.shuttle/plan/pipeline.groovy
//...
stage('Tests') {
  try {
    parallel 'project tests': { /* project test spec */ },
    'snyk security scan': { try {
        sh "./shuttle -v run security_scan:code tag=${project.imageTag()} upload=true"
      } catch(err) {
        echo "Failed snyk security scan - continuing during evaluation phase. Error: $err"
      } finally {
        publishHTML([allowMissing: true, alwaysLinkToLastBuild: false, keepAll: false, reportDir: 'reports', reportFiles: 'snyk-code.html', reportName: 'snyk: code', reportTitles: ''])
      }
    },
    'sourceclear scan': {
      try {
        sh "./shuttle -v run sourceclear_scan tag=${project.imageTag()}"
      } catch(err) {
        echo "Failed sourceclear scan - continuing during evaluation phase. Error: $err"
      }
    },
    'aqua scan': {
      try {
        sh "./shuttle -v run aqua_scan tag=${project.imageTag()}"
      } catch(err) {
        echo "Failed aqua scan - continuing during evaluation phase. Error: $err"
      }
    }
  } finally { /* ... */ }
}