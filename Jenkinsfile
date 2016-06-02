node {
  stage 'build'
  checkout scm
  sh 'debuild -us -uc'
  step ([artifacts: "**/*.deb", fingerprint: true])
}
