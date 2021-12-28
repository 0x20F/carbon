use clap::{ App, Arg, SubCommand };



pub fn component() -> App<'static, 'static> {
    SubCommand::with_name("service")
        .alias("s")
        .about("Manage services")
        .subcommand(SubCommand::with_name("start")
            .alias("s")
            .about("Start a service")
            .arg(Arg::with_name("services")
                .help("Services to start")
                .required(true)
                .multiple(true)
                .index(1))
            .arg(Arg::with_name("display")
                .short("d")
                .long("display")
                .help("Display the compose file"))
            .arg(Arg::with_name("isotope")
                .short("i")
                .long("isotope")
                .help("Pick carbon-isotope.yml instead of carbon.yml"))
        )
        .subcommand(SubCommand::with_name("stop")
            .alias("p")
            .about("Stop a service")
            .arg(Arg::with_name("services")
                .help("Services to stop")
                .required(true)
                .multiple(true)
                .index(1))
        )
        .subcommand(SubCommand::with_name("list")
            .alias("ls")
            .about("List services")
            .arg(Arg::with_name("available")
                .short("a")
                .long("available")
                .help("List available services"))
        )
        .subcommand(SubCommand::with_name("rebuild")
            .alias("rb")
            .about("Rebuild a service")
            .arg(Arg::with_name("services")
                .help("Services to rebuild")
                .required(true)
                .multiple(true)
                .index(1))
        )
        .subcommand(SubCommand::with_name("add")
                .alias("a")
                .about("Start services and add them to a network")
                .arg(Arg::with_name("isotope")
                    .short("i")
                    .long("isotope")
                    .help("Pick carbon-isotope.yml instead of carbon.yml"))
                .arg(Arg::with_name("network")
                    .help("The network to add all the services to")
                    .required(true)
                    .index(1))
                .arg(Arg::with_name("services")
                    .help("Services to start")
                    .required(true)
                    .multiple(true)
                    .index(2))
        )
}