#! /bin/bash

# Builds the keybase binary and packages it into two ".deb" files, one for i386
# and one for amd64. The argument to this script is the output directory of a
# build_binaries.sh build. The package files are created there, in their
# respective architecture folders.
#
# Usage:
#   ./package_binaries.sh <build_root>

set -e -u -o pipefail

here="$(dirname "$BASH_SOURCE")"

build_root="${1:-}"
if [ -z "$build_root" ] ; then
  echo 'Usage:  ./package_binaries.sh <build_root>'
  exit 1
fi

version="$(cat "$build_root/VERSION")"
mode="$(cat "$build_root/MODE")"

name="$("$here/../../binary_name.sh" "$mode")"

if [ "$mode" = "production" ] ; then
  repo_url="http://dist.keybase.io/linux/deb/repo"
elif [ "$mode" = "prerelease" ] ; then
  repo_url="http://s3.amazonaws.com/prerelease.keybase.io/deb"
elif [ "$mode" = "staging" ] ; then
  # Note: This doesn't exist yet. But we need to be distinct from the
  # production URL, because we're moving to a model where we build a clean repo
  # every time, rather than adding to an existing one. (For S3 compatibility.)
  repo_url="http://dist.keybase.io/linux/deb_staging/repo"
else
  # We don't actually publish devel builds. This URL is a dream within a dream.
  repo_url="http://dist.keybase.io/linux/deb_devel/repo"
fi

build_one_architecture() {
  echo "Making .deb package for $debian_arch."
  dest="$build_root/deb/$debian_arch"
  mkdir -p "$dest/build/DEBIAN"

  # Copy the entire filesystem layout, binaries and all, into the debian build
  # folder. TODO: Something less wasteful of disk space?
  cp -r "$build_root"/binaries/"$debian_arch"/* "$dest/build"

  # Installed-Size is a required field in the control file. Without it Ubuntu
  # users will see warnings.
  size="$(du --summarize --block-size=1024 "$dest" | awk '{print $1}')"

  cat "$here/control.template" \
    | sed "s/@@NAME@@/$name/" \
    | sed "s/@@VERSION@@/$version/" \
    | sed "s/@@ARCHITECTURE@@/$debian_arch/" \
    | sed "s/@@SIZE@@/$size/" \
    > "$dest/build/DEBIAN/control"
  postinst_file="$dest/build/DEBIAN/postinst"
  cat "$here/postinst.template" \
    | sed "s|@@REPO_URL@@|$repo_url|" \
    | sed "s/@@SOURCE_LIST_NAME@@/$name/" \
    > "$postinst_file"
  chmod 755 "$postinst_file"

  fakeroot dpkg-deb --build "$dest/build" "$dest/$name-$version-$debian_arch.deb"
}

export debian_arch=amd64
build_one_architecture

export debian_arch=i386
build_one_architecture
