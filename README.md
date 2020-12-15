### bin

An online service that can help you easily install binary files.

It is mainly aimed at Golang, of course other types are also available.

> I am tired of writing installation guide for so many cli.
>
> I am trying to use this tool to solve it all at once

WIP:

- [ ] Custom binary file template
- [ ] Support different decompression formats
- [ ] Testing

### Usage

eg, I want to install the release file of [https://github.com/axetroy/dvm](https://github.com/axetroy/dvm)

You need to run the following commands in your operating systems

#### Install binary for Linux/MacOS

```bash
curl https://bin2.herokuapp.com/axetroy/dvm | bash
```

#### Install binary for Windows

```bash
iwr https://bin2.herokuapp.com/axetroy/dvm -useb | iex
```

### How it works?

Depending on the URL you visit with different tools, different content will be returned.

In PowerShell, it returns the ps1 file, in Unix systems it returns to Shell

### License

The [MIT License](LICENSE)
