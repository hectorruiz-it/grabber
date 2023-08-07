use clap::{ArgGroup, Parser, Subcommand};
use std::process::exit;
mod add;
mod list;
mod new;
mod setup;

#[derive(Parser)]
#[clap(about, version, author)]
struct Value {
    #[clap(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    /// Add one or multiple repositories to a client
    Add {
        #[clap(short, long)]
        /// Client to add the repository
        client: String,

        #[clap(short, long)]
        /// Platform to add the repository
        platform: Option<String>,
    },
    /// List client platforms and repositories
    #[command(arg_required_else_help = true)]
    #[command(group = ArgGroup::new("simple").conflicts_with_all(["client", "platform"]))]
    List {
        #[clap(short, long)]
        /// Name of the client to list
        client: Option<String>,

        #[clap(short, long)]
        /// Name of the platform key alias
        platform: Option<String>,

        #[clap(group = "simple", long)]
        /// List all platform ssh key alias
        platforms: bool,

        #[clap(group = "simple", long)]
        /// List all clients
        clients: bool,
    },
    /// Adds a new Client
    New {
        #[clap(long)]
        /// Name of the client to add
        client: bool,

        #[clap(long)]
        /// Name of the platform to add
        platform: bool,
    },
    /// Configure script files and directory. You must run this first.
    Setup {
        #[clap(short, long)]
        /// Delete current configuration and start setup
        force: bool,
    },
}

fn main() {
    let value = Value::parse();
    match &value.command {
        Commands::Add { client, platform } => {
            if platform.is_none() {
                match add::add_platform_repository(client, platform) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
            if platform.is_some() {
                match add::add_platform_repository(client, platform) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
        }
        Commands::List {
            client,
            platform,
            platforms,
            clients,
        } => {
            if client.is_some() && platform.is_some() {
                match list::client_platform_repositories(
                    &client.to_owned().unwrap(),
                    &platform.to_owned().unwrap(),
                ) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                };
            };
            if client.is_some() && platform.is_none() {
                match list::client_platform(&client.to_owned().unwrap()) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
            if client.is_none() && platform.is_some() {
                match list::platform_key_alias_config(&platform.to_owned().unwrap()) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                };
            };
            if *platforms {
                match list::platforms() {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                };
            }
            if *clients {
                match list::clients() {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
        }
        Commands::New { client, platform } => {
            if *client && !*platform {
                new::client();
            } else if !*client && *platform {
                match new::platform() {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("{}", err);
                        exit(1)
                    }
                }
            } else {
                eprintln!("ERROR: flags --client and --platform can't be used together")
            }
        }
        Commands::Setup { force } => match setup::setup(*force) {
            Ok(_) => exit(0),
            Err(err) => {
                eprintln!("ERROR: {}", err);
                exit(1)
            }
        },
    }
}
