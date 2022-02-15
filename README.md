<h1 align="center">Carbon</h1>
<p align="center">beast of a wrapper</p>

<br/>

#### Fast Travel:
- [Documentation](#documentation)
- [Getting Help](#getting-help)
- [Reporting Issues](#reporting-issues)
- [Contributing](#contributing)

<br/>

## What is this?
- You create a file called `carbon.yml` in your repository root
- That file contains the definition for how that repo will boot up as a docker container
- You can now build dynamic docker compose files using each `carbon.yml` in the repositories.
- Oooooor you just use this for the docker command wrappers, that works too I guess...

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

## Reporting Issues
Wanna complain? That's fine. Here's how:
- Make sure there aren't any other issues that already relate to yours
- Make sure that it's truly an issue! If you're not sure, start a discussion and we'll go from there.
- If none of the above, open an issue describing what you're trying to accomplish and the ways you've attempted to do that so far.
- When you're done, kindly bring that fork back so the issue gets solved. If you want...

<br/>

## Contributing
- If you've got a new feature in mind, open a discussion about it and we'll go from there
- If you want to pick up an existing issue because it speaks to your heart for some reason, leave a comment on that issue telling the owner that you're on it and go make a fork.

<br/>

## Getting Help
If you're not sure how to find something you're looking for, here are a few things you can try:
- Look at the help menus of each of the commands in your terminal using the `-h` flag.
- Look through the existing issues and see if anyone is already talking about the same thing.
- Look through the existing discussions and see if anyone is already talking about the same thing.
- Open a discussion and ask away, always start with a discussion, don't jump into an issue directly unless it's really obvious. If it truly is an issue, we'll elevate it from there.

<br/>

## Documentation
> **Note**: This program does have helpful wrapper commands for common docker things, however, it's not just that. Therefore it also contains more features that you may or may not choose to use.

> **Double Note**: Keep in mind that all these options also exist in the program help menu which can be accessed by passing the `-h` param to any command or subcommand.

> **Triple Note**:  If you see the 📦 next to something, it's carbon specific and you probably shouldn't care if docker is all you want.

Let's start then. Here are all the command wrappers (and commands related to unique carbon functionality) so far and what they do:

<br/>

### 📦 `carbon.yml`
The `carbon.yml` file is the heart and soul of all carbon specific functionality within the program.
This is just a simple declaration of the docker-compose kind, without any of the docker compose bits added.

> Note: Any valid docker compose field is valid in this file

Looks kind of like this:
```yaml
my-service:
    image: golang
    ports:
        - "80:80"
    volumes:
        - /some/path:/another/path
    command: tail -f /dev/null
```

That's pretty simple right?

Now to run that service:
- First make sure you've registered the parent directory of your repository as a store e.g if your repo is called `A`, and the parent `B` (`/B/A`), you register `B` not `A`. This allows for a single store registration
for multiple repositories that might live in the same directory. (Use the [store add command](#%F0%9F%93%A6-co2-store-add))
- Now that its registered, carbon should be able to find your service so starting it is trivial: `co2 service start my-service`

> Pro Tip: If you ever want more than one service defined in your file, you separate them using the yaml document separator `---`

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