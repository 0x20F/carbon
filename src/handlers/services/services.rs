use crate::docker;
use crate::file;
use crate::util::environment;
use crate::error::{ Result, CarbonError };
use crate::config::Emissions;
use paris::Logger;
use std::collections::HashMap;


/// The filename standards for all the files that
/// every service should use to describe themselves.
static SERVICE_FILE: &'static str = "carbon.yml";
static SERVICE_FILE_ISOTOPE: &'static str = "carbon-isotope.yml";
static COMPOSE_FILE_FORMAT: &'static str = "yml";



/// Handler struct to take care of setting up and
/// taking down all the services based on their defined
/// configuration files.
pub struct Service<'p> {
    /// Paris logger since we want loading spinners.
    logger: Logger<'p>
}

impl<'p> Service<'p> {
    /// Create a new service handler
    pub fn new() -> Self {
        Self {
            logger: Logger::new()
        }
    }


    /// Given a list of services, if none of the services are
    /// already running, attempt to load their configuration files
    /// based on the active .env file and start them according to that.
    pub fn start<'a>(&mut self, services: Vec<&'a str>, display: bool, isotope: bool, save_path: Option<&str>) -> Result<()> {
        let mut carbon_conf = Emissions::get();
        let service_file = if isotope {
            self.logger.info("Loading isotope services..."); 
            SERVICE_FILE_ISOTOPE 
        } else {
            SERVICE_FILE 
        };

        // Save the generated compose file if told to
        if save_path.is_some() {
            let compose = docker::compose::build_compose_file(&services, &service_file)?;

            info!("Saving compose file to <bright-green>{}</>", save_path.unwrap());

            file::save(
                save_path.unwrap(), 
                &compose, 
            )?;

            return Ok(());
        }

        // Check if any of the provided services are already running
        // if they are we don't continue.
        for service in services.iter() {
            let values = carbon_conf.get_running_services().values();

            for value in values {
                if value.contains(&service.to_string()) {
                    return Err(CarbonError::ServiceAlreadyRunning(service.to_string()));
                }
            }
        }

        self.logger.info("Gathering individual service configurations...");

        let compose = docker::compose::build_compose_file(&services, &service_file)?;

        self.logger.info("Building docker-compose file for all services to live in...");

        let cleaned = environment::parse_variables(&compose)?;
        let temp_path = file::write_tmp(COMPOSE_FILE_FORMAT, &cleaned)?;

        self.logger.info("Starting all services...");
        
        if display {
            println!("Saved compose file to: {}", temp_path);
            println!("{}", cleaned);
        }

        docker::compose::start_service_setup(&temp_path)?;
        self.logger.success("Services should be up!");

        carbon_conf.add_running_service(&temp_path, services);
        Emissions::save(&carbon_conf)?;
        Ok(())
    }


    /// Given a list of services, attempt to stop them by carefully
    /// matching them with their individual docker compose files.
    /// Not running in the global scope with `docker container stop <name>` 
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
            docker::compose::stop_service_container(&containers, &compose_file)?;
        }
        self.logger.success("Stopped all required services");


        // Update the running services within the config
        config.set_running_services(to_keep);
        Emissions::save(&config)?;

        Ok(())
    }



    /// Given a list of services, find out which docker-compose
    /// file they belong to and stop them. Then rebuild them from
    /// that exact compose file so they spawn inside the same network
    /// they were previously in, but with updated images and whatnot.
    pub fn rebuild<'a>(&mut self, services: Vec<&'a str>) -> Result<()> {
        let config = Emissions::get();

        for s in services {
            // Find the compose file for each service
            for (compose_file, running) in config.get_running_services().iter() {
                if !running.contains(&s.to_string()) {
                    continue;
                }

                self.logger.loading(format!("Stopping service <bright-green>{}</> in (<magenta>{}</>)", s, compose_file));
                docker::compose::stop_service_container(s, compose_file)?;

                self.logger.info(format!("Rebuilding service: <bright-green>{}</> in (<magenta>{}</>)", s, compose_file));
                docker::compose::rebuild_specific_service_setup(s, &compose_file)?;

                // Only need to run once since docker doesn't allow
                // multiple containers to have the same name
                break;
            }
        }

        Ok(())
    }



    /// Helper function to figure out if a key already
    /// exists within a hashmap and if it does, push the
    /// value onto the vector, otherwise create a new
    /// vector and push the value onto it.
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