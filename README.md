
# gorcedit

[![Go Version](https://img.shields.io/github/go-mod/go-version/azhuravlev/gorcedit)](https://golang.org/)
[![License](https://img.shields.io/github/license/azhuravlev/gorcedit)](./LICENSE)

**gorcedit** is a simple command-line tool written in Go for editing Rails 5.2+ encrypted credentials.  
It is designed for **read/write access to strings encoded via Ruby Marshal format 4.8 only**.


---

## ✨ Features

- Edit encrypted `credentials.yml.enc` and similar files
- No Rails environment needed
- Pure Go implementation
- Supports editing via your favorite `$EDITOR`
- Debug mode for troubleshooting decryption/encryption issues

---

## 📦 Installation

    go get -u github.com/azhuravlev/gorcedit@latest

## Usage 

    gorcedit [--key KEY] [--keyfile KEYFILE] [--debug] [CREDSPATH]
