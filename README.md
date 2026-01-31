zar
===

<!-- pwsh -Command "badges.ps1" | -->
[![License](https://img.shields.io/badge/License-MIT-red)](https://github.com/hymkor/zar/blob/master/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/hymkor/zar.svg)](https://pkg.go.dev/github.com/hymkor/zar)
<!-- -->

`zar` is a zip clone with a tar-like interface.

Motivation
----------

- I often forget the `-r` option in the native zip command.
    - I prefer a tar-like interface, but my colleagues are not familiar with tar archives, so I have to use zip archives.
- I want to know whether filenames are encoded in UTF-8.
    - The fifth field of the output that `zar tvf` indicates the encoding of the filename.
        - `U` ... UTF8
        - `A` ... NonUTF8(ANSI)
- Are the files stored correctly in the archive?
  - `zar --md5    -tvf` shows the MD5    checksums of the files in the archive
  - `zar --sha1   -tvf` shows the sha1   checksums of the files in the archive
  - `zar --sha256 -tvf` shows the sha256 checksums of the files in the archive

```
$ zar --md5 -tvf zar-noversion-linux-amd64.zip
b52fa7c45800545279d3a8264484386b -rw-rw-rw-  694516 2021/10/23 00:28 U zar
```

Create an archive
-----------------

- `zar cvf ARCHIVE.zip [files...]`
- `zar cvf - [files...] > ARCHIVE.ZIP`

Filenames are always encoded in UTF-8.

### Option

- `--remove-files` removes files from disk after they are added to the archive.
- `-C DIRECTORY` changes the working directory

List an archive
----------------

- `zar tvf ARCHIVE.zip [files...]`  
- `zar tvf - [files...] < ARCHIVE.ZIP` \*

### Option

- `--md5` shows MD5 checksums of files in the archive

Extract an archive:
-------------------

- `zar xvf ARCHIVE.zip [files...]`  
- `zar xvf - [files...] < ARCHIVE.ZIP` \*

---

- Reading the archive from STDIN is slower because ZIP stores the file list at the end, so a temporary file is needed for random access.
- By golang's "archive/zip" package.

![](./demo.png)

Install
-------

### Manual Installation

Download the binary package from [Releases](https://github.com/hymkor/zar/releases) and extract the executable.

<!-- pwsh -Command "readme-install.ps1" | -->

### Use [eget] installer (cross-platform)

```sh
brew install eget        # Unix-like systems
# or
scoop install eget       # Windows

cd (YOUR-BIN-DIRECTORY)
eget hymkor/zar
```

[eget]: https://github.com/zyedidia/eget

### Use [scoop]-installer (Windows only)

```
scoop install https://raw.githubusercontent.com/hymkor/zar/master/zar.json
```

or

```
scoop bucket add hymkor https://github.com/hymkor/scoop-bucket
scoop install zar
```

[scoop]: https://scoop.sh/

### Use "go install" (requires Go toolchain)

```
go install github.com/hymkor/zar@latest
```

Note: `go install` places the executable in `$HOME/go/bin` or `$GOPATH/bin`, so you need to add this directory to your `$PATH` to run `zar`.
<!-- -->

License
-------

- MIT License

Acknowledgement
---------------

- [eqawasm (Ethar A Qawasmeh)](https://github.com/eqawasm)

Author
------

- [HAYAMA Kaoru (a.k.a zetamatta)](https://github.com/hymkor)
