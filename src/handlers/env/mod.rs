use clap::ArgMatches;
use crate::error::Result;
use crate::config::Footprint;



/// Handler function for every parameter
/// and subcommand of the `env` command
pub fn handle(matches: &ArgMatches) -> Result<()> {
    if let None = matches.subcommand_matches("env") {
        return Ok(());
    }

    let matches = matches.subcommand_matches("env").unwrap();
    let mut config = Footprint::get();

    if let Some(matches) = matches.subcommand_matches("add") {
        let path: String = matches.value_of("path").unwrap().to_string();
        let id: String = matches.value_of("identifier").unwrap().to_string();

        config.add_env_file(&path, &id);
        Footprint::save(&config)?;
    }

    if let Some(matches) = matches.subcommand_matches("remove") {
        let id: String = matches.value_of("identifier").unwrap().to_string();
        
        config.remove_env_file(&id);
        Footprint::save(&config)?;
    }

    if let Some(_) = matches.subcommand_matches("list") {
        config.print_as_table();
    }

    if let Some(matches) = matches.subcommand_matches("activate") {
        let id: String = matches.value_of("identifier").unwrap().to_string();
        
        config.activate_env_file(&id);
        Footprint::save(&config)?;
    }

    Ok(())
}