# 📘 GopherPad_ 
- A secure personal notepad service written in Go — complete with user authentication, JWT-based access control, and a PostgreSQL backend for safely storing user secrets. Ideal for learning Go server design, authentication patterns, database integration, and deployment workflows.

# 🔐 Features
- JWT Authentication: Users register and log in to receive a signed JWT token. All protected endpoints require this token for access.

- User Registration: Secure signup with password hashing (bcrypt).

- Secret Note Storage: Authenticated users can manage their personal vault-notes, which are stored privately in the database and only accessible via their token.

# CRUD Operations for notes:

1. Create a new note

2. Read all personal notes

3. Update a specific note

4. Delete a specific note

# 🛠️ Tech Stack
- Go (Golang) — backend server and handlers

- PostgreSQL — persistent storage for users and notes

- JWT — secure stateless authentication

- Docker — to run the app in containers

- Splunk UF — logs collection and forwarding (planned)

- Kubernetes — for deployment (planned)

# 📎 Example Use Case
A user signs up, logs in, and gets a JWT token. They can now create and manage secure, encrypted notes like a vault — a lightweight personal notepad that’s authenticated and encrypted via backend validation.

# 💡 Inspiration
This project is a Go-native alternative to apps like Standard Notes or Bitwarden's secure notes — but made for learning how full-stack backend systems work in Go.
