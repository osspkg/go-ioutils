# go-ioutils

[![Go Reference](https://pkg.go.dev/badge/go.osspkg.com/ioutils.svg)](https://pkg.go.dev/go.osspkg.com/ioutils)
[![Go Version](https://img.shields.io/github/go-mod/go-version/osspkg/go-ioutils)](https://github.com/osspkg/go-ioutils/blob/master/go.mod)
[![CI](https://github.com/osspkg/go-ioutils/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/osspkg/go-ioutils/actions/workflows/ci.yml)
[![License](https://img.shields.io/github/license/osspkg/go-ioutils)](LICENSE)

`go-ioutils` — библиотека утилит Go для копирования потоков, работы с файлами, кодирования данных, кэширования, пулов объектов и запуска команд оболочки.

Требуется Go 1.26 или новее.

## Установка

```sh
go get go.osspkg.com/ioutils@latest
```

## Пакеты

| Пакет | Назначение |
| --- | --- |
| [`ioutils`](https://pkg.go.dev/go.osspkg.com/ioutils) | Копирование потоков, `ReadAll` и поиск маркера ASCII EOF |
| [`cache`](https://pkg.go.dev/go.osspkg.com/ioutils/cache) | Потокобезопасный обобщённый кэш и параметры очистки |
| [`codec`](https://pkg.go.dev/go.osspkg.com/ioutils/codec) | Кодирование и декодирование YAML, JSON, TOML, XML и Unic |
| [`fs`](https://pkg.go.dev/go.osspkg.com/ioutils/fs) | Поиск, копирование и перезапись файлов, вычисление хешей |
| [`pool`](https://pkg.go.dev/go.osspkg.com/ioutils/pool) | Обёртка над `sync.Pool` и пул сбрасываемых срезов |
| [`shell`](https://pkg.go.dev/go.osspkg.com/ioutils/shell) | Запуск команд с контекстом, окружением, рабочим каталогом и настройкой вывода |

## Примеры

### Копирование потока

```go
import (
	"io"

	"go.osspkg.com/ioutils"
)

func copyStream(dst io.Writer, src io.Reader) error {
	_, err := ioutils.Copy(dst, src)
	return err
}
```

`Copy`, `CopyN`, `CopyB` и `Pipe` копируют данные до EOF. Параметр `size` у `CopyN` задаёт размер буфера, а не максимальное число копируемых байтов. `ReadAll` читает данные и закрывает переданный `io.ReadCloser`.

### Кодирование и декодирование JSON

```go
import "go.osspkg.com/ioutils/codec"

type Config struct {
	Endpoint string `json:"endpoint"`
	Retries  int    `json:"retries"`
}

func roundTrip() error {
	blob := &codec.BlobEncoder{Ext: codec.ExtJSON}
	if err := blob.Encode(Config{Endpoint: "https://service.example", Retries: 3}); err != nil {
		return err
	}

	var decoded Config
	if err := blob.Decode(&decoded); err != nil {
		return err
	}
	return nil
}
```

Для работы с файлами используйте `codec.FileEncoder("settings.yaml")`. Кодек выбирается по точному расширению файла.

### Использование кэша

```go
import (
	"fmt"

	"go.osspkg.com/ioutils/cache"
)

func example() {
	c := cache.New[string, string]()
	c.Set("language", "Go")

	if value, ok := c.Get("language"); ok {
		fmt.Println(value)
	}
}
```

### Запуск команды оболочки

```go
import (
	"context"
	"fmt"

	"go.osspkg.com/ioutils/shell"
)

func run(ctx context.Context) error {
	sh := shell.New()
	output, err := sh.Call(ctx, "echo hello")
	if err != nil {
		return err
	}
	fmt.Print(string(output))
	return nil
}
```

По умолчанию на Linux и macOS используется `/bin/sh`, на Windows — `%COMSPEC%` (или `cmd.exe`). Синтаксис команды должен соответствовать выбранной оболочке; пакет не преобразует команды между операционными системами. Не вставляйте непроверенные данные в строку команды.

## Разработка

Запускайте команды из корня репозитория:

```sh
go test ./...
go test -race ./...
go vet ./...
```

GitHub Actions запускает тесты и `go vet` на Linux, macOS и Windows. Задача CI на Ubuntu выполняет `make ci`: устанавливает `goppy`, запускает его настройку, проверку лицензий, линтер, тесты и сборку.

## Лицензия

Проект распространяется на условиях [лицензии BSD 3-Clause](LICENSE).
