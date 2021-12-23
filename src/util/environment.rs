use std::env;
use crate::error::{ Result, CarbonError };


pub fn get_root_directory() -> Result<String> {
    let var = "PROJECTS_DIRECTORY";

    match env::var(var) {
        Ok(v) => Ok(v),
        Err(_) => Err(CarbonError::UndefinedEnvVar(var.to_string(), "TODO: .env path".to_string()))
    }
}


pub fn parse_variables(contents: &str) -> Result<String> {
    let replaced = contents.replace("${ROOT}", &get_root_directory()?);

    Ok(replaced)
}