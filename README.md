# vega-asistant

A simple program that helps you manage the vega node.

## Available commands

### `vega-asistant setup postgresql`

This command is optional when you have PostgreSQL already configured. However, the Vega node requires PostgreSQL 14 with the TimescaleDB extension v2.8.0 installed. You can install all of the components on the same or a different server. This command prepare the `docker-compose.yaml` file that is ready to start a PostgreSQL server with [suggested server optimizations](https://docs.vega.xyz/testnet/node-operators/get-started/setup-datanode#postgresql-configuration-tuning)

#### Usage

Execute the following command

```shell
vega-asistant setup postgresql
```

Then fill the data and follow the instructions.
<br /><br />

### `vega-asistant setup data-node`

This command prepares the data node to be usage-ready on your computer. It asks about all custom details like home paths, SQL credentials, etc. Finally, initialize the node and gives you an instruction on how to run it.

#### Usage

Execute the following command

```shell
vega-asistant setup data-node
```

Then fill all the informations and follow the instruction on how to start the node. Optionally you can see the `vega-asistant setup systemd` command to prepare the systemd service.
<br /><br />

### `vega-asistant setup post-start`

You MUST call this command after your node has been started and you confirm it is moving blocks forward.

This command reverts configuration changes required to fast startup of the node. 

This command does not need to be executed when you start from block 0, but there is nothing wrong with calling it.

#### Usage

```shell
vega-asistant setup post-start
```
<br /><br />

### `vega-asistant setup systemd`

This command prepares the systemd service for your data node. 
There are some restrictions for this command:

- You must have initialized your node before you call this command.
- You must execute this command as a root user, otherwise, it will print the content of the systemd service to the stdout.
- This command is supported only for the Linux OS

#### Usage

```shell
vega-asistant setup systemd --visor-home <visor_home>
```

Flags:

- `--visor-home` - The home directory for vegavisor, you provided for the `vega-asistant setup data-node command`

