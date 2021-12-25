use serde::{ Serialize, Deserialize };
use std::collections::HashMap;
use std::fs;
use crate::util::table::Table;
use crate::error::Result;
use std::path::PathBuf;



static CONFIG_PATH: &'static str = "~/.local/carbon-footprint.toml";


#[derive(Serialize, Deserialize, Debug)]
pub struct Dotenv {
    active: bool,
    id: String,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Footprint {
    dotenv: HashMap<String, Dotenv>
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


    pub fn get_current_env(&self) -> Option<String> {
        // Loop through all dotenv files and return the path of the one that's active
        for (path, info) in &self.dotenv {
            if info.active {
                return Some(path.to_string());
            }
        }

        None
    }


    pub fn add_env_file(&mut self, path: &str, id: &str) {
        self.disable_all_files();

        // Canonicalize path before saving
        let path = shellexpand::full(path).unwrap();
        let path = PathBuf::from(&path.to_string());
        
        let expanded = match fs::canonicalize(&path) {
            Ok(path) => path.display().to_string(),
            Err(_) => {
                println!("Couldn't find dotenv file at {}", path.display().to_string());
                return;
            }
        };


        // If ID is longer than 5 characters, shorten it
        let id = if id.len() > 5 {
            id[0..5].to_string()
        } else {
            id.to_string()
        };
        

        self.dotenv.insert(
            expanded.clone(), 
            Dotenv {
                active: true,
                id: id
            }
        );
    }


    pub fn activate_env_file(&mut self, id: &str) {
        self.disable_all_files();

        for (_, info) in self.dotenv.iter_mut() {
            if info.id == id {
                info.active = true;
            }
        }
    }


    pub fn remove_env_file(&mut self, id: &str) {
        let mut to_remove = vec![];

        for (key, info) in self.dotenv.iter() {
            if info.id == id {
                to_remove.push(key.to_string());
            }
        }

        for key in to_remove {
            self.dotenv.remove(&key);
        }
    }


    pub fn print_as_table(&self) {
        let mut table = Table::new(
            3, 
            vec!['^', '<', '^']
        );

        table.header(vec!["ID", "Path", "Active"]);

        // Only print a footer if no entries defined
        if self.dotenv.is_empty() {
            table.footer("No entries defined");
            return;
        }
        
        for (k, info) in self.dotenv.iter() {
            let enabled = if info.active { "<bright-green>1</>" } else { "0" };
            let id = if info.active { format!("<bright-green>{}</>", info.id) } else { info.id.to_string() };
            
            // If the path is longer than 40 characters, show a shortened version of it
            let path = if k.len() > 40 {
                format!("{}...{}", &k[0..10], &k[k.len() - 20..])
            } else {
                k.to_string()
            };
            
            let path = if info.active { format!("<bright-green>{}</>", path) } else { path };

            table.row(vec![&id, &path, enabled]);
        }

        table.display();
    }


    fn disable_all_files(&mut self) {
        for (_, info) in self.dotenv.iter_mut() {
            info.active = false;
        }
    }
}
