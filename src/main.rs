extern crate clap;
extern crate dotenv;
extern crate rand;

mod util;
mod handlers;
mod docker;
mod file;

use clap::{ Arg, App, SubCommand };
use dotenv::dotenv;



fn main() {
    dotenv().ok();

    let matches = App::new("carbon")
        .version("1.0")
        .author("0x20F")
        .about("Container build tool")
        .subcommand(SubCommand::with_name("start")
                        .about("Start a specific service")
                        .version("1.0")
                        .author("0x20F")
                        .arg(Arg::with_name("services")
                                .help("What service and/or services to start (all = *)")
                                .required(true)
                                .min_values(1))
                    )
        .get_matches();

    

    if let Some(matches) = matches.subcommand_matches("start") {
        // Have a subcommand handler here
        let services: Vec<_> = matches.values_of("services").unwrap().collect();

        handlers::Service::start(services);
    }
}
