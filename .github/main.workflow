workflow "PR to master" {
  resolves = ["Run tests"]
  on = "pull_request"
}

action "Run tests" {
  uses = "actions/docker/cli@86ff551d26008267bb89ac11198ba7f1d807b699"
  args = "build --target=BUILD ."
}
