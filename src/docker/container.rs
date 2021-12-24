use serde::Deserialize;
use std::process::Command;
use std::collections::HashMap;




#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Container {
    name: String,

    #[serde(rename = "NetworkSettings")]
    pub settings: NetworkSettings,
}

impl Container {
    pub fn name(&self) -> &str {
        // Remove the leading slash that docker adds
        &self.name[1..]
    }
}


#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct NetworkSettings {
    pub networks: HashMap<String, Network>,
}


#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Network {
    pub aliases: Vec<String>,
}




pub fn all() -> Vec<Container> {
    let output = Command::new("docker")
        .arg("container")
        .arg("ls")
        .arg("--format")
        .arg("{{.Names}}")
        .output()
        .expect("Something went wrong when listing all networks");

    let stdout = std::str::from_utf8(&output.stdout).unwrap();
    let containers: Vec<&str> = stdout
        .trim()
        .split("\n")
        .collect();
    
    inspect(&containers)
}



fn inspect(containers: &[&str]) -> Vec<Container> {
    let mut command = Command::new("docker");

    command
        .arg("container")
        .arg("inspect");

    for container in containers {
        command.arg(container);
    }

    let output = command
        .output()
        .expect("Something went wrong when inspecting the containers");
    let stdout = std::str::from_utf8(&output.stdout).unwrap();

    if !output.status.success() {
        let stderr = std::str::from_utf8(&output.stderr).unwrap();
        println!("{}", stderr);
    }

    let json: Vec<Container> = serde_json::from_str(&stdout.trim().to_string()).unwrap();

    json
}