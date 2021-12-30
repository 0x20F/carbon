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
pub fn build_compose_file(
    services: &Vec<String>, 
    carbon_conf: &str, 
    clean: bool,
    dependencies: bool
) -> Result<String> {
    // Find configs for all services and their dependencies
    let full_services = find_service_dependencies(services, carbon_conf)?;
    let mut compose = vec![];

    if services.len() != full_services.len() && !dependencies {
        let message = "Dependencies weren't provided for the specified services, try using --auto-dependencies or providing them";
        return Err(CarbonError::ServiceDependenciesNotProvided(message.to_string()))
    }

    for service in full_services {
        let string = serde_yaml::to_string(&service).unwrap();
        let clean = string.replace("---", "");

        compose.push(clean);
    }

    // Cleanup the yaml output since the docker-compose file doesn't
    // need to contain multiple documents.
    let output = compose.join("\n");

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



fn find_service_dependencies(services: &Vec<String>, carbon_conf: &str) -> Result<Vec<serde_yaml::Value>> {
    let mut parsed = vec![];
    let mut dependencies = vec![];
    let mut yaml = find_carbon_services(carbon_conf, services)?;

    for (name, config) in yaml {
        if let serde_yaml::Value::Sequence(d) = &config[name]["depends_on"] {
            for dependency in d {
                let dependency = dependency.as_str().unwrap();
                dependencies.push(dependency.to_string());
            }
        }

        parsed.push(config);
    }

    if dependencies.len() > 0 {
        parsed.append(&mut find_service_dependencies(&dependencies, carbon_conf)?);
    }

    Ok(parsed)
}



fn find_carbon_services(carbon_conf: &str, services: &Vec<String>) -> Result<Vec<(String, serde_yaml::Value)>> {
    let dir = util::environment::get_root_directory()?;
    let mut configs = vec![];
    let mut parsed = vec![];

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

    for service in services {
        let mut found = false;

        for (path, config) in configs.iter() {
            // Split file into multiple documents 
            let docs = config.split("\n---\n").collect::<Vec<&str>>();
    
            // Check each document with serde
            for doc in docs.iter() {
                let v: serde_yaml::Value = serde_yaml::from_str(doc).unwrap();
    
                if services.contains(&"*".to_string()) {
                    parsed.push((path.clone(), v));
                    continue;
                }

                match v[service] {
                    serde_yaml::Value::Null => (),
                    _ => {
                        parsed.push((service.clone(), v));
                        found = true;
                        break;
                    }
                }
            }
        }

        if !found {
            return Err(CarbonError::ServiceNotDefined(service.to_string()));
        }
    }


    Ok(parsed)
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

    let configs = find_carbon_services(file, &vec!["*".to_string()])?;
    let mut services = vec![];

    for (_, config) in configs.iter() {
        let config = serde_yaml::to_string(config).unwrap();

        config
            .split("---")
            .map(|s| s.trim())
            .map(|s| s.split(":").next().unwrap())
            .map(|s| s.trim())
            .for_each(|s| services.push(s.to_string()));
    }

    Ok(services)
}



/// Show all services that carbon has access to through the
/// currently active environment file in a nicely formatted table.
pub fn show_available() -> Result<()> {
    let mut configs = find_carbon_services("carbon.yml", &vec!["*".to_string()])?;
    configs.append(&mut find_carbon_services("carbon-isotope.yml", &vec!["*".to_string()])?);

    let mut names = vec![];

    for (path, config) in configs.iter() {
        let string = serde_yaml::to_string(config).unwrap();

        string
            .split("---")
            .map(|s| s.trim())
            .map(|s| s.split(":").next().unwrap())
            .for_each(|s| names.push((path.clone(), s.to_string())));
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
