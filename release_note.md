Release Notes
=============

- Fix path traversal vulnerability during archive extraction. (#2, #3, thanks to @eqawasm)
v0.2.2
------
Mar.2,2022
- `-c` supports wildcard

v0.2.1
------
Feb.19,2022

- Change repository owner's name

v0.2.0
------
Feb.1, 2022

- Support new options: --sha1 and --sha256
- `--remove-files -cvf` change message `remove` to `rm` or `rmdir`
- `-c` and `-x` show prompt before overwriting files.
- Fix: `zar --remove-files cvf` was not invalid command-line
- `-c` and `-x` strip `X:/` from paths
- `--remove-files`: reduce memory consumption

v0.1.0
------
Jan.16, 2022

- The first release
