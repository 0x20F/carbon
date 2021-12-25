use std::process::Command;
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::collections::HashMap;
use serde::Deserialize;
use crate::util::table::Table;




#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Network {
    name: String,
    containers: Option<HashMap<String, Container>>,
}

#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Container {
    name: String,
}





pub fn create(name: &str) -> Result<()> {
    let output = Command::new("docker")
        .arg("network")
        .arg("create")
        .arg(name)
        .output()
        .expect("Something went wrong when creating a new network");

    unwrap_stderr!(output, DockerNetworkCreate)
}



pub fn remove(name: &str) -> Result<()> {
    let output = Command::new("docker")
        .arg("network")
        .arg("remove")
        .arg(name)
        .output()
        .expect("Something went wrong when deleting the network");

    unwrap_stderr!(output, DockerNetworkRemove)
}



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


pub fn connect(network: &str, container_names: &[&str]) -> Result<()> {
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


pub fn disconnect(network: &str, container_names: &[&str]) -> Result<()> {
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



fn run_if_container<F>(to_match: &[&str], f: F) -> Result<()>
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
