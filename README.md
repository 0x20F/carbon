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
> **Note**: This program does have helpful wrapper commands for common docker things, however, it's not just that. Therefore it also contains more features that you may or may not choose to use.

> **Double Note**: Keep in mind that all these options also exist in the program help menu which can be accessed by passing the `-h` param to any command or subcommand.

> **Triple Note**:  If you see the 📦 next to something, it's carbon specific and you probably shouldn't care if docker is all you want.

Let's start then. Here are all the command wrappers (and commands related to unique carbon functionality) so far and what they do:

<br/>

### 📦 `carbon.yml`

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

<br/>

### 📦 `co2 store add`
On its own, this is useless, however with the 2 provided subcommands it does something.
This will _add_ a new directory(store) for carbon to look in when searching for `carbon.yml` files. It comes packed with 2 whole parameters:
- `-s` The path for the store, could be absolute (`/home/whatever/you`) or relative (`../../../sure`)
- `-i` A unique ID for the store you're adding. If not provided, one will be generated automatically so don't worry.
```bash
# Example Usage
$ co2 store add -s ../ -i unique-store
```
> Pro Tip: You can list all the stores with the [show command](#co2-show)

<br/>

### 📦 `co2 store remove`
The opposite of [add](#%F0%9F%93%A6-co2-store-add) as you'd expect. All it takes is that unique ID that you maybe defined, but definitely got with the `add` command.
```bash
$ co2 store remove unique-store
```
> Pro Tip: The remove command can take any number of store IDs

<br/>

### 📦 `co2 service start`
Looks through all the registered stores (see [add](#%F0%9F%93%A6-co2-store-add) on how to register stores) and starts all of the provided services
if they're found. 

As an added bonus, if the service defines any other services it is dependent on, usually within the `depends_on` field in the configuration, it will make sure that
those services are included in the provided list otherwise it'll abort and inform you that you're missing some important things.

Example:
```bash
$ co2 service start A B C
```
> Note: The names you provide here are what you defined within your carbon.yml file

If some of the provided services are already running but you'd like to stop them all and force a refresh, there's a flag for that:
- `-f` forces a service start, meaning all provided services will be stopped before attempting to start them again.

<br/>

### 📦 `co2 service stop`
Looks through the currently running **carbon** services and stops the provided ones.

Example:
```bash
$ co2 service stop A B C
```
> Note: The names you provide here are what you defined within your carbon.yml file