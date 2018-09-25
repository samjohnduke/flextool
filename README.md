# Flextool

A handy tool for managing my services on digital ocean and other hosts. There
is an installable agent and a commandline app that is capable of connecting
to agents and running tasks.

## C&C

The command and control component is used by agents as a registration and central
authority server for storing logs and others.

## Agent

The agent is used for performing a bunch of tasks on the server. The core tasks are

- Export metrics such as CPU and Memory usage
- Connect to any database for running maintainence
- Observing the health of a server

## CLI

Used for interacting with the Agents and the C&C server