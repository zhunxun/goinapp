workflow "On Push" {
  on = "push"
  resolves = ["Build"]
}

action "Run tests" {
  uses = "actions/docker/cli@86ff551d26008267bb89ac11198ba7f1d807b699"
  args = "build --target=TEST ."
}

action "Build" {
  uses = "actions/docker/cli@86ff551d26008267bb89ac11198ba7f1d807b699"
  args = "build ."
  needs = ["Run tests"]
}
