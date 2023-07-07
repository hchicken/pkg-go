#!/bin/zsh

version=${1}
release_pkg=${2}

echo "version",$version,"release_pkg",$release_pkg

if [ -z "$version" ]; then
  echo "version is empty" && exit 1
elif ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "the version does not match the format of v1.0.0" && exit 1
fi


## variable to add version
function release_all_version() {
    version=$1
    for file_name in $(ls -d */|awk -F/ '{print $1}'); do
      _tag="${file_name}/${version}"
      git tag -a "${_tag}" -m "add the ${version}"
      git push origin "${_tag}"
      git archive --format=zip --output=pkg-go-${file_name}-${version}.zip "${_tag}"
    done
}

## add single version
function release_single_version() {
    version=$1
    file_name=$2
    _tag="${file_name}/${version}"
    git tag -a "${_tag}" -m "add the ${version}"
    git push origin "${_tag}"
    git archive --format=zip --output=pkg-go-${file_name}-${version}.zip "${_tag}"
}

# By entering the package name determine whether all released
if [ -z "$release_pkg" ]; then
  release_all_version ${version}
else
  if [ -d "${release_pkg}" ]; then
    release_single_version ${version} ${release_pkg}
  else
    echo "unknown package name" || exit 1
  fi
fi
