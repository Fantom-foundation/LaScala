def call() {
    step([
        $class: 'ClassicUploadStep',
        credentialsId: 'aida-jenkins-service-account',
        bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
        pattern: "*.cpuprofile"
    ])

    step([
        $class: 'ClassicUploadStep',
        credentialsId: 'aida-jenkins-service-account',
        bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
        pattern: "*.memprofile"
    ])

    step([
        $class: 'ClassicUploadStep',
        credentialsId: 'aida-jenkins-service-account',
        bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
        pattern: "*.log"
    ])

    step([
        $class: 'ClassicUploadStep',
        credentialsId: 'aida-jenkins-service-account',
        bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
        pattern: "*.html"
    ])
}