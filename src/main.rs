extern crate clap;
extern crate dotenv;
extern crate rand;
extern crate serde;
#[macro_use] extern crate paris;


mod macros;
mod util;
mod handlers;
mod docker;
mod file;
mod error;
mod config;

use clap::{ Arg, App, SubCommand, ArgMatches };
use dotenv::dotenv;



fn main() {
    dotenv().ok();

    let matches = App::new("carbon")
        .version("1.0")
        .author("0x20F")
        .about("Container build tool")
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
                    )
        .get_matches();
    

    match execute(&matches) {
        Err(e) => error!("{}", e),
        _ => ()
    }
}


pub fn execute(matches: &ArgMatches) -> error::Result<()> {
    // Handle service actions
    if let Some(service_matches) = matches.subcommand_matches("service") {
        let mut service_handler = handlers::Service::new();
    
        if let Some(start_matches) = service_matches.subcommand_matches("start") {
            let services: Vec<_> = start_matches.values_of("services").unwrap().collect();
            let display = start_matches.is_present("display");
    
            service_handler.start(services, display)?;
        }
        
        if let Some(stop_matches) = service_matches.subcommand_matches("stop") {
            let services: Vec<_> = stop_matches.values_of("services").unwrap().collect();
            service_handler.stop(services)?;
        }
    }

    // Handle network actions
    if let Some(matches) = matches.subcommand_matches("network") {
        if let Some(matches) = matches.subcommand_matches("create") {
            let name: String = matches.value_of("name").unwrap().to_string();
            docker::network::create(&name)?;
        }

        if let Some(matches) = matches.subcommand_matches("remove") {
            let name: String = matches.value_of("name").unwrap().to_string();
            docker::network::remove(&name)?;
        }

        if let Some(_) = matches.subcommand_matches("list") {
            docker::network::show_all();
        }

        if let Some(matches) = matches.subcommand_matches("connect") {
            let network: String = matches.value_of("network").unwrap().to_string();
            let containers: Vec<_> = matches.values_of("container").unwrap().collect();
            docker::network::connect(&network, &containers)?;
        }
    }


    Ok(())
}