# Goani

Goani is a webserver for hosting static content.

## Configuration

`config.toml`:

```toml
Port = 8080
Folders = ['./logs']
```

Folders take a list of paths and will serve them with the name of the first folder.

## Building

```sh
$ docker build -t goani .
$ docker run -v $PWD:/go/src/goani goani go build -o /go/src/goani/build/goani /go/src/goani
```

## License

Goani is licensed under GNU GPL 3.0
