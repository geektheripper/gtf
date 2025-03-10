# GTF

GeekTR's command line terraform utilities

## Modules

`gtf module publish module-name path/to/module/dir`

this command will:

1. create a orphan branch named `module-name`
2. add all files in `path/to/module/dir` to the new branch root directory
3. tag the new branch with auto-incremented version, e.g. `module-name/v0.0.4`
4. push the new tag with generated message includes `github.com/linolabx/terraform-modules-meta?ref=module-name/v0.0.4`

more options:

```bash
# if module path is not provided, it will use module name as path
gtf m p module-name

# auto increment different version part or specify version
gtf m p path-to-model --patch
gtf m p path-to-model --minor
gtf m p path-to-model --major
gtf m p path-to-model --version 0.0.4

# if you want to publish to a different remote repo
gtf m p path-to-model --remote=git@github.com:yourname/terraform-modules-meta.git

# LICENSE file will auto copy to the new module
# but you can disable the feature, or specify more files to copy
gtf m p path-to-model --no-license
gtf m p path-to-model --includes "README.md" --includes "another file"

# force publish (remove remote tag if exists)
gtf m p path-to-model --version 0.0.4 --force
```
