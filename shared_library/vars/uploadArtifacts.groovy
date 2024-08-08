/**
* uploadArtifacts takes list of artifacts and stores them in google cloud storage bucket
* under path described in the pattern. It also archives the given artifacts on jenkins controller VM
*
* @param artifacts string[] list of artifact filenames
*/
def call(artifacts) {
    artifacts.each() {
        step([
            $class: 'ClassicUploadStep',
            credentialsId: 'aida-jenkins-service-account',
            bucket: "gs://aida-jenkins-artifacts/$JOB_NAME/$BUILD_NUMBER",
            pattern: it
        ])
    }

    archiveArtifacts artifacts: artifacts.join(',')
}