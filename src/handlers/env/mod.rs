use clap::ArgMatches;
use crate::error::Result;
use crate::config::Footprint;



pub fn handle(matches: &ArgMatches) -> Result<()> {
    // Handle env command actions
    if let Some(matches) = matches.subcommand_matches("env") {
        let config = Footprint::get();


        if let Some(matches) = matches.subcommand_matches("add") {
            let path: String = matches.value_of("name").unwrap().to_string();

            // Add the path to the config
        }

        if let Some(matches) = matches.subcommand_matches("remove") {
            let index: usize = matches.value_of("name").unwrap().parse().unwrap();
            
            // Remove the path from the config
        }

        if let Some(_) = matches.subcommand_matches("list") {
            // Print all paths and their data as a table
        }
    }


    Ok(())
}