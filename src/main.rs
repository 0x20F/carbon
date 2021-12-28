extern crate clap;
extern crate dotenv;
extern crate rand;
extern crate serde;
extern crate shellexpand;
extern crate key_list;
#[macro_use] extern crate paris;


mod util;
mod handlers;
mod docker;
mod file;
mod config;
mod app;

pub use util::error;
pub use util::macros;

use clap::ArgMatches;
use std::{ env, fs };




fn main() {
    let footprint = config::Footprint::get();

    // If there is an active config file, load it.
    // If not, try loading one from the current directory.
    match footprint.get_current_env() {
        Some(path) => { dotenv::from_path(path).ok(); },
        None => { dotenv::dotenv().ok(); },
    }
    

    // If something breaks, print the error and exit
    match execute(&app::start()) {
        Err(e) => error!("{}", e),
        _ => ()
    }
}


/// Wrapper around all the actions that the handlers
/// execute so that we can nicely catch all the errors
/// here and print them out.
pub fn execute(matches: &ArgMatches) -> error::Result<()> {
    handlers::services::handle(matches)?;
    handlers::network::handle(matches)?;
    handlers::env::handle(matches)?;    

    Ok(())
}