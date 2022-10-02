# Contributing

# Debugging options

| Option             | Description                           |
|--------------------|---------------------------------------|
| `--debug`          | Tell the program to run in debug mode | 
| `--no-admin-check` | Bypass admin check                    |

# Adding providers

You can add additional providers by following these steps:
1. Add a provider file in the `providers` directory (like `linode.go`). 
2. Create a struct for that providers with the Provider struct embedded 
3. Implement all Provider functions. (Check `linode.go` for reference)
4. Add provider name and struct in `providers.availableProviders`
5. Add provider to `default.config.yml`
