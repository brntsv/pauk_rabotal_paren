# Go скрипт «pauk_rabotal»

Небольшая утилита на Go, которая глобально слушает клавиатуру и проигрывает `pauk_rabotal.mp3` по хоткею. Нажатие `Esc` завершает программу.

## Требования

- macOS или Windows.
- Go 1.26 или новее.
- Файл `pauk_rabotal.mp3` в корне проекта или рядом с собранным бинарником.

Зависимостей Go-модулей нет.

Платформенные детали:

- macOS: клавиатура слушается через CoreGraphics/cgo, звук запускается через штатный `/usr/bin/afplay`. Нужно разрешение Accessibility для терминала или собранного приложения.
- Windows: клавиатура слушается через WinAPI `WH_KEYBOARD_LL`, звук проигрывается через WinMM/MCI.
- Linux: пока не поддержан. Глобальные хоткеи зависят от X11, Wayland или evdev, поэтому нужен отдельный backend и понятные требования к правам.

Хоткеи:

- macOS: правый `Command` проигрывает звук.
- Windows: правый `Alt` проигрывает звук.
- `Esc` завершает программу.

## Запуск без исполняемого файла

Этот вариант похож на `cargo run`: Go соберет временный бинарник сам, но отдельный файл в проекте не появится.

```bash
cd /Users/brntsv/projects/go/pauk_rabotal_paren
go run ./cmd/white_punk
```

## Сборка исполняемого файла

Так можно получить постоянный бинарник примерно как в Rust-проекте после `cargo build`, только путь задаем явно.

macOS:

```bash
cd /Users/brntsv/projects/go/pauk_rabotal_paren
mkdir -p target/debug
go build -o target/debug/white_punk ./cmd/white_punk
./target/debug/white_punk
```

Бинарник из `target/debug` можно запускать и абсолютным путем из другой директории: он найдет `pauk_rabotal.mp3` в корне проекта.

Windows PowerShell:

```powershell
cd C:\path\to\pauk_rabotal_paren
New-Item -ItemType Directory -Force target\debug | Out-Null
go build -o target\debug\white_punk.exe .\cmd\white_punk
.\target\debug\white_punk.exe
```

## Проверка

```bash
go test ./...
```

## Кросс-сборка

Собрать Windows-бинарник с macOS можно без cgo:

```bash
mkdir -p target/debug
GOOS=windows GOARCH=amd64 go build -o target/debug/white_punk.exe ./cmd/white_punk
```

Для macOS-сборки используется cgo, поэтому собирать ее проще на macOS:

```bash
mkdir -p target/debug
go build -o target/debug/white_punk ./cmd/white_punk
```

## Структура

```text
cmd/white_punk/     entrypoint бинарника
internal/app/       общий сценарий запуска и обработка событий
internal/hotkey/    платформенные глобальные хоткеи
internal/sound/     платформенное воспроизведение звука
```
