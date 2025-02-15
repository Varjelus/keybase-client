#! /bin/bash

# Given a build root, invokes package_binaries.sh and then lays out the
# resulting packages as they would be in a debian package server.
#
# Dependencies:
#   - regular Go setup for building the client
#   - dpkg-deb and fakeroot for building the .deb
#   - reprepro for writing the hierarchy
#
# Note that all of the above works just fine on Arch as well as Debian.

set -e -u -o pipefail

here="$(dirname "$BASH_SOURCE")"

build_root="${1:-}"
if [ -z "$build_root" ] ; then
  echo 'Usage:  ./layout_repo.sh <build_root>'
  exit 1
fi

version="$(cat "$build_root/VERSION")"
mode="$(cat "$build_root/MODE")"

name="$("$here/../../binary_name.sh" "$mode")"

repo_root="$build_root/deb_repo"

# Run the Debian packaging script on this build root.
"$here/package_binaries.sh" "$build_root"

# Write out the reprepro config. We could just check this in, but writing it
# here means that we're free to nuke the entire prod/linux dir (generally just
# to test this build), and it also means we can share the PGP fingerprint.
code_signing_fingerprint="$(cat "$here/../code_signing_fingerprint")"
mkdir -p "$repo_root/repo/conf"
cat > "$repo_root/repo/conf/distributions" << END
Codename: stable
Components: main
Architectures: i386 amd64
SignWith: $code_signing_fingerprint
END

for debian_arch in amd64 i386 ; do
  echo Creating the Debian repo hierarchy.
  # Create/update the Debian repository hierarchy. The layout of this is
  # described here: https://wiki.debian.org/RepositoryFormat We use the
  # reprepro tool to automatically generate this hierarchy and sign it, which
  # works on (at least) Debian or Arch.
  debfile="$(ls "$build_root/deb/$debian_arch"/*.deb)"  # keybase, kbstage, or kbdev
  # reprepro reads "prod/linux/deb/repo/conf/distributions", which includes the
  # PGP fingerprint of our code signing key.
  reprepro --basedir "$repo_root/repo" includedeb stable "$debfile"

  # Update the latest pointer.
  ln -sf "repo/pool/main/k/$name/${name}_${version}_${debian_arch}.deb" \
    "$repo_root/$name-latest-$debian_arch.deb"
done

# Because the reprepro hierarchy only contains one version of each architecture
# at a time, copy the imported packages into our own "all" directory.
mkdir -p "$repo_root/all"
cp "$repo_root"/repo/pool/main/*/"$name"/*.deb "$repo_root/all"
