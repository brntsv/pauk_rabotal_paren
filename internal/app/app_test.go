package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveSoundPathFromWorkingDirectory(t *testing.T) {
	tempDir := t.TempDir()
	soundPath := filepath.Join(tempDir, soundFile)
	if err := os.WriteFile(soundPath, []byte("fake mp3"), 0o644); err != nil {
		t.Fatalf("write sound file: %v", err)
	}

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	got, err := resolveSoundPath(soundFile)
	if err != nil {
		t.Fatalf("resolve sound path: %v", err)
	}

	gotInfo, err := os.Stat(got)
	if err != nil {
		t.Fatalf("stat resolved sound path: %v", err)
	}

	wantInfo, err := os.Stat(soundPath)
	if err != nil {
		t.Fatalf("stat expected sound path: %v", err)
	}

	if !os.SameFile(gotInfo, wantInfo) {
		t.Fatalf("resolve sound path = %q, want %q", got, soundPath)
	}
}

func TestResolveSoundPathMissing(t *testing.T) {
	tempDir := t.TempDir()

	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	})

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}

	if _, err := resolveSoundPath("missing.mp3"); err == nil {
		t.Fatal("resolve sound path succeeded for missing file")
	}
}

func TestResolveSoundPathFromTargetDebugProjectRoot(t *testing.T) {
	projectDir := t.TempDir()
	executableDir := filepath.Join(projectDir, "target", "debug")
	if err := os.MkdirAll(executableDir, 0o755); err != nil {
		t.Fatalf("create executable directory: %v", err)
	}

	soundPath := filepath.Join(projectDir, soundFile)
	if err := os.WriteFile(soundPath, []byte("fake mp3"), 0o644); err != nil {
		t.Fatalf("write sound file: %v", err)
	}

	workingDir := t.TempDir()
	executablePath := filepath.Join(executableDir, "white_punk")

	got, err := resolveSoundPathFrom(soundFile, workingDir, executablePath)
	if err != nil {
		t.Fatalf("resolve sound path: %v", err)
	}

	gotInfo, err := os.Stat(got)
	if err != nil {
		t.Fatalf("stat resolved sound path: %v", err)
	}

	wantInfo, err := os.Stat(soundPath)
	if err != nil {
		t.Fatalf("stat expected sound path: %v", err)
	}

	if !os.SameFile(gotInfo, wantInfo) {
		t.Fatalf("resolve sound path = %q, want %q", got, soundPath)
	}
}
