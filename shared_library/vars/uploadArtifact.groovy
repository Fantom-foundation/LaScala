def call(pattern) {
    step([
        $class: 'ClassicUploadStep',
        credentialsId: 'aida-jenkins-service-account',
        bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
        pattern: "${pattern}"
    ])
}