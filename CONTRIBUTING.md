# Contributing

# Debugging arguments

| Flag               | Description                           |
|--------------------|---------------------------------------|
| `--debug`          | Tell the program to run in debug mode | 
| `--no-admin-check` | Bypass admin check                    |

## Debug mode

Tells the program that we're running this program as a developer, which does this:
* Uses the `config.yml` in the repository root directory, not the executables' directory.
  This is necessary when debugging from an IDE, as those binaries usually exist in a temporary directory.

## No admin check

Setting this flag bypasses checking if the program is run as root/admin. 
This is handy when debugging from an IDE.

# Adding providers

You can add additional providers by following these steps:
1. Add a provider file in the `providers` directory (like `linode.go`). 
2. Create a struct for that providers with the Provider struct embedded 
3. Implement all Provider functions. (Check `linode.go` for reference)
4. Add provider name and struct in `providers.availableProviders`
5. Add provider to `default.config.yml`
