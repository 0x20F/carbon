use std::process::Command;
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::collections::HashMap;
use serde::Deserialize;




#[derive(Deserialize, Debug)]
#[serde(rename_all = "PascalCase")]
pub struct Network {
    name: String,
    containers: HashMap<String, Container>,
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
    let networks: Vec<&str> = stdout.split("\n").collect();

    print_table_header();

    for network in networks {
        let i = inspect(&network);
        let n = serde_json::from_str::<Network>(&i);

        if n.is_err() {
            continue; // Ignore errors for now
        }

        let n = n.unwrap();
        print_network_information(&n);
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




fn inspect(name: &str) -> String {
    let output = Command::new("docker")
        .arg("network")
        .arg("inspect")
        .arg(name)
        .output()
        .expect("Something went wrong when inspecting the network");

    let stdout = std::str::from_utf8(&output.stdout).unwrap();

    // Trim whitespace, [ and ] from the beginning and end of the string
    stdout.trim().trim_matches('[').trim_matches(']').to_string()
}


fn print_table_header() {
    log!("<bright-green>#</> {:20}  <cyan>{}</>", "Networks", "Containers");
    log!("  {:20}  {}", "--------", "---------");
}


fn print_network_information(network: &Network) {
    let containers = network.containers.values().map(|c| c.name.as_str()).collect::<Vec<&str>>();

    log!("<bright-green>#</> {:20}  <cyan>[</> {} <cyan>]</>", network.name, containers.join(", "));
}