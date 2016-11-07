# one-off

Generates a bash script to submit a local [concourse](https://www.concourse.ci) task. The bash 
script will contain all pipeline vars for a given task and generate part of the `fly execute` 
command for easy customization.

## Usage
```
$ one-off -h
Usage of app:
  -j string
    name of job
  -out string
    (optional) directory to write one off script
  -p string
    name of pipeline
  -t string
    name of task
  -ta string
    concourse target alias
$ one-off -ta my-ci -p main -j integration-tests -t third-party-test
```

## Installation
Requires: `fly` [concourse cli](https://concourse.ci/fly-cli.html)
```
$ go get github.com/kkallday/one-off/one-off
```
