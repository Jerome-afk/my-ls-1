package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"my_ls/flags"
	"my_ls/format"
	"my_ls/fileinfo"
)

var nonFlagArgs []string

const (
	ExitSuccess = 0
	ExitFailure = 1
	ExitInvalid = 2
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE]...\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "List information about FILEs (the current directory by default)")
		fmt.Fprintln(os.Stderr, "Options:")
		fmt.Fprintln(os.Stderr, "  -l    use a long listing format")
		fmt.Fprintln(os.Stderr, "  -R    list subdirectories recursively")
		fmt.Fprintln(os.Stderr, "  -a    do not hide entries starting with .")
		fmt.Fprintln(os.Stderr, "  -r    reverse order while sorting")
		fmt.Fprintln(os.Stderr, "  -t    sort by modification time, newest first")
	}

	// Parse flags
	var options flags.Flags

	// Process command line argument
	args := os.Args[1:]
	nonFlagArgs = []string{}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if len(arg) > 0 && arg[0] == '-' && len(arg) > 1 && arg[1] != '-' {
			for j := 1; j < len(arg); j++ {
				switch arg[j] {
				case 'l':
					options.LongFormat = true
				case 'R':
					options.Recursive = true
				case 'a':
					options.ShowHidden = true
				case 'r':
					options.ReverseSort = true
				case 't':
					options.SortByTime = true
				default:
					fmt.Fprintf(os.Stderr, "my-ls: invalid options -- '%c'\n", arg[j])
					flag.Usage()
					os.Exit(ExitInvalid)
				}
			}
		} else {
			// This is not a flag
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}

	// Get the target directory
	if len(nonFlagArgs) == 0 {
		// No directories
		err := processDirectory(".", options, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "my-ls: %v\n", err)
			os.Exit(ExitFailure)
		}
	} else {
		// Process each directory
		hasError := false
		for _, path := range nonFlagArgs {
			err := processDirectory(path, options, "")
			if err != nil {
				fmt.Fprintf(os.Stderr, "my-ls: %s: %v\n", path, err)
				hasError = true
			}
		}

		if hasError {
			os.Exit(ExitFailure)
		}
	}
}

func processDirectory(path string, options flags.Flags, prefix string) error {
	// File info
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	// If path is a file display it
	if !info.IsDir() {
		displaySingleFile(path, options)
		return nil
	}

	// Open the directory
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()

	// Read all files
	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	absPath, _ := filepath.Abs(path)
	if options.Recursive && prefix != "" {
		fmt.Printf("\n%s:\n", absPath)
	} else if len(nonFlagArgs) > 1 || options.Recursive {
		fmt.Printf("%s:\n", path)
	}

	// Sort entries
	entries = fileinfo.FilterAndSortEntries(entries, options)

	// Display entries
	if options.LongFormat {
		// Get total block size
		totalBlocks := fileinfo.GetTotalBlocks(entries, path)
		fmt.Printf("total %d\n", totalBlocks/2)

		for _, entry := range entries {
			entryPath := filepath.Join(path, entry.Name())
			formatting.DisplayLongFormat(entryPath, entry)
		}
	} else {
		// Simple format
		var names []string
		for _, entry := range entries {
			entryName := entry.Name()

			if entry.IsDir() {
				entryName = formatting.ColorString(entryName + "/", "blue", true)
			} else if entry.Mode()&os.ModeSymlink != 0 {
				entryName = formatting.ColorString(entryName, "cyan", true)
			} else if formatting.IsExecutable(filepath.Join(path, entry.Name())) {
				entryName = formatting.ColorString(entryName, "green", true)
			}
			names = append(names, entryName)
		}

		if len(names) > 0 {
			formatting.DisplayInColumns(names, path)
		}
	}

	// Recursively process subdirectories
	if options.Recursive {
		for _, entry := range entries {
			if entry.IsDir() {
				// Skip
				if entry.Name() == "." || entry.Name() == ".." {
					continue
				}

				subPath := filepath.Join(path, entry.Name())
				err := processDirectory(subPath, options, prefix+"  ")
				if err != nil {
					fmt.Fprintf(os.Stderr , "my-ls: %s: %v\n", subPath, err)
				}
			}
		}
	}

	return nil
}

func displaySingleFile(path string, options flags.Flags) {
	info, err := os.Lstat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "my-ls: %s: %v\n", path, err)
		return
	}

	if options.LongFormat {
		formatting.DisplayLongFormat(path, info)
	} else {
		name := info.Name()

		if info.IsDir() {
			name = formatting.ColorString(name+"/", "blue", true)
	    } else if info.Mode()&os.ModeSymlink != 0 {
			name = formatting.ColorString(name, "cyan", true)
	    } else if formatting.IsExecutable(path) {
			name = formatting.ColorString(name, "green", true)
	    }
	    fmt.Println(name)
	}
}