use std::env;


pub fn get_root_directory() -> String {
    env::var("PROJECTS_DIRECTORY").expect("Projects directory must be set in dotenv file")
}


pub fn parse_variables(contents: &str) -> String {
    contents.replace("${ROOT}", &get_root_directory())
}