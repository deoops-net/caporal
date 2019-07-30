docker {
  repo = "registry.cn-beijing.aliyuncs.com/deoops/caporal"
  tag = "{{release.prod}}"
  auth = {
    username = "_args.u"
    password = "_args.p"
  }
}

var "release" {
  prod = "v0.1.0"
  release = "v0.1.0"
}