<h1 align="center">Carbon</h1>
<p align="center">beast of a wrapper</p>

## What is this?
- You create a file called `carbon.yml` in your repository root
- That file contains the definition for how that repo will boot up as a docker container
- You can now build dynamic docker compose files using each `carbon.yml` in the repositories.

That's the gist of it. A script with nice output that basically concatenates small, chunked, docker compose service definitions
into a big file and runs them.

<br/>

## Why would you use this?
**I don't know, honestly.** Maybe you have multiple small docker services that are all defined in the same compose file but you'd like them to be separated into their own 
configuration files within each project. Maybe you're just here for the helpers, that's a valid reason too. 

I just built this because I needed it and it's now
public because I thought other people might want to use it too.

Apart from that, this provides:
- Unique container names for all started services.
- Easy wrapper commands for common things you might want to do, such as, getting a shell into a container is as simple as `co2 shell <container-name> | iex` (powershell).
- Easy executing of commands within one or multiple containers `co2 exec <container-A> <container-B> <container-C> -c echo "Hello World"`.
- Give each running docker container a unique ID which you can use to interact with it, preventing you from writing tedious 20+ character names by hand, or worse, copy pasting... (`co2 show -r` for a fresh table of all running containers).
- Neatly colored output.

Read more about the commands at the [Documentation Section](#documentation)


<br/>



<br/>

## Documentation
> **Note**: This program does have helpful wrapper commands for common docker things, however, it's not just that. Therefore it also contains more features that you may or may not choose to use. If you see the 📦 next to something, it's carbon specific and you probably shouldn't care if docker is all you want.

> **Double Note**: Keep in mind that all these options also exist in the program help menu which can be accessed by passing the `-h` param to any command or subcommand.

Let's start then. Here are all the command wrappers (and commands related to unique carbon functionality) so far and what they do:

<br/>

### `co2 show`
This one handles multiple things depending on the set flag:
- `-r` Will show **all** the running docker containers.
- 📦 `-a` Will show all the `carbon.yml` service files that are available for use.
- 📦 `-s` Shows all the _stores_ that carbon has access to

> Pro Tip: These can all be used together

<br/>

### `co2 shell`
This one will build a docker command that gets you a shell into whatever container or service you specified.
Do keep in mind, however, that since this **only returns** the composed command you still have to run it somehow.
Example:
```bash
$ co2 shell my-container
$ docker exec -it my-container /bin/bash
```
> Pro Tip: You can use the unique IDs that the [show](#co2-show) command displays to quickly specify a container

#### Valid Modifiers
- `-sh` To get `/bin/sh` instead of `bash`
- `-c` To execute a custom shell (or command if you really wanna) as in `co2 shell -c /bin/but-something-else my-container`