# AutoVPN

A tool for cheaper VPN connections.

This tool provisions single-session VPN servers at chosen VPS providers. Then it
automatically connects to the VPN server and destroys the VPN server when you
disconnect from the session.

## Table of contents
- [1. Setup](#setup)
    - [1.1 Install OpenVPN](#1-install-openvpn)
    - [1.2 Download binary](#2-download-binary)
    - [1.3 Unzip archive and run](#3-unzip-archive-and-run)
- [2. Configuration](#configuration)
- [3. Usage](#usage)

# Setup

## 1: Install OpenVPN

OpenVPN is required to connect to the VPN server.

**Note:** This script doesn't support "OpenVPN Connect", you need to install
the "OpenVPN" CLI tool.
If you're on Windows, you should find installers etc. here:
https://openvpn.net/community-downloads/

## 2: Download binary

Go to [Releases](https://github.com/Dekamik/autovpn/releases) and download the appropriate binary for your operating
system and architecture.

### Binary types:

|           | Windows | macOS (darwin) | Linux |
|-----------|---------|----------------|-------|
| **386**   | ✔️      | ❌              | ✔️    |
| **amd64** | ✔️      | ✔️             | ✔️    |
| **arm64** | ❌       | ✔️             | ❌     |

* 386 = x86, a.k.a. 32-bit computers
* amd64 = most 64-bit computers
* arm64 = Apple M-series (M1, M2 etc.)

## 3: Unzip archive and configure `config.yml`

Copy the `default.config.yml` and setup according to [Configuration](#configuration) below.

# Configuration

Configuration is done through the `config.yml` file. Here is an example:

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

OpenVPN must be run as administrator/root.

```
Usage: autovpn <provider> <region>  Provision a VPN server at <provider> on <region> and connects to it
       autovpn <provider> purge     Purges all AutoVPN servers at provider
       autovpn <provider>           Lists all regions at <provider>
       autovpn providers            Lists all available providers
       autovpn purge                Purges all AutoVPN servers at all providers
       autovpn (-h | --help)        Shows further help and options
       autovpn --version            Shows binary version
```
