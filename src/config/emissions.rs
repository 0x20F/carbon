use serde::{ Serialize, Deserialize };
use std::collections::HashMap;
use std::fs;
use crate::error::Result;




static CONFIG_PATH: &'static str = "carbon-emissions.toml";

/// The "database" for all carbon related things
///
/// When something that's only relevant during
/// the current session needs to be kept track of
/// it gets stored here.
#[derive(Serialize, Deserialize, Debug)]
pub struct Emissions {
    /// All the containers that are currently running
    /// and what docker-compose file they're running in
    running: HashMap<String, Vec<String>>
}

impl Default for Emissions {
    fn default() -> Emissions {
        Self {
            running: HashMap::new()
        }
    }
}

impl Emissions {
    /// Try to load the emissions database from disk if 
    /// it exists. If it doesn't exist, create a new one
    /// with all the default values
    pub fn get() -> Self {
        // Create a fresh struct if it's not already written to file
        match fs::read_to_string(CONFIG_PATH) {
            Ok(s) => toml::from_str(&s).unwrap(),
            Err(_) => Self::default()
        }
    }


    /// Write the emissions database to disk
    pub fn save(config: &Self) -> Result<()> {
        let content = toml::to_string(config).unwrap();
        let path = format!("{}/{}", std::env::temp_dir().display(), CONFIG_PATH);
        fs::write(path, &content).expect("Couldn't save config file");

        Ok(())
    }


    /// Get the list of containers that are currently running
    /// and what docker-compose file they're running in
    pub fn get_running_services(&self) -> &HashMap<String, Vec<String>> {
        &self.running
    }


    /// Update the running services with a new list.
    /// Use this if you want to update the contents of the
    /// running services list
    pub fn set_running_services(&mut self, services: HashMap<String, Vec<String>>) {
        self.running = services;
    }


    /// Add a new service to the running services list
    pub fn add_running_service(&mut self, path: &str, services: &Vec<String>) {
        self.running.insert(
            path.to_string(),
            services
                .iter()
                .map(|s| String::from(s))
                .collect()
        );
    }
}
