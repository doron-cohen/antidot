# antidot :house: :small_orange_diamond: :boom:

![Pipeline](https://github.com/doron-cohen/antidot/workflows/Pipeline/badge.svg?branch=master)

Cleans up your `$HOME` from those pesky dotfiles.

### :construction: This is a work in progress! :construction:

## Intro

For years I stood by and saw how countless applications populate my home dir with dotfiles.

No more! `antidot` is a tool to automatically detect and remove dotfiles from `$HOME` without any risks. It will move files to more appropriate locations (based on [XDG base directory specifications](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)). It will also set environment variables, declare aliases and use symlinks to ensure apps can find their files.
