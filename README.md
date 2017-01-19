# one-off [![Build Status](https://travis-ci.org/kkallday/one-off.svg?branch=master)](https://travis-ci.org/kkallday/one-off)

Generates a bash script to submit a local [concourse](https://www.concourse.ci) task. The bash
script will contain all pipeline vars for a given task and generate part of the `fly execute`
command for easy customization. `one-off` writes the resulting script to `stdout`.

`one-off` supports all versions up to 2.6.0.

## Usage
```
$ one-off -h
Usage of app:
  -fo string
    (optional) override path to fly. one-off uses "fly" in $PATH by default
  -j string
    name of job
  -p string
    name of pipeline
  -t string
    name of task
  -ta string
    concourse target alias
$ one-off -ta my-ci -p main -j integration-tests -t third-party-test > third-party-test-one-off.sh
$ cat third-part-test-one-off.sh
#!/bin/bash

export SOME_SECRET_PIPELINE_VAR1="foo"
export ANOTHER_ONE="keys for something"

fly -t my-ci execute --config=REPLACE/ME/PATH/TO/TASK \
    --inputs-from main/integration-tests
```

## Installation
Requires: `fly` [concourse cli](https://concourse.ci/fly-cli.html)
```
$ go get github.com/kkallday/one-off/one-off
```
## Testing

- To run tests you will need [ginkgo] (https://onsi.github.io/ginkgo) and [gomega] (https://onsi.github.io/gomega)
- Run `ginkgo -r` or `./bin/test` to run all tests
