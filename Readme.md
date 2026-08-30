# Trak CLI 🚀

> **Trak** is a local CLI that resolves structured learning templates and materializes organized, hands-on learning workspaces directly on your filesystem.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Registry](https://img.shields.io/badge/Registry-GitHub%20Raw-black?style=flat&logo=github)](https://github.com/ndk123-web/trak-registry)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

---

## 💡 Why Trak?

Learning a new programming language, framework, or DevOps tool often starts with messy folders, scattered notes, and fragmented tutorials.

**Trak solves this by instantly materializing an end-to-end, structured workspace on your machine:**
- 📂 **Step-by-step modular folders** (e.g., `01-basics`, `02-concurrency`, etc.)
- 📝 **Pre-configured boilerplate & exercises** with runnable code and readme files
- 🌐 **Dynamic Cloud Registry** that updates without requiring you to update the CLI binary
- 🪶 **Zero runtime dependencies** — fast, lightweight single binary written in Go

---

## 📦 Installation

### From Source (Go 1.22+)
```bash
go install github.com/ndk123-web/trak/cmd/cli@latest
```

Or clone and build locally:
```bash
git clone https://github.com/ndk123-web/trak.git
cd trak
go build -o trak ./cmd/cli
```

---

## 🚀 Usage

### 1. Initialize a Learning Workspace
Create a structured learning track in your current directory:
```bash
trak init lang/go
```

Or specify a custom destination path:
```bash
trak init lang/go --path=./my-learning-folder
```

### 2. Available Template Examples
```bash
# Programming Languages
trak init lang/go

# DevOps & Tools
trak init tool/jenkins
```

### 3. List Available Templates
```bash
trak list
```

---

## 🏗️ How It Works (Architecture)

```text
User CLI Command:  trak init lang/go
                           │
                           ▼
                 ┌──────────────────┐
                 │  Registry Client │ ──▶ Fetches registry.json from GitHub
                 └─────────┬────────┘
                           │
                    lang/go exists?
                           │
                           ▼
                 ┌──────────────────┐
                 │ Template Fetcher │ ──▶ Fetches templates/lang/go.json
                 └─────────┬────────┘
                           │
                           ▼
                 ┌──────────────────┐
                 │ Filesystem Engine│ ──▶ Recursively creates directories & files
                 └─────────┬────────┘
                           │
                           ▼
                 ┌──────────────────┐
                 │ Workspace Stamped│ ──▶ Writes trak.json metadata locally
                 └──────────────────┘
```

---

## 📁 Generated Workspace Structure

When you run `trak init lang/go`, Trak generates:

```text
go-workspace/
├── go.mod
├── README.md
├── trak.json                  # Local workspace metadata
├── 01-hello-world/
│   ├── main.go
│   └── README.md
├── 02-variables-and-types/
│   ├── main.go
│   └── README.md
├── 03-control-flow/
│   └── main.go
├── 04-structs-and-methods/
│   └── main.go
├── 05-interfaces/
│   └── main.go
└── 06-concurrency/
    └── main.go
```

---

## 🤝 Template Ecosystem

Templates are hosted openly at [trak-registry](https://github.com/ndk123-web/trak-registry). Anyone can submit new tracks, languages, and tools!

---

## 📜 License

MIT License. See [LICENSE](LICENSE) for details.