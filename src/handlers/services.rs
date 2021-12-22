use crate::docker;
use crate::file;
use std::env;


pub struct ServiceStart {}

impl ServiceStart {
    pub fn with_services<'a>(services: Vec<&'a str>) {
        // Handle error a bit more nicely
        let environment = Self::get_root_directory();
        let mut configs = vec![];

        for service in services {
            // 1. Load config for each service and put it in a vector
            let path = format!("{}/{}/carbon.yml", environment, service);
            configs.push(file::get_contents(&path));
        }

        let compose = docker::build_compose_file(&configs);
        let cleaned = Self::parse_service_file(&compose);
        let temp_path = file::write_tmp("yml", &cleaned);
        // Log the path that was saved as well
        println!("Saved compose file to: {}", temp_path);

        // Check if argument for viewing the file was passed
        // and print the file if it has
        println!("{}", cleaned);

        docker::start_service_setup(&temp_path);
    }


    fn parse_service_file(contents: &str) -> String {
        contents.replace("${ROOT}", &Self::get_root_directory())
    }


    fn get_root_directory() -> String {
        env::var("PROJECTS_DIRECTORY").expect("Projects directory must be set in dotenv file")
    }
}