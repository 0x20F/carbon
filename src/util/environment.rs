use std::env;


pub fn get_root_directory() -> String {
    env::var("PROJECTS_DIRECTORY").expect("Projects directory must be set in dotenv file")
}