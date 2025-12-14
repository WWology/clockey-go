# Copilot Instructions

This repository contains a Discord bot written in Go using:
- DisGo for Discord interactions
- sqlc with PostgreSQL for persistence

GitHub Copilot should follow the rules below when generating or modifying code.

---

## Project Responsibilities

The bot has two primary domains:

1. **Moderator Hours Tracking**
   - Track, store, and query moderator working hours
   - Commands may create, update, and summarize time-based records

2. **Prediction Leaderboard**
   - Allow members to predict game scores
   - Calculate and store prediction results
   - Maintain leaderboards based on accumulated scores

Copilot should keep these domains logically separated in code.

---

## Repository Structure (Authoritative)

Copilot MUST respect the existing structure:

- `app/`
  Contains all bot-related code

- `app/commands/`
  Contains all Discord application commands (slash commands, context menus, autocomplete, UI interactions)

- `database/`
  Contains sqlc queries, schemas, and generated Go code

- `test/`
  Contains unit tests

Do NOT invent new top-level folders unless explicitly requested.

---

## Discord Command Guidelines

- Use **slash commands**, **context menu commands**, **autocomplete**, and **UI components** (buttons, selects, modals)
- Command handlers may contain business logic directly
- Avoid introducing unnecessary service or domain layers
- Keep each command handler focused and readable

---

## Database & sqlc Usage

- Database: **PostgreSQL**
- Use sqlc-generated queries directly inside command handlers
- Use explicit transactions **only when necessary**
- Follow sqlc’s default naming conventions
- The database is the source of truth; avoid in-memory state

---

## Error Handling & Logging

- Handle errors inline where they occur
- Log errors using **`log/slog`**
- Do not panic for expected runtime errors
- Use `context.Context` consistently (commands, DB, background work)
- User-facing errors:
  - Use **ephemeral replies** for sensitive or user-specific errors
  - Use **public replies** when appropriate for shared context

---

## Concurrency & Performance

- Avoid goroutines unless there is a clear benefit
- Assume a **large Discord server (10k+ members)**
- Prefer stateless command handlers backed by the database
- Do not introduce shared mutable state without strong justification

---

## Configuration & Startup

- Load configuration into a struct at startup
- Maintain separate configurations for development and production
- Use a **single shared PostgreSQL connection pool**
- Never hardcode secrets, tokens, or DSNs

---

## Testing Guidelines

- Tests are nice to have, not mandatory
- When written:
  - Focus on unit tests for command handlers
  - Use Go’s standard `testing` package only
- Avoid over-mocking DisGo internals

---

## Go Coding Style

Copilot MUST:

- Follow `gofmt` and `go vet` conventions
- Prefer explicit error handling over panics
- Avoid unnecessary interfaces and abstractions
- Prefer clarity over cleverness
- Use short, descriptive function and variable names
- Add comments only for **complex logic**, not obvious code

---

## General Rules

- Generate idiomatic, production-ready Go code
- Do not invent APIs or DisGo features that do not exist
- Match existing code style before introducing new patterns
- Keep solutions simple, readable, and maintainable
