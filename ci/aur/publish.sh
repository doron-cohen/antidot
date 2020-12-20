#!/usr/bin/env bash

set -e

ROOT="$(dirname $(dirname $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )))"

export VERSION=$1
BINARY_PATH=$2
echo "Publishing to AUR as version ${VERSION}"

# Ensure the SSH private key ends with a new line
printf -- "${AUR_SSH_KEY}\n" > ~/.ssh/aur.key
chmod 600 ~/.ssh/aur.key
export GIT_SSH_COMMAND="ssh -i ~/.ssh/aur.key -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

cd ${ROOT}/ci/aur

rm -rf aur-package-repo
git clone aur@aur.archlinux.org:antidot-bin aur-package-repo 2>&1

export SHA256SUM=$(sha256sum ${BINARY_PATH} | awk '{ print $1 }')

CURRENT_PKGVER=$(cat aur-package-repo/.SRCINFO | grep pkgver | awk '{ print $3 }')
CURRENT_RELEASE=$(cat aur-package-repo/.SRCINFO | grep pkgrel | awk '{ print $3 }')

# AUR doesn't allow hiphens (-) so we replace it with underscores
# https://wiki.archlinux.org/index.php/PKGBUILD#pkgver
export PKGVER=${VERSION//-/_}

if [[ "${CURRENT_PKGVER}" == "${PKGVER}" ]]; then
    export RELEASE=$((CURRENT_RELEASE+1))
else
    export RELEASE=1
fi

export DESCRIPTION="Cleans up your \$HOME from those pesky dotfiles"
export BINARY_NAME="antidot"
export REPO_URL="https://github.com/doron-cohen/antidot"

ENVSUBST_VARS="\$PKGVER \$VERSION \$RELEASE \$SHA256SUM \$DESCRIPTION \$BINARY_NAME \$REPO_URL"

envsubst "$ENVSUBST_VARS" < templates/.SRCINFO > aur-package-repo/.SRCINFO
envsubst "$ENVSUBST_VARS" < templates/PKGBUILD > aur-package-repo/PKGBUILD

cd aur-package-repo
git config user.name "goreleaser"
git config user.email "goreleaserbot@doron.dev"
git --no-pager diff
git add .
if [ -z "$(git status --porcelain)" ]; then
  echo "No changes, skipping AUR repo update."
else
  git commit -m "Updated to version ${VERSION} release ${RELEASE}"
  git push origin master
fi

rm -f ~/.ssh/aur.key
