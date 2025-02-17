# Versioning and Release

**Important:** Before starting the release process, ensure that all [E2E tests](https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei/actions/workflows/test-e2e.yml) have passed.

## Create tag, and build images

1. Please fork <https://github.com/HuaweiCloudDeveloper/cluster-api-provider-huawei> and clone your own repository with e.g. `git clone git@github.com:YourGitHubUserName/cluster-api-provider-huawei.git`.
2. Add a git remote to the upstream project. `git remote add upstream git@github.com:HuaweiCloudDeveloper/cluster-api-provider-huawei.git`
3. If this is a major or minor release, create a new release branch and push to GitHub, otherwise switch to it, e.g. `git checkout release-0.0`.
4. If this is a major or minor release, update `metadata.yaml` by adding a new section with the version, and make a commit.
5. Update the release branch on the repository, e.g. `git push origin HEAD:release-0.0`. `origin` refers to the remote git reference to your fork.
6. Update the release branch on the repository, e.g. `git push upstream HEAD:release-0.0`. `upstream` refers to the upstream git reference.
7. Make sure your repo is clean by git standards.
8. Set environment variables which is the last release tag and `VERSION` which is the current release version, e.g. `export VERSION=v0.0.1`.
    * _**Note**_: the version MUST contain a `v` in front.
    * _**Note**_: you must have a gpg signing configured with git and registered with GitHub.

9. Create a tag `git tag -s -m $VERSION $VERSION`. `-s` flag is for GNU Privacy Guard (GPG) signing.
10. Make sure you have push permissions to the upstream CAPHW repo. Push tag you've just created (`git push <upstream-repo-remote> $VERSION`). Pushing this tag will kick off a GitHub Action that will create the release, attach the binaries and YAML templates to it, and build multi-architecture images: `ghcr.io/huaweiclouddeveloper/cluster-api-huawei-controller:$VERSION`.

## Verify and Publish the draft release

1. Verify that all the files below are attached to the drafted release:
    1. `infrastructure-components.yaml`
    2. `cluster-template.yaml`
    3. `metadata.yaml`
2. Update the release description to link to the CAPHW image.
3. Update the draft release's title and description, and publish the release.

## Support Horizon

| CAPHW Version | kubernetes Version   | huawei cloud provider Version |
|---------------|----------------------|-------------------------------|
| v0.0.1        | v1.32.0+             | -                             |
