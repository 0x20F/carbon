use crate::docker;
use crate::file;
use crate::util::environment;
use crate::error::Result;


pub struct Service {}

impl Service {
    pub fn start<'a>(services: Vec<&'a str>) -> Result<()> {
        // Handle error a bit more nicely
        let environment = environment::get_root_directory()?;
        let mut configs = vec![];

        for service in services {
            // 1. Load config for each service and put it in a vector
            let path = format!("{}/{}/carbon.yml", environment, service);
            configs.push(file::get_contents(&path)?);
        }

        let compose = docker::build_compose_file(&configs);
        let cleaned = environment::parse_variables(&compose)?;
        let temp_path = file::write_tmp("yml", &cleaned)?;

        // Log the path that was saved as well
        println!("Saved compose file to: {}", temp_path);

        // Check if argument for viewing the file was passed
        // and print the file if it has
        println!("{}", cleaned);

        docker::start_service_setup(&temp_path);

        Ok(())
    }


    pub fn stop<'a>(services: Vec<&'a str>) -> Result<()> {
        for service in services {
            println!("Stopping service: {}", service);
        }

        Ok(())
    }
}