# Smart Scan

Smart scan is an application built in golang that is used to cleanup duplicate files and
folders.

# Usage

Used to cleanup duplicate files in the directory. Currently supports MAC OS only. Will
update support for other versions later on.

To run the dev version: `go run main.go <absolute path to clean up>`

# Responsibilities

1. Remove directories mapped in env variables.
2. Removes duplicate files in provided directories.

# Rules

1. Removed files and folders and moved to trash; not permanently deleted.