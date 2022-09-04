# Cleanup utility

The utility will keep up to **N** (defaults to 5) of the latest subdirectories and removes all the others.

Modification time will be used when determining the sort order of subdirectories.

## Usage examples

* `cleanup <PATH-TO-DIRECTORY> -l 5` keeps up to 5 latest subdirectories of the path provided
* `cleanup <PATH-TO-DIRECTORY>` keeps up to 5 (default value) subdirectories of the path provided
* `cleanup -v` displays the version information of the utility script
* `cleanup` - displays help text of the utility script

## Building

Build the binaries by running `./go-executable-build.bash` and check the **data** directory for resulting binaries. You
may change the defined platform as needed.

Default platforms are `linux/386`, `linux/amd64` and `darwin/amd64`.