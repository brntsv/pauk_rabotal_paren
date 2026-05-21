//go:build darwin

package sound

import (
	"context"
	"log"
	"os/exec"
	"time"
)

func Play(soundPath string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	cmd := exec.CommandContext(ctx, "/usr/bin/afplay", soundPath)

	if err := cmd.Start(); err != nil {
		cancel()
		log.Printf("не удалось запустить afplay: %v", err)
		return
	}

	go func() {
		defer cancel()
		if err := cmd.Wait(); err != nil && ctx.Err() == nil {
			log.Printf("ошибка воспроизведения %s: %v", soundPath, err)
		}
	}()
}
