use std::process::{ Command };
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::str;


/// Initial compose file structure
static COMPOSE_FILE: &'static str = r#"
version: '3.8'

services:
"#;



/// Build a new docker-compose file by combining all the carbon.yml
/// files for each of the provided services into one file.
/// Making sure to indent everything properly.
pub fn build_compose_file(services: &[String]) -> String {
    let mut definitions = vec![];

    for service in services {
        // Indent the whole thing by one level (4 spaces)
        let lines: Vec<String> = service
            .split("\n")
            .map(|s| format!("    {}", s))
            .collect();

        definitions.push(lines.join("\n"));
    }

    format!("{}{}", COMPOSE_FILE, definitions.join("\n\n"))
}



/// Given a docker compose file, start all the services 
/// within it and bubble up any errors that may arise.
pub fn start_service_setup(configuration: &str) -> Result<()> {
    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--file")
                    .arg(configuration)
                    .arg("up")
                    .arg("-d")
                    .output()
                    .expect("Something went wrong when building the generated compose file");

    unwrap_stderr!(output, DockerServiceStartup)
}


/// Given a docker compose file and a service that is already stopped
/// do a hard start and rebuild everything about that service while keeping
/// it within the same compose file it was in before.
pub fn rebuild_specific_service_setup(name: &str, configuration: &str) -> Result<()> {
    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--file")
                    .arg(configuration)
                    .arg("up")
                    .arg("-d")
                    .arg("--no-deps")
                    .arg("--build")
                    .arg(name)
                    .output()
                    .expect("Something went wrong when trying to start a service in a running compose file");

    unwrap_stderr!(output, DockerServiceStartup)
}


/// Given a docker compose file and a service that is already running
/// Stop the running service.
pub fn stop_service_container(name: &str, configuration: &str) -> Result<()> {
    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--file")
                    .arg(configuration)
                    .arg("down")
                    .arg(name)
                    .output()
                    .expect("Something went wrong when trying to stop a service in a running compose file");

    unwrap_stderr!(output, DockerServiceShutdown)
}