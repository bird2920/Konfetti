package scanner

import (
	"os"
	"path/filepath"
	"strings"
)

// ScanDirs scans the provided directories for files with specified extensions.
// It returns a slice of file paths that match the given extensions, and a slice of error messages for any access errors encountered.
func ScanDirs(paths []string, extensions []string) ([]string, []string) {
	configFiles := make([]string, 0)
	errors := make([]string, 0)

	for _, path := range paths {
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				errors = append(errors, p+": "+err.Error())
				return nil
			}
			if info == nil {
				return nil
			}
			if !info.IsDir() {
				for _, ext := range extensions {
					if strings.HasSuffix(strings.ToLower(info.Name()), ext) {
						configFiles = append(configFiles, p)
						break
					}
				}
			}
			return nil
		})
	}
	return configFiles, errors
}

// ScanDirs scans the provided directories for files with specified extensions.
// It returns a slice of file paths that match the given extensions.
//// Parameters:
// - paths: A slice of directory paths to scan.
// - extensions: A slice of file extensions to look for (e.g., ".json", ".yaml").
//// Returns:
// - A slice of strings containing the paths of the found configuration files.
//// Example usage:
//   files := scanner.ScanDirs([]string{"/etc", "/home/user/.config"}, []string{".json", ".yaml"})
//   fmt.Println("Found config files:", files)
//// This function uses filepath.Walk to traverse directories and checks each file's extension.
// It collects paths of files that match the specified extensions and returns them as a slice.
// It handles errors gracefully and skips directories that cannot be accessed.
