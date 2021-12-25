mod services;


use services::Service;
use clap::ArgMatches;
use crate::error::Result;
use crate::docker;



pub fn handle(matches: &ArgMatches) -> Result<()> {
    if let Some(service_matches) = matches.subcommand_matches("service") {
        let mut service_handler = Service::new();
    
        if let Some(start_matches) = service_matches.subcommand_matches("start") {
            let services: Vec<_> = start_matches.values_of("services").unwrap().collect();
            let display = start_matches.is_present("display");
    
            service_handler.start(services, display)?;
        }
        
        if let Some(stop_matches) = service_matches.subcommand_matches("stop") {
            let services: Vec<_> = stop_matches.values_of("services").unwrap().collect();
            service_handler.stop(services)?;
        }

        if let Some(list_matches) = service_matches.subcommand_matches("list") {
            docker::container::show_all();
        }
    }

    Ok(())
}