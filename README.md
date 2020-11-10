# BerrFuzz
BerrFuzz is a binary fuzzer. This is different from something like AFL in that is language agnostic and takes a black box approach, meaning there is no source code required.

It is currently in its infancy stage!

## Example Use
```bash
berrfuzz -i <inputFile> -n 3 -g
```
This is a simple generation of 3 files that have randomly mutated bytes based on the input file.

## File Generator Mode

