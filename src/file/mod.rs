use std::fs;
use crate::util::generators;
use crate::error::{ Result, CarbonError };



/// Get contents of a file at a given path
pub fn get_contents(path: &str) -> Result<String> {
    match fs::read_to_string(path) {
        Ok(s) => Ok(s),
        _ => Err(CarbonError::FileReadError(path.to_string()))
    }
}


/// Save a given string in a file at the given path
pub fn save(path: &str, contents: &str) -> Result<()> {
    match fs::write(path, contents) {
        Ok(_) => Ok(()),
        _ => Err(CarbonError::FileWriteError(path.to_string()))
    }
}


/// Create a new file with the given extension and contents
/// in the /tmp directory with a random name.
pub fn write_tmp(extension: &str, content: &str) -> Result<String> {
    let name = generators::random_string(15);
    let path = format!("/tmp/{}.{}", name, extension);

    save(&path, content)?;

    Ok(path)
}
