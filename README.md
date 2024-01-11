# Alpacon-CLI

`Alpacon CLI` is a powerful command-line tool designed for managing **AIP** (Alpaca Infra Platform) seamlessly from the terminal. This tool simplifies complex operations, making it easier for developers to interact with Alpacon services.

**NOTE** :  `Alpacon CLI` is currently in its test release phase. We're focusing on collecting user feedback and refining features

## Prerequisites
For the optimal use of `Alpacon CLI`, ensure that both [**Alpacon-Server**](https://github.com/alpacanetworks/alpacon-server) and [**Alpamon**](https://github.com/alpacanetworks/alpamon) are operational.
These components are integral for the CLI to function effectively.

## Documentation

**Note**: Detailed documentation, including usage guides and best practices, is in progress and will be available soon.

## Installation
Download the latest `Alpacon CLI` directly from our releases page or install it using package managers on Linux.

### Docker
For every release and Release Candidate (RC), we push a corresponding container image to our Docker Hub repository at `alpacanetworks/alpacon-cli`. For example:

```bash
docker run --rm -it alpacanetworks/alpacon-cli version  
```

### Build the binary
- Make sure you have go installed:
```bash
git clone https://github.com/alpacanetworks/alpacon-cli.git
go build
sudo mv alpacon-cli /usr/local/bin/alpacon
```

### macOS
```bash
VERSION=<latest-version> # Replace with the actual version
wget https://github.com/alpacanetworks/alpacon-cli/releases/download/${VERSION}/alpacon-${VERSION}-darwin-arm64.tar.gz
tar -xvf alpacon-${VERSION}-darwin-arm64.tar.gz
chmod +x alpacon
sudo mv alpacon /usr/local/bin
```

### Linux

#### Debian and Ubuntu
```bash
curl -s https://packagecloud.io/install/repositories/alpacanetworks/alpacon/script.deb.sh?any=true | sudo bash

sudo apt-get install alpacon
```

#### CentOS and RHEL
```bash
curl -s https://packagecloud.io/install/repositories/alpacanetworks/alpacon/script.rpm.sh?any=true | sudo bash

sudo yum install alpacon
```

#### Download from GitHub Releases:
```bash
VERSION=<latest-version> # Replace with the actual version
wget https://github.com/alpacanetworks/alpacon-cli/releases/download/${VERSION}/alpacon-${VERSION}-linux-amd64.tar.gz
tar -xvf alpacon-${VERSION}-linux-amd64.tar.gz
chmod +x alpacon
sudo mv alpacon /usr/local/bin
```

### Windows
Installation instructions for Windows will be provided soon.


### Login
To access and utilize all features of `Alpacon CLI`, first authenticate with the Alpacon API:

```bash
$ alpacon login

$ alpacon login -s [SERVER URL] -u [USERNAME] -p [PASSWORD]
```
Successful login creates `config.json` in `~/.alpacon`, containing server address, API token, and token expiration (~1 week).
This file is essential for command execution and a new login is required upon token expiration.

#### Default Server URL
If you do not explicitly specify the server URL (-s) in the command, the default value `https://alpacon.io` is used. 
Therefore, you only need to use the `-s` option to specify a server URL if you wish to connect to a server other than the default one.

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
  event       Retrieve and display recent Alpacon events.
  group       Manage Group (Identity and Access Management) resources
  help        Help about any command
  log         Retrieve and display server logs
  login       Log in to Alpacon Server
  note        Manage and view server notes
  package     Commands to manage and interact with packages
  server      Commands to manage and interact with servers
  upload      Transfer a file to a remote server
  user        Manage User resources
  version     Displays the current CLI version.
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

# Delete server
$ alpacon server delete [SERVER NAME]

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

# Create a new user
$ alpacon user create

# Delete user
$ alpacon user delete [USER NAME]

# Managing Groups

# List all groups.
$ alpacon group ls

# Detailed group information.
$ alpacon group describe [GROUP NAME]

# Delete group
$ alpacon group delete [GROUP NAME]

# Add a member to a group with a specific role
$ alpacon group member add
$ alpacon group member add --group=[GROUP NAME] --member=[MEMBER NAME] --role=[ROLE]

# Remove a member from a group
$ alpacon group member delete --group=[GROUP NAME] --member=[MEMBER NAME]
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
- Download command: This command downloads the specified file to the **current directory** and also provides its **download URL**

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

#### Events Management
Retrieve and monitor events in the Alpacon:
```bash
# Display a list of recent events in the Alpacon
$ alpacon event
$ alpacon events

# Tail the last 10 events related to a specific server and requested by a specific user
$ alpacon event -tail 10 -s myserver -u admin
$ alpacon event --tail=10 --server=myserver --user=admin
```

#### Agent (Alpamon) Commands
Manage server agents(Alpamon) with ease:
```bash
# Commands to control and upgrade server agents.
$ alpacon agent restart [SERVER NAME]
$ alpacon agent upgrade [SERVER NAME]
$ alpacon agent shutdown [SERVER NAME]
```

#### Note Commands
Manage and view server notes:
```bash
# Display a list of all notes
$ alpacon note ls / list / all
$ alpacon note ls -s [SERVER NAME] --tail=10

# Create a note on the specified server
$ alpacon note create
$ alpacon note create -s [SERVER NAME] -c [CONTENT] -p [PRIVATE(true or false)]

# Delete a specified note
$ alpacon note delete [NOTE ID]
```


### Contributing
We welcome bug reports and pull requests on our GitHub repository at https://github.com/alpacanetworks/alpacon-cli.
