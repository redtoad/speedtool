# speedtool

This is a small CLI tool to measure connection speed via the Fast.com API (by way of [ddo/fast](https://github.com/ddo/go-fast)).

    $ ./speedtool
    Initialise measurements...
    52854.784000
    55623.680000
    56535.722667
    57024.512000
    57232.588800
    57349.461333
    57493.796571
    57638.912000
    57682.602667
    57733.939200
    Done.
    {"ts": "2021-11-07T22:14:37+0100", "speed": 56716.999924}
    
## Development

New releases are done with [GoReleaser](https://goreleaser.com/):

    git tag v...
    goreleaser release --rm-dist   

## Changes

### 0.1.3 (2021-11-08)

* Write all status updates to stderr.

### 0.1.2 (2021-11-07)

* Add option `--output` to choose JSON or CSV as output format
* Errors will also be measured (with `null`)

### 0.1.1 (2021-11-07)

* Add option `-q` to suppress status updates
### 0.1.0 (2021-11-07)

* Initial release.