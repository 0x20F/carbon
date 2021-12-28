use serde::Deserialize;
use std::process::Command;
use std::collections::HashMap;
use crate::util::table::Table;



/// Code representation of the values we want
/// serde to retrieve from the massive JSON that
/// docker yields when inspecting a container
/// with `docker inspect container <name>`
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Container {
    /// The name of the container
    name: String,

    /// Network settings for the container
    #[serde(rename = "NetworkSettings")]
    pub settings: NetworkSettings,

    /// Current container state
    pub state: State,

    /// Current container configuration
    pub config: Config,

    /// Current container platform
    pub platform: String
}

impl Container {
    /// Get the name of the container.
    /// This is the only getter method since docker
    /// adds a / to the container name and it needs
    /// to be removed since they're not actually called by that name.
    pub fn name(&self) -> &str {
        &self.name[1..]
    }
}


/// Network settings for the container
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct NetworkSettings {
    pub networks: HashMap<String, Network>,
}


/// Individual network spec for a network that a
/// container is connected to
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Network {
    pub aliases: Option<Vec<String>>,
}


/// Current container state
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct State {
    status: String
}


/// Current container configuration
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Config {
    image: String,
}




/// Get all the containers that are available on
/// the current system and convert them from JSON
/// to a more usable format.
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



/// Get all the containers that are available on
/// the current system and display them in a nicely
/// colored and formatted table.
pub fn show_all() {
    let containers = all();
    let mut table = Table::new(5, vec![]);

    table.header(vec![ "Name", "Networks", "Status", "Image", "Platform" ]);

    for container in containers {
        let color = match container.state.status.as_str() {
            "running" => "cyan",
            "exited" => "black",
            "stopped" => "black",
            _ => "yellow"
        };
        let networks = container.settings.networks.keys().map(|s| &**s).collect::<Vec<_>>();
    
        let image = if container.config.image.len() > 10 {
            format!("{}...", &container.config.image[0..10])
        } else {
            container.config.image.to_string()
        };

        // If the container name is longer than 20 characters, show a shortened version of it
        let name = if container.name().len() > 20 {
            format!("{}...", &container.name()[0..20])
        } else {
            container.name().to_string()
        };


        table.row(vec![
            &format!("<{}>{}</>", color, name),
            &format!("<{}>{}</>", color, networks.join(", ")),
            &format!("<{}>{}</>", color, container.state.status),
            &format!("<{}>{}</>", color, image),
            &format!("<{}>{}</>", color, container.platform)
        ]);
    }

    table.display();
}



/// Inspect multiple containers and return a vector of
/// Deserialized structs that can be used throughout the
/// application.
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