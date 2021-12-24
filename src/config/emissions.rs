use serde::{ Serialize, Deserialize };
use std::collections::HashMap;
use std::fs;
use crate::error::Result;




static CONFIG_PATH: &'static str = "/tmp/carbon-emissions.toml";


#[derive(Serialize, Deserialize, Debug)]
pub struct Emissions {
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
    pub fn get() -> Self {
        // Create a fresh struct if it's not already written to file
        match fs::read_to_string(CONFIG_PATH) {
            Ok(s) => toml::from_str(&s).unwrap(),
            Err(_) => Self::default()
        }
    }


    pub fn save(config: &Self) -> Result<()> {
        let content = toml::to_string(config).unwrap();
        fs::write(CONFIG_PATH, &content).expect("Couldn't save config file");

        Ok(())
    }


    pub fn get_running_services(&self) -> &HashMap<String, Vec<String>> {
        &self.running
    }


    pub fn set_running_services(&mut self, services: HashMap<String, Vec<String>>) {
        self.running = services;
    }


    pub fn add_running_service(&mut self, path: &str, services: Vec<&str>) {
        self.running.insert(
            path.to_string(),
            services
                .iter()
                .map(|s| String::from(*s))
                .collect()
        );
    }
}
