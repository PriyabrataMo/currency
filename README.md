💱 Currency Converter (TUI)

A simple terminal-based currency converter app.

## Installation

### 1. Tap & Install via Homebrew
```sh
brew tap PriyabrataMo/homebrew-taps
brew install currency
```

### 2. Acquire API Key
Obtain an API key from ExchangeRatesAPI.io.

### 3. Set the API Key
For macOS & Linux:
```sh
export API_KEY=actual_api_key
```
For Windows (Command Prompt):
```sh
set API_KEY=actual_api_key
```
For Windows (PowerShell):
```sh
$env:API_KEY="actual_api_key"
```

### 4. Run the Currency Converter
```sh
currency
```

## Uninstall
```sh
brew uninstall currency
```

## Troubleshooting
• "API_KEY is missing" → ensure the key is exported  
• "Command Not Found" → try:
```sh
brew link currency
```
• Permission issues → consider:
```sh
sudo brew install currency
```

## License
MIT License

## Contributing
Open issues or submit pull requests.

## Support
Give a ⭐ on GitHub.