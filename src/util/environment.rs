use std::env;
use crate::error::{ Result, CarbonError };
use crate::config::Footprint;


/// Get the root directory for all the services.
/// Hopefully defined in the .env file that is active.
pub fn get_root_directory() -> Result<String> {
    let var = "PROJECTS_DIRECTORY";
    let config = Footprint::get();

    let active = match config.get_current_env() {
        Some(s) => s,
        None => return Err(CarbonError::NoActiveEnv)
    };

    match env::var(var) {
        Ok(v) => Ok(v),
        Err(_) => Err(CarbonError::UndefinedEnvVar(var.to_string(), active))
    }
}


/// Get the directory of the .env file that's currently in use
pub fn current_env_path() -> Result<String> {
    let config = Footprint::get();

    match config.get_current_env() {
        Some(s) => Ok(s),
        None => Err(CarbonError::NoActiveEnv)
    }
}


/// Replace environment variables found
/// within each service configuration file with custom values
/// provided by carbon.
pub fn parse_variables(contents: &str) -> Result<String> {
    let replaced = contents.replace("${ROOT}", &get_root_directory()?);

    Ok(replaced)
}