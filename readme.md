
# 🜏 Gaghiel — Remote Management Tool

  

**Gaghiel** is a lightweight Go client–server system for remotely managing your own devices.

Simple, extensible, and built for experimentation — it can execute commands, return basic system info, and serve as a foundation for further functionality.

  

---

  

## ⚙️ Overview

  

-  **Client:** connects to the **server** via TCP.

-  **Server:** lists connected clients, retrieves system info, and executes remote commands.

- Designed for **legitimate administration** of devices you own or control.

  

---

  

## ✨ Key Features

  

- 🔁 Auto-reconnect from client to server

- 🧠 `INFO` — returns username, hostname, OS, and IP

- 💻 `EXEC` — runs a command on the client and returns combined stdout/stderr

- 📦 Lightweight single binaries (no dependencies)

  

---

  

## 🚀 Quick Start (Local Test)

  

### 1️⃣ Build

```bash

go  build  -o  server  ./server

go  build  -o  client.exe  ./client  # or: go build -o client ./client (Linux/macOS)

```

  

### 2️⃣ Run the Server

```bash

./server

```

  

### 3️⃣ Run the Client

```bash

./client  127.0.0.1  1337

```

  

### 4️⃣ Interact via the Server Console

```bash

list  # show connected clients

exec <cmd> # execute command on selected client

info  # get system information

```

  ---

### 🧩 Usage Examples

```bash

info

# → Displays client username, hostname, OS, and IP

  

exec  mkdir  C:\test

# → Runs a Windows command remotely

```
---
  

## 🔒 Security & Responsible Use

  

>  **Warning**

> Use this tool **only** on systems you own or have explicit permission to manage.

  

- This project currently **does not include authentication or encryption**.

- Operate it **only within private and trusted networks**.

- Intended strictly for **educational purposes**.

  

---

  

## 📄 License

  

This project is licensed under the **MIT License**.

See the [LICENSE](LICENSE) file for full details.