use std::process::{ Command };
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use crate::util;
use crate::file;
use yaml_rust::{ YamlLoader, Yaml, YamlEmitter };
use std::{ str, fs };


/// Initial compose file structure
static COMPOSE_FILE: &'static str = r#"
version: '3.8'

services:
"#;



/// Build a new docker-compose file by combining all the carbon.yml
/// files for each of the provided services into one file.
/// Making sure to indent everything properly.
pub fn build_compose_file(services: &Vec<&str>, carbon_conf: &str) -> Result<String> {
    let dir = util::environment::get_root_directory()?;
    let mut configs: Vec<String> = vec![];

    let mut output = String::new();
    let mut emitter = YamlEmitter::new(&mut output);

    // Find all the carbon.yml/carbon-isotope.yml files in the project directory
    for entry in fs::read_dir(dir).unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();

        if !path.is_dir() { continue; }

        let path = path.join(carbon_conf);
        match file::get_contents(&path.display().to_string()) {
            Ok(contents) => configs.push(contents),
            Err(_) => continue,
        };
    }

    // For each service that the user has requested,
    // look through all the configs that were found and see
    // if any of them contain that service.
    for service in services {
        let mut found = false;

        for config in configs.iter() {
            let docs = YamlLoader::load_from_str(config).unwrap();

            for doc in docs.iter() {
                match doc[*service] {
                    Yaml::BadValue => (),
                    _ => {
                        found = true;
                        emitter.dump(&doc).unwrap();
                        break;
                    }
                }
            }
        }

        if !found {
            return Err(CarbonError::ServiceNotDefined(service.to_string()));
        }
    }

    drop(emitter);

    // Cleanup the yaml output since the docker-compose file doesn't
    // need to contain multiple documents.
    let output = output.replace("---", "");
    Ok(merge_compose_file(&output))
}



/// Merge docker compose service declarations to the
/// main structure of a docker compose file.
/// Making sure to indent everything that needs indenting
fn merge_compose_file(services: &str) -> String {
    let mut definitions = vec![];
    
    // Indent the whole thing by one level (4 spaces)
    let lines: Vec<String> = services
        .split("\n")
        .map(|s| format!("    {}", s))
        .collect();

    definitions.push(lines.join("\n"));

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
