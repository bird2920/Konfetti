package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirs_Empty(t *testing.T) {
	files, errs := ScanDirs([]string{}, []string{".json"})
	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
	if len(errs) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errs))
	}
}

func TestScanDirs_NonExistentPath(t *testing.T) {
	files, errs := ScanDirs([]string{"/unlikely/path/that/does/not/exist"}, []string{".json"})
	if len(files) != 0 {
		t.Errorf("Expected 0 files for non-existent path, got %d", len(files))
	}
	if len(errs) == 0 {
		t.Errorf("Expected at least one error for non-existent path")
	}
}

func TestScanDirs_SimpleCase(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "test1.json")
	file2 := filepath.Join(dir, "test2.yaml")
	file3 := filepath.Join(dir, "not_config.txt")
	os.WriteFile(file1, []byte(`{"foo": "bar"}`), 0644)
	os.WriteFile(file2, []byte(`foo: bar`), 0644)
	os.WriteFile(file3, []byte(`not config`), 0644)

	files, errs := ScanDirs([]string{dir}, []string{".json", ".yaml"})
	if len(files) != 2 {
		t.Errorf("Expected 2 config files, got %d", len(files))
	}
	found := map[string]bool{file1: false, file2: false}
	for _, f := range files {
		if _, ok := found[f]; ok {
			found[f] = true
		}
	}
	for f, v := range found {
		if !v {
			t.Errorf("Expected to find file %s", f)
		}
	}
	if len(errs) != 0 {
		t.Errorf("Expected 0 errors, got %d", len(errs))
	}
}

func TestScanDirs_PermissionDenied(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "noaccess")
	os.Mkdir(sub, 0000)
	defer os.Chmod(sub, 0700) // restore so TempDir cleanup works
	files, errs := ScanDirs([]string{dir}, []string{".json"})
	// Should not panic or error, just skip
	if files == nil {
		t.Errorf("Expected files slice, got nil")
	}
	if len(errs) == 0 {
		t.Errorf("Expected at least one error for permission denied directory")
	}
}

func TestScanDirs_SymlinkLoop(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sub")
	os.Mkdir(sub, 0755)
	loop := filepath.Join(sub, "loop")
	os.Symlink(dir, loop) // symlink to parent
	file := filepath.Join(sub, "test.json")
	os.WriteFile(file, []byte(`{"foo": "bar"}`), 0644)
	files, _ := ScanDirs([]string{dir}, []string{".json"})
	found := false
	for _, f := range files {
		if f == file {
			found = true
		}
	}
	if !found {
		t.Errorf("Expected to find %s", file)
	}
	// Should not hang or panic due to symlink loop
	// Symlink loops may or may not produce errors depending on platform, so no strict error check here
}
