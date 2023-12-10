# Alpacon-CLI

`Alpacon CLI` is a powerful command-line tool designed for managing **AIP** (Alpaca Infra Platform) seamlessly from the terminal. This tool simplifies complex operations, making it easier for developers to interact with Alpacon services.

**NOTE** :  `Alpacon CLI` is currently in its test release phase. We're focusing on collecting user feedback and refining features

## Prerequisites
For the optimal use of `Alpacon CLI`, ensure that both [**Alpacon-Server**](https://github.com/alpacanetworks/alpacon-server) and [**Alpamon**](https://github.com/alpacanetworks/alpamon) are operational.
These components are integral for the CLI to function effectively.

## Documentation

**Note**: Detailed documentation, including usage guides and best practices, is in progress and will be available soon.

## Installation
Download `Alpacon CLI` as a binary directly from our [releases page](https://github.com/alpacanetworks/alpacon-cli/releases).
This ensures quick access to the most up-to-date version.

**Coming Soon**: Plans are in motion to provide Alpacon CLI through Homebrew for macOS users and as a native package for Linux distributions.

### macOS
```bash
wget https://github.com/alpacanetworks/alpacon-cli/releases/download/v0.0.1/alpacon-cli-0.0.1-darwin-arm64.tar.gz
tar -xvf alpacon-cli-0.0.1-darwin-arm64.tar.gz
chmod +x alpacon-cli
sudo mv alpacon-cli /usr/local/bin/alpacon
```

### Linux
```bash
wget https://github.com/alpacanetworks/alpacon-cli/releases/download/v0.0.1/alpacon-cli-0.0.1-linux-amd64.tar.gz
tar -xvf alpacon-cli-0.0.1-linux-amd64.tar.gz
chmod +x alpacon-cli
sudo mv alpacon-cli /usr/local/bin/alpacon
```

### Windows
Installation instructions for Windows will be provided soon.


### Login
To access and utilize all features of `Alpacon CLI`, first authenticate with the Alpacon API:

```bash
$ alpacon login

$ alpacon login -s=[SERVER URL] -u=[USERNAME] -p=[PASSWORD]
```
Successful login creates `config.json` in `~/.alpacon`, containing server address, API token, and token expiration (~1 week).
This file is essential for command execution and a new login is required upon token expiration.

## Usage
Explore Alpacon CLI's capabilities with the `-h` or `help` command.

```bash
$ alpacon -h

Use this tool to interact with the alpacon service.

Usage:
  alpacon [flags]
  alpacon [command]

Available Commands:
  agent       Commands to manage server's agent
  completion  Generate the autocompletion script for the specified shell
  download    Transfer a file from a remote server
  group       Manage Group (Identity and Access Management) resources
  help        Help about any command
  log         Retrieve and display server logs
  login       Log in to Alpacon Server
  package     Commands to manage and interact with packages
  server      Commands to manage and interact with servers
  upload      Transfer a file to a remote server
  user        Manage User resources
  websh       Open a websh terminal for a server
```


### Examples of Use Cases

#### Server Management
Manage and interact with servers efficiently using Alpacon CLI:
```bash
# List all servers.
$ alpacon server ls / list / all

# Get detailed information about a specific server.
$ alpacon server describe [SERVER NAME]

# Interactive server creation process.
$ alpacon server create

Server Name: 
Platform(debian, rhel): 
Groups:
[1] alpacon
[2] auditors
[3] designers
[4] developers
[5] managers
[6] operators
Select groups that are authorized to access this server. (e.g., 1,2):
```

#### Connect Websh
Access a server's websh terminal:
```bash
# Opens a websh terminal for the specified server.
$ alpacon websh [SERVER NAME]
```

#### Identity and Access Management (IAM)
Efficiently manage user and group resources:
```bash
# Managing Users

# List all users.
$ alpacon user ls / list / all

# Detailed user information.
$ alpacon user describe [USER NAME]

# Managing Groups

# List all groups.
$ alpacon group ls

# Detailed group information.
$ alpacon group describe [GROUP NAME]
```

#### File Transfer Protocol (FTP)
Facilitate file uploads and downloads:
```bash
# Upload files
$ alpacon upload alpacon.txt myserver:/home/alpacon/
$ alpacon cp /Users/alpacon.txt myserver:/home/alpacon/

# Download files
$ alpacon download myserver:/home/alpacon/alpacon.txt
$ alpacon cp myserver:/home/alpacon/alpacon.txt
```

#### Package Management
Handle Python and system packages effortlessly:
```bash
# python
$ alpacon package python ls / list / all
$ alpacon package python upload alpamon-1.1.0-py3-none-any.whl
$ alpacon package python cp /home/alpacon/alpamon-1.1.0-py3-none-any.whl

# system
$ alpacon package system ls / list /all
$ alpacon package system upload osquery-5.10.2-1.linux.x86_64.rpm
$ alpacon package system cp /home/alpacon/osquery_5.8.2-1.linux_amd64.deb
```

#### Logs Management
Retrieve and monitor server logs:
```bash
# View recent logs or tail specific logs.
$ alpacon logs [SERVER_NAME]
$ alpacon logs [SERVER NAME] --tail=10
```

#### Agent (Alpamon) Commands
Manage server agents(Alpamon) with ease:
```bash
# Commands to control and upgrade server agents.
$ alpacon agent restart [SERVER NAME]
$ alpacon agent upgrade [SERVER NAME]
$ alpacon agent shutdown [SERVER NAME]
```


### Contributing
We welcome bug reports and pull requests on our GitHub repository at https://github.com/alpacanetworks/alpacon-cli.
