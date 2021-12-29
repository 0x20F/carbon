use std::process::{ Command };
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use crate::util;
use crate::file;
use crate::util::table::Table;
use std::{ str, fs };


/// Initial compose file structure
static COMPOSE_FILE: &'static str = r#"
version: '3.8'

services:
"#;



/// Build a new docker-compose file by combining all the carbon.yml
/// files for each of the provided services into one file.
/// Making sure to indent everything properly.
pub fn build_compose_file(services: &Vec<String>, carbon_conf: &str, clean: bool) -> Result<String> {
    let mut compose = vec![];
    let configs = find_carbon_services(carbon_conf)?;


    // For each service that the user has requested,
    // look through all the configs that were found and see
    // if any of them contain that service.
    for service in services {
        let mut found = false;

        for (_, config) in configs.iter() {
            // Split file into multiple documents 
            let docs = config.split("\n---\n").collect::<Vec<&str>>();

            // Check each document with serde
            for doc in docs.iter() {
                let mut v: serde_yaml::Value = serde_yaml::from_str(doc).unwrap();

                match v[service] {
                    serde_yaml::Value::Null => (),
                    _ => {
                        found = true;

                        // Generate a container name if the container doesn't already have one
                        if let serde_yaml::Value::Null = v[service]["container_name"] {
                            if !clean {
                                let name = format!("{}", service);
                                v[service]["container_name"] = name.into();
                            }
                        }

                        compose.push(serde_yaml::to_string(&v).unwrap());
                        break;
                    }
                }
            }
        }

        if !found {
            return Err(CarbonError::ServiceNotDefined(service.to_string()));
        }
    }

    // Cleanup the yaml output since the docker-compose file doesn't
    // need to contain multiple documents.
    let output = compose.join("\n").replace("---", "");
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



fn find_carbon_services(carbon_conf: &str) -> Result<Vec<(String, String)>> {
    let dir = util::environment::get_root_directory()?;
    let mut configs = vec![];

    // Find all the carbon.yml/carbon-isotope.yml files in the project directory
    for entry in fs::read_dir(dir).unwrap() {
        let entry = entry.unwrap();
        let path = entry.path();

        if !path.is_dir() { continue; }

        let path = path.join(carbon_conf);
        match file::get_contents(&path.display().to_string()) {
            Ok(contents) => configs.push((path.display().to_string(), contents)),
            Err(_) => continue,
        };
    }

    Ok(configs)
}



/// Given a docker compose file, start all the services 
/// within it and bubble up any errors that may arise.
pub fn start_service_setup(configuration: &str) -> Result<()> {
    let current_env = util::environment::current_env_path()?;

    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--env-file")
                    .arg(current_env)
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
    let current_env = util::environment::current_env_path()?;

    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--env-file")
                    .arg(current_env)
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




pub fn get_all_services(isotope: bool) -> Result<Vec<String>> {
    let file = if isotope { "carbon-isotope.yml" } else { "carbon.yml" };

    let configs = find_carbon_services(file)?;
    let mut services = vec![];

    for (_, config) in configs.iter() {
        let mut names: Vec<String> = config
            .split("---")
            .map(|s| s.trim())
            .map(|s| s.split(":").next().unwrap())
            .map(|s| s.trim())
            .map(|s| s.to_string())
            .collect();

        services.append(&mut names);
    }

    Ok(services)
}



/// Show all services that carbon has access to through the
/// currently active environment file in a nicely formatted table.
pub fn show_available() -> Result<()> {
    let mut configs = find_carbon_services("carbon.yml")?;
    configs.append(&mut find_carbon_services("carbon-isotope.yml")?);

    let mut names = vec![];

    for (path, config) in configs.iter() {
        let mut services: Vec<_> = config
            .split("---")
            .map(|s| s.trim())
            .map(|s| s.split(":").next().unwrap())
            .map(|s| (path.clone(), s))
            .collect();

        names.append(&mut services);
    }

    let mut table = Table::new(2, vec![]);

    table.header(vec!["Service", "Path"]);

    for (path, name) in names {
        let color = if path.contains("isotope.yml") {
            "cyan"
        } else {
            "black"
        };

        // Shorten the name
        let path = util::str::cut(&path, 30);

        table.row(vec![
            &format!("<{}>{}</>", color, name), 
            &format!("<{}>{}</>", color, path)
        ]);
    }

    table.display();
    Ok(())
}
