zar
===

zar is the zip-clone like tar-interface.

**Create an archive:**

- `zar cvf ARCHIVE.zip [files...]`
- `zar cvf - [files...] > ARCHIVE.ZIP`

The filename encoding is always set as UTF8.

**List an archive:**

- `zar tvf ARCHIVE.zip [files...]`  
- `zar tvf - [files...] < ARCHIVE.ZIP` \*

**Extract an archive:**

- `zar xvf ARCHIVE.zip [files...]`  
- `zar xvf - [files...] < ARCHIVE.ZIP` \*

---

- Reading the archive file from STDIN is slow because the ZIP file has the file list in the tail of the archive and we have to make a temporary file to access randomly.
- By golang's "archive/zip" package.
- The fifth field of the output that `zar tvf` indicates the encoding of the filename.
    - `U` ... UTF8
    - `A` ... NonUTF8(ANSI)

![](./demo.png)
