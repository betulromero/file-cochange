# File Co-change

## Motivation
Often in a large codebase it can be very hard to understand the hidden structure between source files. There is the obvious structure of packages, libraries, etc. but there are also more subtle connections that only people with experience in the codebase know of. There are many times where a change in one part of the codebase requires a change elsewhere.

## Output

With this tool, you can find the files that are related with the current modified files in a git repository. Running the command in a git repository will score the files based on how likely they will need to be changed.

## Usage
You can simply download the source and run `go build ./cmd/cochanged` to generate the binary. All then you have to do is run the binary in a git repository.