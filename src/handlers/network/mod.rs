use clap::ArgMatches;
use crate::error::Result;
use crate::docker;



/// Handler function for every parameter
/// and subcommand of the `network` command.
pub fn handle(matches: &ArgMatches) -> Result<()> {
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

        if let Some(matches) = matches.subcommand_matches("disconnect") {
            let network: String = matches.value_of("network").unwrap().to_string();
            let containers: Vec<_> = matches.values_of("container").unwrap().collect();
            docker::network::disconnect(&network, &containers)?;
        }
    }

    Ok(())
}