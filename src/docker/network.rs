use std::process::Command;
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::collections::HashMap;
use serde::Deserialize;
use crate::util::table::Table;



/// Code representation of the values we want
/// serde to retrieve from the massive JSON that
/// docker yields when inspecting a network
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Network {
    /// The name of the network
    name: String,

    /// All the containers that are currently connected to the network
    /// if any...
    containers: Option<HashMap<String, Container>>,
}

/// Brief info about all the containers connected to the network
#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Container {
    name: String,
}




/// Given a name, create a new network with it.
pub fn create(name: &str) -> Result<()> {
    let output = Command::new("docker")
        .arg("network")
        .arg("create")
        .arg(name)
        .output()
        .expect("Something went wrong when creating a new network");

    unwrap_stderr!(output, DockerNetworkCreate)
}



/// Given a name, remove any network that has that name
pub fn remove(name: &str) -> Result<()> {
    let output = Command::new("docker")
        .arg("network")
        .arg("remove")
        .arg(name)
        .output()
        .expect("Something went wrong when deleting the network");

    unwrap_stderr!(output, DockerNetworkRemove)
}



/// Given a network name and a list of container names,
/// connect all the containers to the network.
pub fn connect(network: &str, container_names: &Vec<String>) -> Result<()> {
    run_if_container(
        container_names,
        |container| {
            Command::new("docker")
                .arg("network")
                .arg("connect")
                .arg(network)
                .arg(container.name())
                .output()
                .expect("Something went wrong when connecting to the network");

            // TODO: Report container being already connected?

            success!("Container <cyan>{}</> connected to the network <magenta>{}</>", container.name(), network);
        }
    )
}



/// Given a network name and a list of container names,
/// disconnect all the containers from the network.
pub fn disconnect(network: &str, container_names: &Vec<String>) -> Result<()> {
    run_if_container(
        container_names,
        |container| {
            Command::new("docker")
                .arg("network")
                .arg("disconnect")
                .arg(network)
                .arg(container.name())
                .output()
                .expect("Something went wrong when disconnecting from the network");

            // TODO: Report container not being part of the network?

            success!("Container <cyan>{}</> disconnected from the network <magenta>{}</>", container.name(), network);
        }
    )
}



/// Show all information we saved about all the networks
/// in a nicely colored and formatted table;
pub fn show_all() {
    let output = Command::new("docker")
        .arg("network")
        .arg("ls")
        .arg("--format")
        .arg("{{.Name}}")
        .output()
        .expect("Something went wrong when listing all networks");

    let stdout = std::str::from_utf8(&output.stdout).unwrap();
    let networks: Vec<&str> = stdout
        .trim()
        .split("\n")
        .collect();
    let json = inspect(&networks);

    let mut table = Table::new(2, vec![]);

    table.header(vec![
        "Name",
        "Containers"
    ]);

    for network in json {
        // Build a string of all the container names in a network
        let container_names = network
            .containers
            .as_ref()
            .unwrap()
            .iter()
            .map(|(_, container)| container.name.as_str())
            .collect::<Vec<&str>>()
            .join(", ");

       table.row(vec![
           network.name.as_str(),
           &format!("<cyan>[</> {} <cyan>]</>", container_names),
       ]); 
    }

    table.display();
}



/// Helper function to run a closure if the container
/// is found in the list of all available containers
fn run_if_container<F>(to_match: &Vec<String>, f: F) -> Result<()>
    where F: Fn(&super::container::Container) 
{
    let containers = super::container::all();

    for name in to_match {
        let mut exists = false;
        
        for container in containers.iter() {
            // Check if container actually exists
            if container.name() != *name {
                continue;
            }

            exists = true;

            f(container);
        }

        // No need to fail hard if the container doesn't exist
        // just inform the user that the action won't happen for hat container
        if !exists {
            warn!("(ignoring) Container <cyan>{}</> doesn't exist, <bright-green>try</> starting the service first", name);
        }
    }

    Ok(())
}



/// Given a list of network names, ask the docker API
/// for the JSON representation of each network and convert
/// it to a list of Network structs that we can use.
fn inspect(networks: &Vec<&str>) -> Vec<Network> {
    let mut command = Command::new("docker");
    
    command
        .arg("network")
        .arg("inspect");

    for network in networks.iter() {
        command.arg(*network);
    }

    let output = command.output().expect("Something went wrong when inspecting networks");
    let stdout = std::str::from_utf8(&output.stdout).unwrap();
    
    if !output.status.success() {
        let stderr = std::str::from_utf8(&output.stderr).unwrap();
        println!("{}", stderr);
    }

    let json: Vec<Network> = serde_json::from_str(&stdout.trim().to_string()).unwrap();

    json
}
