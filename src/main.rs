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
        .subcommand(SubCommand::with_name("start")
                        .about("Start one or multiple services")
                        .version("1.0")
                        .author("0x20F")
                        .arg(Arg::with_name("services")
                                .help("What service and/or services to start (all = *)")
                                .required(true)
                                .min_values(1))
                    )
        .subcommand(SubCommand::with_name("stop")
                        .about("Stop one or multiple services")
                        .version("1.0")
                        .author("0x20F")
                        .arg(Arg::with_name("services")
                                .help("What service and/or services to stop (all = *)")
                                .required(true)
                                .min_values(1))
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
                    )
        .get_matches();
    

    match execute(&matches) {
        Err(e) => error!("{}", e),
        _ => ()
    }
}


pub fn execute(matches: &ArgMatches) -> error::Result<()> {
    let mut service_handler = handlers::Service::new();

    // Handle service start
    if let Some(matches) = matches.subcommand_matches("start") {
        let services: Vec<_> = matches.values_of("services").unwrap().collect();
        service_handler.start(services)?;
    }

    // Handle service stop
    if let Some(matches) = matches.subcommand_matches("stop") {
        let services: Vec<_> = matches.values_of("services").unwrap().collect();
        service_handler.stop(services)?;
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

        if let Some(matches) = matches.subcommand_matches("list") {
            docker::network::show_all();
        }
    }


    Ok(())
}