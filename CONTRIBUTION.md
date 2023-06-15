Print the code of lines:

```shell
git ls-files | xargs cloc
```

## Development between different remote git repository

I highly recommend you to use [the GitHub official git CLI](https://github.com/cli/cli).
And make sure you use the same email address on all the git repository platform.

First, clone the upstream git repository, and fork it.

```shell
# login if this is your first time to use it
# gh auth login
gh repo clone linuxsuren/api-testing
cd api-testing
gh repo fork --remote
git remote -v
```

then, add the corresponding git remote repository to your local repository. Such as:

```shell
# Please make sure clone the code from https://github.com/LinuxSuRen/api-testing
git remote add gitee https://gitee.com/linuxsuren/api-testing
git checkout -b gitee-master
git branch --set-upstream-to gitee/master

# sync it to GitHub once finished your work

git checkout master
git reset --hard gitee/master
git checkout -b gitee/feat/xxx
git push -u origin gitee/feat/xxx
gh pr create # or using other ways to create a PR

# sync to the downstream once the GitHub PR was merged
git fetch --all
git checkout gitee-master
git reset --hard gitee/master
git rebase origin/master
git rebase gitee/master
git push -u gitee master -f
```
