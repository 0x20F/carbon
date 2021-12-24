use serde::{ Serialize, Deserialize };
use std::collections::HashMap;
use std::fs;
use crate::error::Result;



static CONFIG_PATH: &'static str = "~/.local/carbon-footprint.toml";


#[derive(Serialize, Deserialize, Debug)]
pub struct Footprint {
    dotenv: HashMap<String, bool>
}

impl Default for Footprint {
    fn default() -> Footprint {
        Self {
            dotenv: HashMap::new()
        }
    }
}

impl Footprint {
    pub fn get() -> Self {
        let expanded = shellexpand::full(CONFIG_PATH).unwrap().to_string();

        // Create a fresh struct if it's not already written to file
        match fs::read_to_string(expanded) {
            Ok(s) => toml::from_str(&s).unwrap(),
            Err(_) => Self::default()
        }
    }


    pub fn save(config: &Self) -> Result<()> {
        let expanded = shellexpand::full(CONFIG_PATH).unwrap().to_string();
        let content = toml::to_string(config).unwrap();
        
        fs::write(expanded, &content).expect("Couldn't save config file");

        Ok(())
    }
}
