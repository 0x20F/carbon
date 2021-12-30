use crate::error::{ Result, CarbonError };
use crate::{ app, util, file };
use clap::ArgMatches;
use std::fs;




pub struct ServiceData {
    pub name: String,
    pub path: String,
    pub yaml: serde_yaml::Value
}



pub struct Carbon {}



impl Carbon {
    pub fn init() {
        let footpring = crate::config::Footprint::get();

        // If there is an active config file, load it.
        // If not, try loading one from the current directory.
        match footpring.get_current_env() {
            Some(path) => { dotenv::from_path(path).ok(); },
            None => { dotenv::dotenv().ok(); },
        }
    }


    pub fn run() {
        match Carbon::execute(&app::start()) {
            Err(e) => error!("{}", e),
            _ => ()
        }
    }


    fn execute(matches: &ArgMatches) -> Result<()> {
        use crate::handlers;
        
        handlers::services::handle(matches)?;
        handlers::network::handle(matches)?;
        handlers::env::handle(matches)?;

        Ok(())
    }



    /// Returns yaml representations of all the given services,
    /// making sure to go down their dependency tree and add those
    /// services to the list as well.
    /// 
    /// Keep in mind, the amount of services you pass in,
    /// might not be the same amount as the services you get back.
    /// 
    pub fn expand(services: &Vec<String>, yml: &str) -> Result<Vec<serde_yaml::Value>> {
        let mut parsed = vec![];
        let mut dependencies = vec![];
        let yaml = Self::as_yaml(yml, services)?;

        for service in yaml {
            let config = service.yaml.clone();
            let name = service.name.clone();

            if let serde_yaml::Value::Sequence(d) = &config[name]["depends_on"] {
                for dependency in d {
                    let dependency = dependency.as_str().unwrap();
                    info!("Adding dependency: {}<cyan>-></>{}", service.name, dependency);
                    dependencies.push(dependency.to_string());
                }
            }

            parsed.push(config);
        }

        if dependencies.len() > 0 {
            parsed.append(&mut Self::expand(&dependencies, yml)?);
        }

        Ok(parsed)
    }



    /// Returns the name and yaml representation of all found
    /// services for a given configuration file (isotope or not).
    /// 
    pub fn as_yaml(yml: &str, services: &Vec<String>) -> Result<Vec<ServiceData>> {
        let mut parsed = vec![];
        let configs = Carbon::available(yml)?;

        for service in services {
            let mut found = false;

            for (path, config) in configs.iter() {
                // Split file into multiple documents 
                let docs = Self::service_as_yaml(config);

                for (name, yml) in docs {
                    if &name == service || services.contains(&"*".to_string()) {
                        found = true;
                        parsed.push(ServiceData {
                            name: name.to_string(),
                            path: path.to_string(),
                            yaml: yml
                        });
                    }

                    if &name == service {
                        break;
                    }
                }
            }

            if !found {
                return Err(CarbonError::ServiceNotDefined(service.to_string()));
            }
        }


        Ok(parsed)
    }



    /// Returns the name and path for all found services
    /// for a given configuration file (isotope or not).
    /// 
    pub fn available(yml: &str) -> Result<Vec<(String, String)>> {
        let dir = util::environment::get_root_directory()?;
        let mut configs = vec![];

        for entry in fs::read_dir(dir).unwrap() {
            let entry = entry.unwrap();
            let path = entry.path();

            if !path.is_dir() { continue; }

            let path = path.join(yml);
            match file::get_contents(&path.display().to_string()) {
                Ok(contents) => configs.push((path.display().to_string(), contents.to_string())),
                Err(_) => continue,
            };
        }

        Ok(configs)
    }



    fn service_as_yaml(service_config: &str) -> Vec<(String, serde_yaml::Value)> {
        // Split file into multiple documents 
        let docs = service_config.split("\n---\n").collect::<Vec<&str>>();
        let mut parsed = vec![];
            
        // Check each document with serde
        for doc in docs.iter() {
            let v: serde_yaml::Value = serde_yaml::from_str(doc).unwrap();

            if let serde_yaml::Value::Mapping(ref m) = v {
                let service_name = m.iter().next().unwrap().0;
                parsed.push((service_name.as_str().unwrap().to_string(), v));
            }
        }

        parsed
    }
}