# Ownershipper

Ownershipper generates best-effort [CODEOWNERS files](https://help.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners#example-of-a-codeowners-file). It assigns owners based on the numbers of commits for each commiter.

## Usage

```go
go build ownershipper
./ownershipper <path to git project> -numowners=<defaults to 1> -out=<file or stdout, defaults to file> -overwright=<true or false>
```

## Help

```Markdown
Usage:
  -numowners int
        max number of owners to assign (default 1)
  -out string
        output format is file or stdout (default "stdout")
  -overwrite
        whether to overwrite an existing CODEOWNERS file
```
