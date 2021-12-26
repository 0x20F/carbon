pub mod network;
pub mod container;
pub mod compose;


use crate::error::{ Result, CarbonError };
use std::process::Command;
use crate::macros::unwrap_stderr;



/// Check if the docker daemon is running.
/// This is useful so we don't perform unnecessary
/// actions if the daemon is not running.
pub fn running() -> Result<()> {
    let output = Command::new("docker")
        .arg("ps")
        .output()
        .map_err(|e| CarbonError::Docker(e.to_string()))?;

    
    unwrap_stderr!(output, Docker)
}