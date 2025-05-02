# My-LS: A Go Implementation of the Unix ls Command

This project implements a Go version of the Unix `ls` command with support for several common flags. The implementation is designed to be modular, robust, and closely match the behavior of the standard ls command.

## Features

The command supports the following flags:
- `-l`: Use a long listing format showing permissions, ownership, size, and date
- `-R`: List subdirectories recursively
- `-a`: Show all files (including hidden files starting with '.')
- `-r`: Reverse the order of the sort
- `-t`: Sort by modification time (newest first)

Flags can be combined (e.g., `-la`, `-Rrt`) for more complex operations.

## Implementation Details

The code is organized into several packages:
- `main.go`: Entry point and command-line processing
- `fileinfo`: File information handling, sorting, and filtering
- `formatting`: Display formatting and visual styling
- `flags`: Flag definitions and parsing

## Examples

Basic listing:
```
go run main.go
```

Showing all files including hidden ones:
```
go run main.go -a
```

Long format listing:
```
go run main.go -l
```

Recursive listing of subdirectories:
```
go run main.go -R
```

Complex combinations:
```
go run main.go -laRt
```

Multiple directories:
```
go run main.go dir1 dir2
```

## Testing

The project includes a comprehensive test script `test_script.sh` that verifies all supported flags and combinations.

## Error Handling

The implementation includes robust error handling with appropriate exit codes:
- `0`: Success
- `1`: General failure/error
- `2`: Invalid command usage

## Author

Created by Jerome on May 2, 2025.