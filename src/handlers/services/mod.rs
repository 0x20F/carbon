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

    // Handle start command
    start(&service_matches, &mut service_handler)?;

    // Handle add command
    add(&service_matches, &mut service_handler)?;
    
    // Handle stop command
    if let Some(stop_matches) = service_matches.subcommand_matches("stop") {
        let services: Vec<String> = stop_matches.values_of("services").unwrap().map(|s| s.to_string()).collect();
        service_handler.stop(&services)?;
    }

    // Handle list command
    if let Some(matches) = service_matches.subcommand_matches("list") {
        let available = matches.is_present("available");

        if !available {
            docker::container::show_all();
        } else {
            docker::compose::show_available()?;
        }
    }

    if let Some(rebuild_matches) = service_matches.subcommand_matches("rebuild") {
        let services: Vec<_> = rebuild_matches.values_of("services").unwrap().collect();
        service_handler.rebuild(services)?;
    }


    Ok(())
}



fn start(matches: &ArgMatches, service_handler: &mut Service) -> Result<()> {
    if let Some(start_matches) = matches.subcommand_matches("start") {
        let services: Vec<_> = start_matches.values_of("services").unwrap().collect();
        let display = start_matches.is_present("display");
        let isotope = start_matches.is_present("isotope");
        let save = start_matches.value_of("save");
        let dependencies = start_matches.is_present("dependencies");

        let services = all_or_provided(services, isotope)?;

        service_handler.start(
            &services, 
            display, 
            isotope, 
            save,
            dependencies
        )?;
    }

    Ok(())
}



fn add(matches: &ArgMatches, service_handler: &mut Service) -> Result<()> {
    if let Some(add_matches) = matches.subcommand_matches("add") {
        let network: String = add_matches.value_of("network").unwrap().to_string();
        let services: Vec<_> = add_matches.values_of("services").unwrap().collect();
        let isotope = add_matches.is_present("isotope");
        let dependencies = add_matches.is_present("dependencies");

        // Start services
        let services = all_or_provided(services, isotope)?;
        service_handler.start(&services, false, isotope, None, dependencies)?;

        // Add them to the network
        docker::network::connect(&network, &services)?;
    }

    Ok(())
}



/// Return all available services if the provided services
/// array contains a * wildcard. Otherwise return the
/// original services.
fn all_or_provided(services: Vec<&str>, isotope: bool) -> Result<Vec<String>> {
    if services.contains(&"all") {
        let names = docker::compose::get_all_services(isotope)?;

        return Ok(names);
    }

    Ok(services.iter().map(|s| s.to_string()).collect())
}