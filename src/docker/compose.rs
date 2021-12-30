use std::process::{ Command };
use crate::error::{ Result, CarbonError };
use crate::macros::unwrap_stderr;
use crate::util;
use crate::util::table::Table;
use crate::carbon::{ Carbon, ServiceData };
use std::str;


/// Initial compose file structure
static COMPOSE_FILE: &'static str = r#"
version: '3.8'

services:
"#;



/// Build a new docker-compose file by combining all the carbon.yml
/// files for each of the provided services into one file.
/// Making sure to indent everything properly.
/// 
/// ## Example
/// ```no_run
/// build_compose_file(
///     &vec!["service-a", "service-b"],
///     "carbon.yml",
///    true
/// )
/// ```
/// 
pub fn build_compose_file(
    services: &Vec<String>, 
    carbon_conf: &str, 
    dependencies: bool
) -> Result<String> {
    let full_services = Carbon::expand(services, carbon_conf)?;
    let mut compose = vec![];

    // If services have dependencies, and they're not provided
    // by the user, abort.
    if services.len() != full_services.len() && !dependencies {
        let message = "Dependencies weren't provided for the specified services, try using <cyan>--auto-dependencies</> or providing them manually.";
        return Err(CarbonError::ServiceDependenciesNotProvided(message.to_string()))
    }

    for service in full_services {
        let string = serde_yaml::to_string(&service).unwrap();
        let clean = string.replace("---", ""); // Serde outputs separate documents for each string, we don't want that

        compose.push(clean);
    }

    let output = compose.join("\n");
    Ok(merge_compose_file(&output))
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

    let configs = parse_all_services(file)?;
    let mut services = vec![];

    for s in configs.iter() {
        services.push(s.name.clone());
    }

    Ok(services)
}



/// Show all services that carbon has access to through the
/// currently active environment file in a nicely formatted table.
pub fn show_available() -> Result<()> {
    let mut services = parse_all_services("carbon.yml")?;
    services.append(&mut parse_all_services("carbon-isotope.yml")?);

    let mut names = vec![];

    for service in services.iter() {
        names.push((&service.path, &service.name));
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



/// Merge docker compose service declarations to the
/// main structure of a docker compose file.
/// Making sure to indent everything that needs indenting
/// 
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



/// Simple wrapper function to pass in '*' as service
/// to carbon's conversion to yaml
fn parse_all_services(carbon_conf: &str) -> Result<Vec<ServiceData>> {
    let services = vec!["*".to_string()];
    
    Carbon::as_yaml(carbon_conf, &services)
}
