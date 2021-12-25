use crate::docker;
use crate::file;
use crate::util::environment;
use crate::error::{ Result, CarbonError };
use crate::config::Emissions;
use paris::Logger;
use std::collections::HashMap;



static SERVICE_FILE: &'static str = "carbon.yml";
static COMPOSE_FILE_FORMAT: &'static str = "yml";



pub struct Service<'p> {
    logger: Logger<'p>
}

impl<'p> Service<'p> {
    pub fn new() -> Self {
        Self {
            logger: Logger::new()
        }
    }


    pub fn start<'a>(&mut self, services: Vec<&'a str>, display: bool) -> Result<()> {
        let environment = environment::get_root_directory()?;
        let mut carbon_conf = Emissions::get();
        let mut configs = vec![];

        // Check if any of the provided services are already running
        for service in services.iter() {
            // If they are, we don't continue
            let values = carbon_conf.get_running_services().values();

            for value in values {
                if value.contains(&service.to_string()) {
                    return Err(CarbonError::ServiceAlreadyRunning(service.to_string()));
                }
            }
        }

        self.logger.info("Gathering individual service configurations...");

        for service in services.iter() {
            let path = format!("{}/{}/{}", environment, service, SERVICE_FILE);
            configs.push(file::get_contents(&path)?);
        }

        self.logger.info("Building docker-compose file for all services to live in...");

        let compose = docker::build_compose_file(&configs);
        let cleaned = environment::parse_variables(&compose)?;
        let temp_path = file::write_tmp(COMPOSE_FILE_FORMAT, &cleaned)?;

        self.logger.info("Starting all services...");
        
        if display {
            println!("Saved compose file to: {}", temp_path);
            println!("{}", cleaned);
        }

        docker::start_service_setup(&temp_path)?;
        self.logger.success("Services should be up!");

        carbon_conf.add_running_service(&temp_path, services);
        Emissions::save(&carbon_conf)?;
        Ok(())
    }


    pub fn stop<'a>(&mut self, services: Vec<&'a str>) -> Result<()> {
        let mut config = Emissions::get();
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

            self.logger.loading(
                format!(
                    "Stopping services [ <cyan>{}</> ] in compose file <bright-green>{}</>", 
                    containers, 
                    compose_file
                )
            );
            docker::stop_service_container(&containers, &compose_file)?;
        }
        self.logger.success("Stopped all required services");


        // Update the running services within the config
        config.set_running_services(to_keep);
        Emissions::save(&config)?;

        Ok(())
    }


    pub fn rebuild<'a>(&mut self, services: Vec<&'a str>) -> Result<()> {
        let config = Emissions::get();

        for s in services {
            // Find the compose file for each service
            for (compose_file, running) in config.get_running_services().iter() {
                if !running.contains(&s.to_string()) {
                    continue;
                }

                self.logger.loading(format!("Stopping service <bright-green>{}</> in (<magenta>{}</>)", s, compose_file));
                docker::container::stop(s);

                self.logger.info(format!("Rebuilding service: <bright-green>{}</> in (<magenta>{}</>)", s, compose_file));
                docker::rebuild_specific_service_setup(s, &compose_file)?;

                // Only need to run once since docker doesn't allow
                // multiple containers to have the same name
                break;
            }
        }

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