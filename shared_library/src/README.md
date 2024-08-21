# Groovy source files

## Introduction
The aim for the future is to keep the code of the pipelines as clean as possible, so we opted for Jenkins shared
library feature. This allows us to create global methods which can be easily imported and used throughout LaScala 
and Jenkins. We use jenkins shared library in its most simple form.

## Basic usage and development
All methods of shared library are stored in the `shared_library/vars` directory.
These files are written in groovy and every file contains a method `call`
which will be initialized when calling a method from a pipeline. In the pipeline you call a method by its filename.

You can see a basic usage example in `utils/template.jenkinsfile` with `uploadArtifacts` method.

You can find more information in the official [documentation](https://www.jenkins.io/doc/book/pipeline/shared-libraries/).
