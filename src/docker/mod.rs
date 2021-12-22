use std::process::Command;
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



pub fn start_service_setup(configuration: &str) {
    let output = Command::new("docker")
                    .arg("compose")
                    .arg("--file")
                    .arg(configuration)
                    .arg("up")
                    .arg("-d")
                    .output()
                    .expect("Something went wrong when building the generated compose file");

    let stdout = str::from_utf8(&output.stdout).unwrap();
    let stderr = str::from_utf8(&output.stderr).unwrap();

    if stderr == "" {
        // Everything went OK (with the command execution, can't vouch for docker)
        println!("{}", stdout);
    } else {
        // Something broke
        println!("{}", stderr);
    }
}