use clap::ArgMatches;
use crate::error::Result;
use crate::config::Footprint;



pub fn handle(matches: &ArgMatches) -> Result<()> {
    // Handle env command actions
    if let Some(matches) = matches.subcommand_matches("env") {
        let mut config = Footprint::get();


        if let Some(matches) = matches.subcommand_matches("add") {
            let path: String = matches.value_of("path").unwrap().to_string();
            let id: String = matches.value_of("identifier").unwrap().to_string();

            // Add the path to the config
            config.add_env_file(&path, &id);
            Footprint::save(&config)?;
        }

        if let Some(matches) = matches.subcommand_matches("remove") {
            let id: String = matches.value_of("identifier").unwrap().to_string();
            
            // Remove the path from the config
            config.remove_env_file(&id);
            Footprint::save(&config)?;
        }

        if let Some(_) = matches.subcommand_matches("list") {
            // Print all paths and their data as a table
            config.print_as_table();
        }

        if let Some(matches) = matches.subcommand_matches("activate") {
            let id: String = matches.value_of("identifier").unwrap().to_string();
            
            // Set the active path
            config.activate_env_file(&id);
            Footprint::save(&config)?;
        }
    }


    Ok(())
}