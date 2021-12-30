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
mod carbon;

pub use util::error;
pub use util::macros;

use carbon::Carbon;



fn main() {
    Carbon::init();
    Carbon::run();
}