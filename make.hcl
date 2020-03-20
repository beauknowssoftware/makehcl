default_goal = "build"

command build {
  dependencies = concat(rule.bins, rule.env_bins, ".lint")
}

import {
  file = "build/tools.hcl"
}

import {
  file = "build/shell.hcl"
}

import {
  file = "build/common_vars.hcl"
}

import {
  file = "build/prebuild.hcl"
}

import {
  file = "build/local_exec.hcl"
}

import {
  file = "build/cross_platform.hcl"
}

import {
  file = "build/install.hcl"
}

command clean {
  command = "git clean -f -fdX"
}

command tidy {
  dependencies = ["tools/hclfmt"]
  command      = file("build/tidy.sh")
}
