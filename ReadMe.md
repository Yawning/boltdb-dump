# boltdb-dump #

This is a fork of boltdb-dump(https://github.com/chilts/boltdb-dump), with
changes to make it more generally useful.

Command to dump a human-readable BoltDB to stdout. This works with buckets
at any level, not just the top-level.

Note: The main reason for the changes was to be able to dump binary keys and
values.  If a key contains non-printable ASCII characters OR newlines, the
hexdecimal encoding will be displayed.  If a value contains non-printable
ASCII characters, output compatible with `hexdump -C` will be displayed.

## Install ##

```sh
go get -u github.com/yawning/boltdb-dump
```

## Usage ##

There are (currently) no options, nor anything fancy. Just pass the db file you want to dump:

```sh
boltdb-dump database.db
```

An example of a blog site with users and domains:

```sh
[users]
  chilts
    {"UserName":"chilts","Email":"andychilton@gmail.com"}
[domains]
  [chilts.org]
    [authors]
      andrew-chilton
    [posts]
      first-post
        {"PostName":"first-post","Title":"First Post","Content":"Hello, World!"}
  [blog.appsattic.com]
    [authors]
      andrew-chilton
    [posts]
```

## Author ##

[Andrew Chilton](https://chilts.org/) - [@andychilton](https://twitter.com/andychilton).
Yawning Angel

## License ##

MIT.

(Ends)
