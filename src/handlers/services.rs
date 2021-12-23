use crate::docker;
use crate::file;
use crate::util::environment;
use crate::error::Result;
use crate::config::Config;
use paris::Logger;
use std::collections::HashMap;


pub struct Service<'p> {
    logger: Logger<'p>
}

impl<'p> Service<'p> {
    pub fn new() -> Self {
        Self {
            logger: Logger::new()
        }
    }


    pub fn start<'a>(&mut self, services: Vec<&'a str>) -> Result<()> {
        // Handle error a bit more nicely
        let environment = environment::get_root_directory()?;
        let mut carbon_conf = Config::get();
        let mut configs = vec![];

        for service in services.iter() {
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

        carbon_conf.add_running_service(&temp_path, services);

        docker::start_service_setup(&temp_path)?;

        // Only save to config if service startup succeeded!
        Config::save(&carbon_conf)?;
        Ok(())
    }


    pub fn stop<'a>(&mut self, services: Vec<&'a str>) -> Result<()> {
        let mut config = Config::get();
        let mut to_stop: HashMap<String, Vec<String>> = HashMap::new();
        let mut to_keep: HashMap<String, Vec<String>> = HashMap::new();

        // Create a nice list of all the services that need to be stopped
        // and their corresponding docker compose files, making sure to 
        // remove those services from the config.
        for service in services {
            for (compose_file, running) in config.get_running_services().iter() {
                let own = service.to_string();

                if !running.contains(&own) {
                    Self::push_or_init(&mut to_keep, compose_file, own);
                    continue;
                }

                Self::push_or_init(&mut to_stop, compose_file, own);
            }
        }

        // Loop through all the compose files that have services
        // that need stopping and stop them
        for (compose_file, services) in to_stop.iter() {
            let containers = services.join(" ");

            self.logger.loading(format!("Stopping services [ <cyan>{}</> ] in compose file <bright-green>{}</>", containers, compose_file));
            docker::stop_service_container(&containers, &compose_file)?;
        }

        self.logger.success("Stopped all required services");

        // Update the running services within the config
        config.set_running_services(to_keep);
        Config::save(&config);
        
        Ok(())
    }


    fn push_or_init(map: &mut HashMap<String, Vec<String>>, key: &str, value: String) {
        if let Some(x) = map.get_mut(key) {
            x.push(value);
        } else {
            map.insert(
                key.to_string(),
                vec![ value ]
            );
        }
    }
}