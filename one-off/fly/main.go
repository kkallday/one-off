package main

import "fmt"

func main() {
	fmt.Println(`---
groups:
- name: some-pipeline-name
  jobs:
  - some-job
resources:
- name: some-resource
  type: git
jobs:
- name: some-job
  plan:
  - aggregate:
    - get: some-resource
  - task: some-task
    file: /path/to/task.yml
    params:
      random-non-param:
      - something-non-param
      VAR1: something
      VAR2: "another random var"
      VAR3: "maybe a cert value"`)
}
