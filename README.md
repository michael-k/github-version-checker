# github-version-checker

github-version-checker checks if a newer version (aka tag) exists on github.


## Configuration

Put your github token in `env.list`:

```conf
GITHUB_TOKEN=1a2b3c4d
```

## Run

```sh
docker build -t github-version-checker .
docker run --rm --env-file env.list github-version-checker [repoOwner] [repoName] [versionInUse]
```
