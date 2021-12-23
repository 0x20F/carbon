extern crate clap;
extern crate dotenv;
extern crate rand;
#[macro_use] extern crate paris;

mod util;
mod handlers;
mod docker;
mod file;
mod error;

use clap::{ Arg, App, SubCommand, ArgMatches };
use dotenv::dotenv;
use paris::log;



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
        .get_matches();
    

    match execute(&matches) {
        Err(e) => error!("{}", e),
        _ => ()
    }
}


pub fn execute(matches: &ArgMatches) -> error::Result<()> {
    // Handle service start
    if let Some(matches) = matches.subcommand_matches("start") {
        let services: Vec<_> = matches.values_of("services").unwrap().collect();
        handlers::Service::start(services)?;
    }

    // Handle service stop
    if let Some(matches) = matches.subcommand_matches("stop") {
        let services: Vec<_> = matches.values_of("services").unwrap().collect();
        handlers::Service::stop(services)?;
    }


    Ok(())
}