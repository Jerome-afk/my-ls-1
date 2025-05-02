package formatting

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"my_ls/fileinfo"
)

func DisplayLongFormat(path string, info os.FileInfo) {
	modeStr := GetModeString(info.Mode())
	linkCount := fileinfo.GetNumberOfLinks(info)

	// Owner and group
	ownerName := GetOwnerName(fileinfo.GetFileOwner(info))
	groupName := GetGroupName(fileinfo.GetFileGroup(info))

	size := info.Size()
	modTime := info.ModTime()
	timeStr := FormatTime(modTime)

	name := info.Name()
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err == nil {
			name = fmt.Sprintf("%s -> %s", name, target)
		}
	}

	fmt.Printf("%s %2d %s %6d %s %s\n",
		modeStr, linkCount, ownerName, groupName, size, timeStr, name)
}

func GetModeString(mode os.FileMode) string {
	var typeChar string
	switch {
	case mode&os.ModeDir != 0:
		typeChar = "d"
	case mode&os.ModeSymlink != 0:
		typeChar = "l"
	case mode&os.ModeNamedPipe != 0:
		typeChar = "p"
	case mode&os.ModeSocket != 0:
		typeChar = "s"
	case mode&os.ModeDevice != 0:
		typeChar = "b"
	case mode&os.ModeCharDevice != 0:
		typeChar = "c"
	default:
		typeChar = "-"
	}

	// Convert permissions
	perm := mode.Perm()
	permStr := ""

	// Owner permissions
	permStr += formatPermissionsBits((perm & 0o700) >> 6)

	// Group permission
	permStr += formatPermissionsBits((perm & 0o070) >> 3)

	// Others permissions
	permStr += formatPermissionsBits(perm & 0o007)

	return typeChar + permStr
}

func formatPermissionsBits(bits os.FileMode) string {
	result := ""
	result += bit(bits, 4, "r")
	result += bit(bits, 2, "w")
	result += bit(bits, 1, "x")
	return result
}

func bit(bits os.FileMode, mask uint32, char string) string {
	if bits&os.FileMode(mask) != 0 {
		return char
	}
	return "-"
}

func GetOwnerName(uid uint32) string {
	u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	if err != nil {
		return strconv.FormatUint(uint64(uid), 10)
	}
	return u.Username
}

// GetGroupName gets the group name from GID
func GetGroupName(gid uint32) string {
	g, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))
	if err != nil {
		return strconv.FormatUint(uint64(gid), 10)
	}
	return g.Name
}

// FormatTime formats the modification time according to ls rules
func FormatTime(t time.Time) string {
	now := time.Now()
	sixMonthsAgo := now.AddDate(0, -6, 0)

	if t.After(sixMonthsAgo) && t.Before(now) {
		// Recent files: "May 16 10:49"
		return t.Format("Jan _2 15:04")
	} else {
		// Older files: "May 16  2022"
		return t.Format("Jan _2  2006")
	}
}

// IsExecutable checks if a file is executable
func IsExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode()&0o111 != 0
}

// IsSymlink checks if a file is a symbolic link
func IsSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}

// GetAbsolutePath gets the absolute path of a file
func GetAbsolutePath(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absPath
}

// FormatFileName formats a filename with appropriate decorations
func FormatFileName(path string) string {
	info, err := os.Lstat(path)
	if err != nil {
		return filepath.Base(path)
	}

	name := info.Name()

	// Add trailing slash for directories
	if info.IsDir() {
		name += "/"
	}

	return name
}

// PadString pads a string to a given width
func PadString(s string, width int) string {
	padding := width - len(s)
	if padding <= 0 {
		return s
	}
	return s + strings.Repeat(" ", padding)
}

// ColorString returns a colored string using ANSI escape codes
// This is used to simulate the colored output of ls
// If the terminal doesn't support colors, this will still look okay
func ColorString(s, color string, bold bool) string {
	// ANSI color codes
	var colorCode string
	switch color {
	case "black":
		colorCode = "30"
	case "red":
		colorCode = "31"
	case "green":
		colorCode = "32"
	case "yellow":
		colorCode = "33"
	case "blue":
		colorCode = "34"
	case "magenta":
		colorCode = "35"
	case "cyan":
		colorCode = "36"
	case "white":
		colorCode = "37"
	default:
		return s // No color code available
	}

	// Bold text if requested
	if bold {
		colorCode = "1;" + colorCode
	}

	// Return colored string using ANSI escape codes
	return fmt.Sprintf("\033[%sm%s\033[0m", colorCode, s)
}

// DisplayInColumns displays a list of names in columns
func DisplayInColumns(names []string, path string) {
	// Calculate terminal width (use a reasonable default)
	termWidth := 80

	// Find the longest name to determine column width
	maxLen := 0
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
	}

	// Add padding between columns
	colWidth := maxLen + 2

	// Calculate number of columns that can fit
	numCols := termWidth / colWidth
	if numCols < 1 {
		numCols = 1
	}

	// Calculate number of rows needed
	numRows := (len(names) + numCols - 1) / numCols

	// Create a 2D array for the display
	display := make([][]string, numRows)
	for i := 0; i < numRows; i++ {
		display[i] = make([]string, numCols)
		for j := 0; j < numCols; j++ {
			idx := i + j*numRows
			if idx < len(names) {
				display[i][j] = names[idx]
			}
		}
	}

	// Display the names in columns
	for i := 0; i < numRows; i++ {
		for j := 0; j < numCols; j++ {
			if j < len(display[i]) && display[i][j] != "" {
				fmt.Print(PadString(display[i][j], colWidth))
			}
		}
		fmt.Println()
	}
}
