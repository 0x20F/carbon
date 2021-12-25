use serde::Deserialize;
use std::process::Command;
use std::collections::HashMap;
use crate::util::table::Table;




#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Container {
    name: String,

    #[serde(rename = "NetworkSettings")]
    pub settings: NetworkSettings,

    pub state: State 
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
    pub aliases: Option<Vec<String>>,
}


#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct State {
    status: String
}




pub fn all() -> Vec<Container> {
    let output = Command::new("docker")
        .arg("container")
        .arg("ls")
        .arg("-a")
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



pub fn show_all() {
    let containers = all();
    let mut table = Table::new(3, vec![]);

    table.header(vec![ "Name", "Network", "Status" ]);

    for container in containers {
        let color = match container.state.status.as_str() {
            "running" => "cyan",
            "exited" => "black",
            "stopped" => "black",
            _ => "yellow"
        };
        let networks = container.settings.networks.keys().map(|s| &**s).collect::<Vec<_>>();
    
        table.row(vec![
            &format!("<{}>{}</>", color, container.name()),
            &format!("<{}>{}</>", color, networks.join(", ")),
            &format!("<{}>{}</>", color, container.state.status)
        ]);
    }

    table.display();
}



pub fn stop(container: &str) {
    Command::new("docker")
        .arg("container")
        .arg("stop")
        .arg(container)
        .output()
        .expect("Something went wrong when stopping a container");
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