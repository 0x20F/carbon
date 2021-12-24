pub mod network;
pub mod container;


use std::process::{ Command };
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use std::str;



static COMPOSE_FILE: &'static str = r#"
version: '3.8'

services:
"#;



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
