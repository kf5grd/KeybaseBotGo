# Config

The bot's config is held in a file called `config.json` in the same folder as the bot. If no config file is found, a default/empty config will be written. At the very least, you should fill in the `botOwner` field with your username. The rest of the configuration can be done by interacting with the bot.

## Import the Config Package
```go
import(
    keybot/config
)
```

## Default Config
```json
{
  "botOwner": ""
}
```

## Create Default Config
```go
const ConfigFile = "config.json"
c := config.ConfigJSON{}

// Create default config if none exists
if _, err := os.Stat(ConfigFile); os.IsNotExist(err) {
    c.Write(ConfigFile)
}
```

## Read Config
```go
c := config.ConfigJSON{}

c.Read("config.json")
```

## Write Config
After making changes to a `ConfigJSON` interface, you can call the `Write()` method, passing just the filename.

```go
c := config.ConfigJSON{}

c.BotOwner = "dxb"
c.Write("config.json")
```

# Interacting With the Keybase API

## Import the API Package
```go
import(
    "keybot/api"
)
```

## Chat

### Send a Message To a User
```go
user := api.Channel{Name: "dxb"}
user.SendMessage("Hello, dxb!")
```

### Send a Message To a Team
```go
team := api.Channel{Name: "crbot.public", Channel: "bots", IsTeam: true}
team.SendMessage("Hello everyone!")
```

## Teams

### Add Members
```go
members := map[string]string{
    // "username": "role"
    "cagingroyals": "admin",
}

t := api.Team{Name: "crbot.public"}
teamAdd := t.AddMembers(members)
if teamAdd.Error.Message != "" {
    fmt.Println(teamAdd.Error.Message)
} else {
    fmt.Println("Users successfully added to team '%s'.", t.Name)
}
```

### Get List of Members
The `ListMembers()` method returns a `map[string]string` with the key being the username and value being the role. This is the same as the input for `AddMembers()`.  
  
The following code snippet will send a message to user `dxb` with a multiline code block containing all users in the team `crbot.public`, one user per line in the format `user: role`:  

```go
t := api.Team{Name: "crbot.public"}
crbotMembers := "```\n"

m := t.ListMembers()
for u, r := range m {
    line := fmt.Sprintf("%s: %s\n", u, r)
    crbotMembers += line
}
crbotMembers += "```"

u := api.Channel{Name: "dxb"}
u.SendMessage(crbotMembers)
```
