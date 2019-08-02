docker {
  repo = "deoops/caporal"
  tag = "{{release.release}}"
  auth = {
    username = "{{_args.u}}"
    password = "{{_args.p}}"
  }
}

var "release" {
  prod = "v0.1.0"
  // add mount support
  release = "v0.1.54"
}