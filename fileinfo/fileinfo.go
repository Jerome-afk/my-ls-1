package fileinfo

import (
        "os"
        "path/filepath"
        "sort"
        "syscall"

        "my_ls/flags"
)

// FilterAndSortEntries filters and sorts directory entries based on provided flags
func FilterAndSortEntries(entries []os.FileInfo, options flags.Flags) []os.FileInfo {
        var filtered []os.FileInfo

        // Safeguard against nil entries
        if entries == nil {
                return filtered
        }

        // Filter hidden files if -a is not set
        for _, entry := range entries {
                // Safeguard against nil entries
                if entry == nil {
                        continue
                }
                
                name := entry.Name()
                if len(name) > 0 && (options.ShowHidden || name[0] != '.') {
                        filtered = append(filtered, entry)
                }
        }

        // If no entries to sort, return early
        if len(filtered) == 0 {
                return filtered
        }

        // Always sort alphabetically first (this is what Unix ls does)
        sort.Slice(filtered, func(i, j int) bool {
                // Safeguard against index errors (shouldn't happen, but just in case)
                if i >= len(filtered) || j >= len(filtered) {
                        return false
                }
                return filtered[i].Name() < filtered[j].Name()
        })

        // Apply time sorting if requested
        if options.SortByTime && len(filtered) > 1 {
                // Create a stable sort function for modification time
                sort.SliceStable(filtered, func(i, j int) bool {
                        // Safeguard against index errors
                        if i >= len(filtered) || j >= len(filtered) {
                                return false
                        }
                        return filtered[i].ModTime().After(filtered[j].ModTime())
                })
        }

        // Apply reverse order if requested
        if options.ReverseSort && len(filtered) > 1 {
                // Reverse the slice
                for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
                        filtered[i], filtered[j] = filtered[j], filtered[i]
                }
        }

        return filtered
}

// GetTotalBlocks calculates the total number of blocks used by files in a directory
func GetTotalBlocks(entries []os.FileInfo, dirPath string) int64 {
        var total int64
        
        // Safeguard against nil entries
        if entries == nil {
                return 0
        }
        
        for _, entry := range entries {
                // Safeguard against nil entries
                if entry == nil {
                        continue
                }
                
                // Use Lstat to get file information without following symlinks
                fullPath := filepath.Join(dirPath, entry.Name())
                info, err := os.Lstat(fullPath)
                if err != nil {
                        continue
                }

                // Get system-specific information
                sys := info.Sys()
                if sys == nil {
                        continue
                }
                
                if stat, ok := sys.(*syscall.Stat_t); ok {
                        // Add block count (512-byte blocks)
                        total += stat.Blocks
                }
        }
        return total
}

// GetFileMode gets the file mode including permission bits
func GetFileMode(info os.FileInfo) os.FileMode {
        if info == nil {
                return 0
        }
        return info.Mode()
}

// GetFileOwner gets the owner UID of a file
func GetFileOwner(info os.FileInfo) uint32 {
        if info == nil {
                return 0
        }
        
        sys := info.Sys()
        if sys == nil {
                return 0
        }
        
        if stat, ok := sys.(*syscall.Stat_t); ok {
                return stat.Uid
        }
        return 0
}

// GetFileGroup gets the group GID of a file
func GetFileGroup(info os.FileInfo) uint32 {
        if info == nil {
                return 0
        }
        
        sys := info.Sys()
        if sys == nil {
                return 0
        }
        
        if stat, ok := sys.(*syscall.Stat_t); ok {
                return stat.Gid
        }
        return 0
}

// GetFileSize gets the size of a file in bytes
func GetFileSize(info os.FileInfo) int64 {
        if info == nil {
                return 0
        }
        return info.Size()
}

// GetFileModTime gets the modification time of a file
func GetFileModTime(info os.FileInfo) int64 {
        if info == nil {
                return 0
        }
        return info.ModTime().Unix()
}

// GetFileBlocks gets the number of blocks allocated for a file
func GetFileBlocks(info os.FileInfo) int64 {
        if stat, ok := info.Sys().(*syscall.Stat_t); ok {
                return stat.Blocks
        }
        return 0
}

// GetNumberOfLinks gets the number of hard links to a file
func GetNumberOfLinks(info os.FileInfo) uint64 {
        if stat, ok := info.Sys().(*syscall.Stat_t); ok {
                return uint64(stat.Nlink)
        }
        return 0
}