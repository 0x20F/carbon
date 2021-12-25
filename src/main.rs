extern crate clap;
extern crate dotenv;
extern crate rand;
extern crate serde;
extern crate shellexpand;
extern crate key_list;
#[macro_use] extern crate paris;


mod util;
mod handlers;
mod docker;
mod file;
mod config;

pub use util::error;
pub use util::macros;

use clap::{ Arg, App, SubCommand, ArgMatches };



fn main() {
    let footprint = config::Footprint::get();

    match footprint.get_current_env() {
        // Take from the predefined path if provided
        Some(path) => { dotenv::from_path(path).ok(); },

        // Otherwise take from the current running directory
        None => { dotenv::dotenv().ok(); },
    }


    let matches = App::new("carbon")
        .version("1.0")
        .author("0x20F")
        .about("Container build tool")
        .subcommand(SubCommand::with_name("env")
                        .about("Manage dotenv files")
                        .subcommand(SubCommand::with_name("add")
                                        .about("Add a new dotenv file path")
                                        .arg(Arg::with_name("path")
                                                .help("Path to the dotenv file")
                                                .required(true)
                                                .index(1)
                                        )
                                        .arg(Arg::with_name("identifier")
                                                .help("Identifier for the dotenv file")
                                                .required(true)
                                                .index(2)
                                        )
                                    )
                        .subcommand(SubCommand::with_name("list")
                                        .about("List all dotenv files")
                                    )
                        .subcommand(SubCommand::with_name("remove")
                                        .about("Remove a dotenv file path")
                                        .arg(Arg::with_name("identifier")
                                                .help("The ID of the path you want to remove")
                                                .required(true)
                                                .index(1)
                                            )
                                    )
                        .subcommand(SubCommand::with_name("activate")
                                        .about("Set a dotenv file as active")
                                        .arg(Arg::with_name("identifier")
                                                .help("The ID of the path you want to activate")
                                                .required(true)
                                                .index(1)
                                            )
                                    )
                    )
        .subcommand(SubCommand::with_name("service")
                        .about("Manage services")
                        .subcommand(SubCommand::with_name("start")
                            .about("Start a service")
                            .arg(Arg::with_name("services")
                                .help("Services to start")
                                .required(true)
                                .multiple(true)
                                .index(1))
                            .arg(Arg::with_name("display")
                                .short("d")
                                .long("display")
                                .help("Display the compose file"))
                        )
                        .subcommand(SubCommand::with_name("stop")
                            .about("Stop a service")
                            .arg(Arg::with_name("services")
                                .help("Services to stop")
                                .required(true)
                                .multiple(true)
                                .index(1))
                        )
                        .subcommand(SubCommand::with_name("list")
                            .about("List all services")
                        )
                    )
        .subcommand(SubCommand::with_name("network")
                        .about("Perform actions on docker networks")
                        .version("1.0")
                        .author("0x20F")
                        .subcommand(SubCommand::with_name("create")
                                        .about("Create a new docker network")
                                        .version("1.0")
                                        .author("0x20F")
                                        .arg(Arg::with_name("name")
                                                .help("The name of the network")
                                                .required(true)
                                                .index(1))
                                    )
                        .subcommand(SubCommand::with_name("remove")
                                        .about("Remove a docker network")
                                        .version("1.0")
                                        .author("0x20F")
                                        .arg(Arg::with_name("name")
                                                .help("The name of the network")
                                                .required(true)
                                                .index(1))
                                    )
                        .subcommand(SubCommand::with_name("list")
                                        .about("List all docker networks")
                                        .version("1.0")
                                        .author("0x20F")
                                    )
                        .subcommand(SubCommand::with_name("connect")
                                        .about("Connect a container to a network")
                                        .version("1.0")
                                        .author("0x20F")
                                        .arg(Arg::with_name("network")
                                                .help("The name of the network")
                                                .required(true)
                                                .index(1))
                                        .arg(Arg::with_name("container")
                                                .help("The name/names of all containers that should connect to the network")
                                                .required(true)
                                                .index(2)
                                                .min_values(1))
                                    )
                        .subcommand(SubCommand::with_name("disconnect")
                                        .about("Disconnect a container from a network")
                                        .version("1.0")
                                        .author("0x20F")
                                        .arg(Arg::with_name("network")
                                                .help("The name of the network")
                                                .required(true)
                                                .index(1))
                                        .arg(Arg::with_name("container")
                                                .help("The name/names of all containers that should disconnect from the network")
                                                .required(true)
                                                .index(2)
                                                .min_values(1))
                                    )
                    )
        .get_matches();
    

    match execute(&matches) {
        Err(e) => error!("{}", e),
        _ => ()
    }
}


pub fn execute(matches: &ArgMatches) -> error::Result<()> {
    handlers::services::handle(matches)?;
    handlers::network::handle(matches)?;
    handlers::env::handle(matches)?;    

    Ok(())
}