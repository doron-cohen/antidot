# antidot :house: :small_orange_diamond: :boom:

![Pipeline](https://github.com/doron-cohen/antidot/workflows/Pipeline/badge.svg?branch=master)

Cleans up your `$HOME` from those pesky dotfiles.

## Migration from Pre 0.6.0 Versions

Please backup your environment variable and aliases files (in `$XDG_DATA_HOME/antidot/{env,alias}.*`). After version 0.6.0 antidot stores env exports and alias definitions in a JSON file and generates from it shell definition files.

## Intro

For years I stood by and saw how countless applications populate my home dir with dotfiles.

No more! `antidot` is a tool to automatically detect and remove dotfiles from `$HOME` without any risks. It will move files to more appropriate locations (based on [XDG base directory specifications](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html)). It will also set environment variables, declare aliases and use symlinks to ensure apps can find their files.

## Installation

### Arch Linux

```shell
yay -Sy antidot-bin
```

### Homebrew

```shell
brew install doron-cohen/tap/antidot
```

Go to the [releases](https://github.com/doron-cohen/antidot/releases) section and grab the one that fits your OS.

After installing run `antidot update` to download the latest rules file and you're all set!

## How does it Work?

Dotfiles pollution is a complex problem to solve. There are many approaches to solve this annoying issue and `antidot` is taking the safest one.

We maintain a rule for each dotfile which applies actions when the file is detected. The main goal is to move the files to the most appropriate location while keeping the application working as expected.

There are a few types of actions:

* **Migrate** - Move a file to a new location, optionally symlink the old location to the new one (this is for the hardest cases).
* **Delete** - Delete a file or a directory (only if it's empty).
* **Export** - Export an environment variable.
* **Alias** - Set an alias to a command.

## Examples

This is the rule for the Docker configuration directory:

```yaml
  - name: docker
    dotfile:
      name: .docker
      is_dir: true
    actions:
      - type: migrate
        source: ${HOME}/.docker
        dest: ${XDG_CONFIG_HOME}/docker
      - type: export
        key: DOCKER_CONFIG
        value: ${XDG_CONFIG_HOME}/docker
```

When running `antidot clean` we will be prompted about this directory:

```bash
❯ ./antidot clean
Rule docker:
  MOVE   /Users/doroncohen/.docker → /Users/doroncohen/.config/docker
  EXPORT DOCKER_CONFIG="${XDG_CONFIG_HOME}/docker"
? Apply rule docker? (y/N)
```

Answering yes will move the directory and write the environment variable to a file that can be easily sourced by the shell. Running `antidot init` will create a shell script that will do just that.

Adding `eval "$(antidot init)"` to your `.bashrc` or `.zshrc` will make sure you shell sessions will see these variables and aliases. In Fish the proper way is to run `antidot init | source`. You could add it to `$__fish_config_dir/conf.d/antidot.fish`.
