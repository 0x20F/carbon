mod services;


use services::Service;
use clap::ArgMatches;
use crate::error::Result;
use crate::docker;


/// Handler function for every parameter
/// and subcommand of the `service` command.
pub fn handle(matches: &ArgMatches) -> Result<()> {
    if let None = matches.subcommand_matches("service") {
        return Ok(());
    }

    let _ = docker::running()?;
    let service_matches = matches.subcommand_matches("service").unwrap();
    let mut service_handler = Service::new();

    if let Some(start_matches) = service_matches.subcommand_matches("start") {
        let services: Vec<_> = start_matches.values_of("services").unwrap().collect();
        let display = start_matches.is_present("display");
        let isotope = start_matches.is_present("isotope");

        service_handler.start(services, display, isotope)?;
    }
    
    if let Some(stop_matches) = service_matches.subcommand_matches("stop") {
        let services: Vec<_> = stop_matches.values_of("services").unwrap().collect();
        service_handler.stop(services)?;
    }

    if let Some(_) = service_matches.subcommand_matches("list") {
        docker::container::show_all();
    }

    if let Some(rebuild_matches) = service_matches.subcommand_matches("rebuild") {
        let services: Vec<_> = rebuild_matches.values_of("services").unwrap().collect();
        service_handler.rebuild(services)?;
    }

    if let Some(add_matches) = service_matches.subcommand_matches("add") {
        let network: String = add_matches.value_of("network").unwrap().to_string();
        let services: Vec<_> = add_matches.values_of("services").unwrap().collect();
        let isotope = add_matches.is_present("isotope");

        // Start services
        service_handler.start(services.clone(), false, isotope)?;

        // Add them to the network
        docker::network::connect(&network, &services)?;
    }

    Ok(())
}