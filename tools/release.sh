# Serulian language release tool
#
# USAGE: Run from the *root directory* of the `compiler` repository:
# ./tools/release.sh v1.0.0 {github-api-token}
#
# NOTE: The `corelib`, `serulian-langserver` and `homebrew-serulian` repositories must
# be up-to-date and checked out under the same directory as `compiler` for this script
# to function.
#
# REQUIREMENTS:
#  - git
#  - go
#  - make
#  - hub (Github CLI tool)
#  - github_changelog_generator

check_repo() {
    echo "Checking $1 for any uncommited changes"
    if [[ -e "ls ../$1" ]] ; then
    echo "Could not find checkout of $1"
    exit 1
    fi

    cd ../$1
    if [[ -e "git diff master" ]] ; then
    echo "Found changes in $1"
    exit 1
    fi
}

check_dep() {
  hash $1 2>/dev/null || { echo >&2 "Could not find required utility $1"; exit 1; }
}

cut_release() {
  cd ../$1
  git checkout master
  git pull

  attach=""
  if $2; then
    echo ">> Building release for $1"
    make -f Makefile.release TOOLKIT_VERSION=$versionwithv
    for f in releases/*; do
      attach="$attach -a $f"
    done
  fi

  echo ">> Pushing version $versionwithv for $1"
  git tag -a $versionwithv -m "Version $versionwithv"
  git push origin $versionwithv

  # Generate the changelog.
  echo ">> Generating changelog and release notes for $1"
  releasenotes="$versionwithv.notes.md"

  github_changelog_generator -u serulian -p $1 --token=$token
  github_changelog_generator -u serulian -p $1 --token=$token --between-tags $versionwithv --no-unreleased  --header-label "Version $versionwithv" --output $releasenotes

  $editor CHANGELOG.md
  $editor $releasenotes

  git add CHANGELOG.md
  git commit -m "Update CHANGELOG.md"
  git push

  # Cut a release.
  hub release create $prerelease $versionwithv -f $releasenotes $attach $versionwithv
  rm $releasenotes

  echo ""
  echo ""
}

update_homebrew() {
  cd ../homebrew-serulian
  git checkout master
  git pull

  echo ">> Updating homebrew formula"

  previous_version_with_quotes=`grep "version \"*\"" serulian.rb | awk '{print $2}'`
  previous_version=${previous_version_with_quotes%\"}
  previous_version=${previous_version#\"}

  sed -i '' "s/$previous_version/$versionwithv/" serulian.rb
  git add serulian.rb
  git commit -m "Update to toolkit version $versionwithv"
  git push
}

check_dep "git"
check_dep "github_changelog_generator"
check_dep "make"
check_dep "hub"
check_dep "go"
check_dep "sed"
check_dep "grep"
check_dep "awk"

version="$1"
if [[ -z "$version" ]] ; then
  echo "Missing version"
  exit 1
fi

versionwithv="$version"
if [[ $version != "v"* ]]; then
  versionwithv="v$version"
fi

token="$2"
if [[ -z "$token" ]] ; then
  echo "Missing GitHub token"
  exit 1
fi

prerelease=""
if [[ $version == *"+"* ]]; then
  prerelease="-p "
fi

editor=${EDITOR:-emacs}

check_repo "corelib"
check_repo "compiler"
check_repo "serulian-langserver"
check_repo "homebrew-serulian"

cut_release "corelib" false
cut_release "compiler" true
cut_release "serulian-langserver" true
update_homebrew