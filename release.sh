#!/bin/zsh

# 添加帮助文档
usage() {
    echo "Usage: $0 <package_name|all> <version>"
    echo "Example: $0 all v1.0.0"
    echo "         $0 my-package v1.0.0"
    exit 1
}

[[ $# -ne 2 ]] && usage

release_pkg=${1}
version=${2}

echo "Release package: $release_pkg"
echo "Version: $version"

if [[ -z "$version" ]]; then
    echo "Error: Version cannot be empty" >&2
    exit 1
elif ! [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must match format v1.0.0" >&2
    exit 1
fi

function release_all_version() {
    local version=$1
    local file_name

    find . -maxdepth 1 -type d ! -name ".*" -printf "%f\n" | while read -r file_name; do
        [[ "$file_name" == "." ]] && continue
        release_single_version "$file_name" "$version"
    done
}

function release_single_version() {
    local file_name=$1
    local version=$2
    local _tag="${file_name}/${version}"

    echo "Creating tag: $_tag"
    if ! git tag -a "${_tag}" -m "Release ${version}"; then
        echo "Error: Failed to create tag" >&2
        return 1
    fi

    if ! git push origin "${_tag}"; then
        echo "Error: Failed to push tag" >&2
        return 1
    fi

    if ! git archive --format=zip --output="pkg-go-${file_name}-${version}.zip" "${_tag}"; then
        echo "Error: Failed to create archive" >&2
        return 1
    fi

    echo "Successfully released ${file_name} version ${version}"
}

if [[ "$release_pkg" == "all" ]]; then
    release_all_version "$version"
else
    if [[ -d "$release_pkg" ]]; then
        release_single_version "$release_pkg" "$version"
    else
        echo "Error: Unknown package name: $release_pkg" >&2
        exit 1
    fi
fi