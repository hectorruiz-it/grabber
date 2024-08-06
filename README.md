# Grabber
A GoLang tool to manage multiple repositories in different platforms.

## What is it intended for?
This is intended for people who have to manage multiple repositories from different repository platforms and authentications methods.

## What can I do with grabber?
With `grabber` you can do:
- Easily manage respositories using different keys and passwords by leverage in it to the OS Keyring.
- Get full list of current repositories.
- Get a list of all configured profiles (Identified Authentication methods).
- Clone, Pull and Push from/to repositories.
- Install all repositories through a JSON file to easily migrate from devices **(WIP)**.
- Work with a DynamoDB table to have a shared storage for teams to ease distribution **(WIP)**.

## Files
- `.grabber-config.json`: JSON Database. Every entry is a Profile with their respective repositories.
- `grabber-ssh-profiles`: Stores SSH Authentication profiles.
- `grabber-token-profiles`: Stores Token Authentication profiles.


## Commands
### Add profile
Creates a new Grabber Profile based on the authentication method you provide (ssh or token):

```shell
grabber add-profile profile_name --[ssh|token]
```

#### What this command does:
- Creates a new profile (Identified Authentication method).
- Based on wether you use ssh or token adds an entry to their respective file.

### Clone

With clone you perform a simple clone like you do with git. The unique difference is that you must specify a profile:

```shell
❯ grabber clone git@github.com:hectorruiz-it/grabber.git -p github
Enumerating objects: 986, done.
Counting objects: 100% (111/111), done.
Compressing objects: 100% (52/52), done.
Total 986 (delta 75), reused 59 (delta 59), pack-reused 875
```

#### What this command does:
- Clones the repository on your current path.
- Stores the repository url on `.grabber-config.json` on the profile entry you've specified.
- This is later used by push command.

### Pull

Pulls repository from the origin on your current directory.

```shell
grabber pull
```

### Push

Pushes your repository changes to your upstream origin:

```shell
grabber push
```

#### What this command does:
- Retrieves the profile that you used to clone this repository on `.grabber-config.json` file.
- Then based on your authentication method:
  - **ssh**: Retrieves private key path and retrieves your password from the keyring system.
  - **token**: Retrieves your token from the keyring system.
- Pushes to the remote origin.

### List

Gives you information of your current profiles and their respective repositories.

```shell
~/grabber# grabber list profiles
Profile                 AuthMethod
────────────────────────────────────
bitbucket-seidor        token
github                  ssh
gitlab                  ssh
github-iot              token
```

```shell
~/grabber# grabber list repositories -p github
Repositories
────────────────────────────────────
git@github.com:hectorruiz-it/grabber.git
git@github.com:hectorruiz-it/letme.git
git@github.com:lockedinspace/letme.git
```
