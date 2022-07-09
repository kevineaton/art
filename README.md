# art

An experimental package for generative art built as a learning exercise. Not intended for use.

## Inspirations

Heavily inspired by Preslav Rachev's [Generative Art in Go](https://p5v.gumroad.com/l/generative-art-in-golang).

## Usage

Each approach is separated into a subpackage and then exposed through the CLI. To see the options, run

`./art --help`

Each command will take in its own configuration, preferably from the command line. Examples coming soon.

### Implemented Approaches

#### Transform

Uses a base image in the `./input/` directory and generates a new piece of art using colors from the original and generating shapes/ Works best with landscapes and images with a lot of different colors.
