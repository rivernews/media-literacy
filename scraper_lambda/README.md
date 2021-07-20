# Development: How to run golang programs

## Install golang environment

Install [`g`](https://github.com/stefanmaric/g) first:

```sh
curl -sSL https://git.io/g-install | sh -s -- zsh
```

Initialize go module and install dependencies

```
cd into golang src dir
go mod init github.com/rivernews/media-literacy
go get -u golang.org/x/net/html/charset
go get -u github.com/aws/aws-lambda-go/lambda
```

Note that the `-u` [is for including all child dependencies](https://blogs.halodoc.io/go-modules-implementation/).

# Reference

- [Running AWS Lambda in Golang](https://docs.aws.amazon.com/lambda/latest/dg/golang-handler.html)
