# AutoVPN

A tool for cheaper VPN connections.

This tool provisions single-session VPN servers at a VPS provider.
When the VPN server is created, it automatically connects to the VPN server and
destroys the VPN server when you disconnect from the session.

![Simplified UML Sequence Diagram](docs/connect_seq_simplified.png)

## Table of content
- [1. Setup](#setup)
    - [1.1 Install Prerequisites](#1-install-prerequisties)
    - [1.2 Create an account at your chosen provider (currently Linode only)](#2-create-an-account-at-your-chosen-provider-currently-linode-only)
    - [1.3 Build and install binary](#3-build-and-install-binary)
    - [1.4 Configure .autovpn.yml](#4-configure-autovpnyml)
    - [1.5 Run](#5-run)
- [2. Configuration](#configuration)
- [3. Usage](#usage)
- [4. Build](#build)

# Setup

## 1: Install prerequisites

OpenVPN is required to connect to the VPN server.

**Note:** This script doesn't support "OpenVPN Connect", you need to install
the "OpenVPN" CLI tool. If you're on Windows, you should find installers etc.
here: https://openvpn.net/community-downloads/

To build the binary, `Go` and `make` have to be installed. (Make might not be
available on Windows, and can be omitted).

## 2: Create an account at your chosen provider (currently Linode only)

Go to the provider's website, sign up and generate the required API key(s).

## 3: Build and install binary

Build the binary by running `make`, followed by `make install` to install it.

NOTE: make may not work on Windows. To build the binary, you need to run the
following commands: `go mod tidy && go build cmd/main.go`.

Then you have to manually install it on your Windows machine and put it on your
`PATH`.

## 4: Configure `.autovpn.yml`

Copy the `.example.autovpn.yml` and setup according to [Configuration](#configuration) below.

## 5: Run

To get started, refer to [Usage](#usage) below, or run `autovpn --help`.

# Configuration

Configuration is done through the `.autovpn.yml` file. The executable will
first and foremost look for `.autovpn.yml` in the current working directory.
If not found, it will look for it in `home` next.

If the configuration file is defined using `-c`, that file will override the
default configuration files.

Here is an example:

`config.yml`
```yaml
agent:
  # script_url is the URL to the OpenVPN installation script
  script_url: https://raw.githubusercontent.com/angristan/openvpn-install/master/openvpn-install.sh

# overrides is optional and can be completely excluded
overrides:
  # Executable override
  openvpn_exe: C:\Program Files\ovpn\openvpn.exe 

# Providers contain a list of available providers and their configurations
providers:
  linode:
    image: linode/ubuntu20.04  # Image slug
    key: <LINODE_API_KEY>      # Your API key at provider
    type_slug: g6-dedicated-2  # Server type/size/tier slug
```

# Usage

When connecting to a VPN server, OpenVPN must be run as administrator/root.

```
Usage:	autovpn <provider> <region> Provision a VPN server at the specified
                                    provider on the specified region and then
                                    connects to it

    	autovpn <provider> list     Lists all region slugs at the provider

    	autovpn list                Lists all available providers

    	autovpn <provider> zombies  Lists all AutoVPN servers that should be
                                    destroyed at the provider

    	autovpn <provider> purge    Destroys all AutoVPN servers at the 
                                    provider

    	autovpn (--help)            Shows further help and options

    	autovpn --version           Shows version
```
