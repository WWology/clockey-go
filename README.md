# Clockey Bot

A modular Discord bot built with Go using the disgo library and SQLite database.

## Features

- ✅ Slash commands support
- ✅ SQLite database integration
- ✅ Modular command structure
- ✅ Clean logging system
- ✅ Environment-based configuration

## Project Structure

```
clockey-go/
├── main.go                  # Entry point
├── config/                  # Configuration management
├── database/                # SQLite database & queries
├── commands/                # Slash commands (modular)
├── bot/                     # Bot instance management
└── utils/                   # Shared utilities
```

## Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Configure environment variables:**
   ```bash
   cp .env.example .env
   ```
   Edit `.env` and add your Discord bot token:
   ```
   DISCORD_TOKEN=your_bot_token_here
   ```

3. **Run the bot:**
   ```bash
   go run main.go
   ```

## Development

### Adding a New Command

1. Create a new file in `commands/` (e.g., `mycommand.go`)
2. Define the command structure:
   ```go
   func (h *Handler) myCommand() discord.ApplicationCommandCreate {
       return discord.SlashCommandCreate{
           Name:        "mycommand",
           Description: "Description of my command",
       }
   }
   ```
3. Add the handler:
   ```go
   func (h *Handler) handleMyCommand(event *events.ApplicationCommandInteractionCreate) {
       // Your command logic here
   }
   ```
4. Register it in `commands/handler.go`:
   - Add to `RegisterCommands()` commands slice
   - Add case in `onInteractionCreate()` switch

### Database Queries

Add new query functions in `database/queries.go`. Follow the existing pattern:
- Return structs for complex data
- Use prepared statements
- Handle `sql.ErrNoRows` appropriately
- Return descriptive errors

## Environment Variables

| Variable       | Required | Default            | Description                              |
|----------------|----------|--------------------|------------------------------------------|
| DISCORD_TOKEN  | Yes      | -                  | Your Discord bot token                   |
| DATABASE_PATH  | No       | `data/clockey.db`  | Path to SQLite database file             |
| GUILD_ID       | No       | -                  | Guild ID for development (faster updates)|
| DEBUG          | No       | `false`            | Enable debug logging                     |

## Tips

- **Development**: Set `GUILD_ID` for instant command updates in your test server
- **Production**: Leave `GUILD_ID` empty to register global commands (takes ~1 hour to propagate)
- **Database**: SQLite file is automatically created on first run

## License

MIT
