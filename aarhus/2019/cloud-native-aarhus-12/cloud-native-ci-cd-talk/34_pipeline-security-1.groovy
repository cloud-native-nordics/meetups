// lunar-way-subscription-service/.shuttle/plan/pipeline.groovy
//...
stage('Tests') {
  try {
    sh './shuttle -v run integration_test report=true'
  }
  finally {
    step([$class: 'JUnitResultArchiver', testResults: 'reports/*.xml', allowEmptyResults: true])
  }
}
//...