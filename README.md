# Discord Username Checker
This script allows you to check the availability of Discord usernames in bulk with multi-threading support.

## Configuration
1. Install Go 1.24+ on your system. (Or use the builded version found here)
2. Fill in the `data/config.json` file with the required information.
3. Add your proxies to the `data/proxies.txt` file (optional).
4. Run the script using `go run .` or build with `go build -o users.exe`. (Or use the pre-built binary found here)

## Misc
- Proxies are optional but recommended for large lists.
- Discord webhook notifications are supported.
- Custom username lists or random generation available.

## Config
```json
{
    "usernames": {
        "custom": false,        // Use custom username list or generate random
        "amount": 1000,         // Number of usernames to generate (if custom=false)
        "length": 3             // Length of generated usernames (if custom=false)
    },
    "retry": {
        "enabled": true,        // Enable retry on failed requests
        "max_attempts": 5       // Maximum retry attempts per username
    },
    "threads": 100,             // Number of concurrent threads
    "timeout": 30,              // Request timeout in seconds
    "webhook": ""               // Discord webhook URL for notifications on valid usernames
}
```

## Other Information
- Make sure to have the `data` directory on the same path as the executable. 

#
### Example
![Example](https://i.imgur.com/d7IlP8P.png)

[Discord](https://discord.gg/undesync)
