//go:build !darwin && !windows

package sound

import "log"

func Play(soundPath string) {
	log.Printf("воспроизведение звука на этой платформе не реализовано: %s", soundPath)
}
