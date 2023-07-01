# vega-assistant

A simple program that helps you manage the vega node. 

- [The Vega official web page](https://vega.xyz/)
- [The Vega documentation](https://docs.vega.xyz/)
- [The Vega GitHub profile](https://github.com/vegaprotocol/)

## Installation

### Supported platforms for vega-assistant

I have tested this binary agains the following systems and architectures, but it may also works with different combinations. If this application is not working on your system, You can always setup the Vega node manually following the [Vega docs.](https://docs.vega.xyz)

- Linux (amd64 & arm64) including WSL2
- Macos (amd64 & arm64)

### Download binary from Github

Download one of pre-built binaries from the [GitHub release page](https://github.com/daniel1302/vega-assistant/releases/)

### Build it from sources

To build vega-assistant, you have to compile it with Go 1.20 or higher.

```shell
git clone https://github.com/daniel1302/vega-assistant.git
cd vega-assistant
go build -o ./vega-assistant ./main.go

./vega-assistant --help
```

## Available commands

### `vega-assistant setup postgresql`

This command is optional when you have PostgreSQL already configured. However, the Vega node requires PostgreSQL 14 with the TimescaleDB extension v2.8.0 installed. You can install all of the components on the same or a different server. This command prepare the `docker-compose.yaml` file that is ready to start a PostgreSQL server with [suggested server optimizations](https://docs.vega.xyz/testnet/node-operators/get-started/setup-datanode#postgresql-configuration-tuning)

#### Usage

Execute the following command

```shell
vega-assistant setup postgresql
```

Then fill the data and follow the instructions.
<br /><br />

### `vega-assistant setup data-node`

This command prepares the data node to be usage-ready on your computer. It asks about all custom details like home paths, SQL credentials, etc. Finally, initialize the node and gives you an instruction on how to run it.

#### Usage

Execute the following command

```shell
vega-assistant setup data-node
```

Then fill all the informations and follow the instruction on how to start the node. Optionally you can see the `vega-assistant setup systemd` command to prepare the systemd service.
<br /><br />

### `vega-assistant setup post-start`

You MUST call this command after your node has been started and you confirm it is moving blocks forward.

This command reverts configuration changes required to fast startup of the node. 

This command does not need to be executed when you start from block 0, but there is nothing wrong with calling it.

#### Usage

```shell
vega-assistant setup post-start
```
<br /><br />

### `vega-assistant setup systemd`

This command prepares the systemd service for your data node. 
There are some restrictions for this command:

- You must have initialized your node before you call this command.
- You must execute this command as a root user, otherwise, it will print the content of the systemd service to the stdout.
- This command is supported only for the Linux OS

#### Usage

```shell
vega-assistant setup systemd --visor-home <visor_home>
```

Flags:

- `--visor-home` - The home directory for vegavisor, you provided for the `vega-assistant setup data-node command`

