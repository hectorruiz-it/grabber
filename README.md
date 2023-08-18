# Grabber
A Rust tool to manage multiple repositories in different platforms.

## What is it intended for?
This is intended for people who have to manage multiple repositories from different clients. This tool helps you to organize all the repositories you want to have cloned in your computer.

## What can I do with grabber?
With `grabber` you can do:
- Easily manage respositories using different keys and passwords.
- Get a list of each client repositories.
- Get a list of the ssh keys used.
- Get a list of ssh keys alias (platform).
- Clone repositories.
- Pull and Push to repositories **(WIP)**.
- Install all repositories through a toml file to easily migrate from devices **(WIP)**.
- Work with a DynamoDB table to have a shared storage for teams to collaborate **(WIP)**.

## Initialize
First you will start configuring the tool:
```shell 
grabber setup
```
### What this command does:
- Creates a new directory called `.grabber` at your **HOME**.
- Creates inside this new directory two files:
  - `grabber-config.toml`: SSH config file.
  - `grabber-repositories.toml`: Repositories database file.
- Will ask you to introduce some values to configure the SSH config file:
  - An alias to identify this keys.
  - The private key absolute path.

#### grabber-config
```toml
[platform-alias]
private_key = "/home/grabber/.ssh/azure"
public_key = "/home/grabber/.ssh/azure.pub"

[personal]
private_key = "/home/grabber/.ssh/github"
public_key = "/home/grabber/.ssh/github.pub"

```

#### grabber-repositories
```toml
[client.alias]
repositories = ["git@github.com:hectorruiz-it/grabber.git"]
[hectorruiz-it.personal]
repositories = ["git@github.com:hectorruiz-it/grabber.git"]
```

## New client
To add a new client and start cloning repositories just type:
```shell 
grabber new --client
```

## Add and clone one or more repositories
```shell 
grabber add -c <CLIENT>
```

## List
### List platforms
```shell
grabber list --platforms
╭─────────────────────────╮
│ Platforms SSH Key alias │
╞═════════════════════════╡
│ github                  │
├╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌┤
│ personal                │
╰─────────────────────────╯
```
### List platform configuration
```shell 
grabber list -p azure
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│ PERSONAL                                                                                                    │
╞═════════════════════════════════════════════════════════════════════════════════════════════════════════════╡
│ { private_key = "/Users/hruiz/.ssh/github-personal", public_key = "/Users/hruiz/.ssh/github-personal.pub" } │
╰─────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
```
### List client platforms
```shell
grabber list -c <CLIENT>
╭─────────────────────────╮
│ HECTORRUIZ-IT PLATFORMS │
╞═════════════════════════╡
│ personal                │
╰─────────────────────────╯
```
### List repositories of a client in a given platform
```shell
grabber list -c hectorruiz-it -p github
╭───────────────────────────────────────────────────────╮
│ HECTORRUIZ-IT PERSONAL REPOSITORIES                   │
╞═══════════════════════════════════════════════════════╡
│ "git@github.com:hectorruiz-it/grabber.git"            │
╰───────────────────────────────────────────────────────╯
```
