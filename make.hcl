default_goal = concat(rule.bins, rule.env_bins, ".lint")

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
  command = file("tidy.sh")
}
