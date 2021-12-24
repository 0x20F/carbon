use std::process::Command;
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::collections::HashMap;
use serde::Deserialize;




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
    #[serde(rename = "IPv4Address")]
    ipv4: String,
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
    let inspections = inspect(&networks);
    let json: Vec<Network> = serde_json::from_str(&inspections).unwrap();


    print_table_header();
    for network in json {
        print_network_information(&network);
    }
}


pub fn connect(network: &str, containers: &[&str]) -> Result<()> {
    info!("Connecting <cyan>[</>{}<cyan>]</> containers to network: <magenta>{}</>", containers.join(", "), network);

    for container in containers.iter() {
        let output = Command::new("docker")
            .arg("network")
            .arg("connect")
            .arg(network)
            .arg(container)
            .output()
            .expect("Something went wrong when connecting to the network");

        // Print the error if it exists
        if !output.status.success() {
            let stderr = std::str::from_utf8(&output.stderr).unwrap();
            return Err(CarbonError::DockerNetworkConnect(stderr.to_string()));
        }
    }

    Ok(())
}




fn inspect(networks: &Vec<&str>) -> String {
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

    stdout
        .trim()
        .to_string()
}


fn print_table_header() {
    log!("<bright-green>#</> {:20}  <cyan>{}</>", "Networks", "Containers");
    log!("  {:20}  {}", "--------", "---------");
}


fn print_network_information(network: &Network) {
    let containers = match &network.containers {
        Some(c) => c.values().map(|c| String::from(&c.name)).collect::<Vec<String>>(),
        None => vec![]
    };

    log!("<bright-green>#</> {:20}  <cyan>[</> {} <cyan>]</>", network.name, containers.join(", "));
}