# Grabber
A GoLang tool to manage multiple repositories in different platforms.

## What is it intended for?
This is intended for people who have to manage multiple repositories from different repository platforms and authentications methods.

## What can I do with grabber?
With `grabber` you can do:
- Easily manage respositories using different keys and passwords by leverage in it to the OS Keyring.
- Get full list of current repositories.
- Get a list of all configured profiles (Identified Authentication methods).
- Clone repositories. **(WIP)**
- Pull and Push from/to repositories **(WIP)**.
- Install all repositories through a JSON file to easily migrate from devices **(WIP)**.
- Work with a DynamoDB table to have a shared storage for teams to ease distribution **(WIP)**.

## Commands
### New profile
```shell 
grabber new-profile
```
### What this command does:
- Creates a new profile (Identified Authentication method).
- Based on wether you use ssh, token or basic adds an entry to their respective file:
  - `grabber-basic-profiles`: Stores Basic Authentication profiles.
  - `grabber-ssh-profiles`: Stores SSH Authentication profiles
  - `grabber-token-profiles`: Stores Token Authentication profiles.

```shell
~/grabber# grabber add-profile github --ssh
> Enter SSH private key absolute path: /home/hruiz/.ssh/github
grabber: Your profile configuration is the following:
> Profile ID: github
> Private Key Path: /home/hruiz/.ssh/github
> Public Key Path: /home/hruiz/.ssh/github.pub
grabber: do you want to apply this configuration [Y/n]: Y
grabber: github profile configured.
```

### List Profiles

```shell 
~/grabber# grabber list-profiles
Profile                 AuthMethod
────────────────────────────────────
azure-devops            basic
azure-devops-certainly  basic
bitbucket-seidor        token
github                  ssh
github-iot              basic

```


