mod env;
mod service;
mod network;


use clap::{ App, Arg, SubCommand, ArgMatches };



/// CLI Application "manifest" using clap.
/// This defines everything about carbon and its
/// CLI interface. (subcommands, arguments, etc.)
pub fn start() -> ArgMatches<'static> {
    App::new("carbon")
        .version("1.0")
        .author("0x20F")
        .about("Container build tool")
        .subcommand(env::component())
        .subcommand(service::component())
        .subcommand(network::component())
        .get_matches()
}