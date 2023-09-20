use clap::{ArgGroup, Parser, Subcommand};
use std::process::exit;
mod add;
mod list;
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
        /// profile to add the repository
        profile: Option<String>,
    },
    /// List client profiles and repositories
    #[command(arg_required_else_help = true)]
    #[command(group = ArgGroup::new("simple").conflicts_with_all(["client", "profile"]))]
    List {
        #[clap(short, long)]
        /// Name of the client to list
        client: Option<String>,

        #[clap(short, long)]
        /// Name of the ssh profile
        profile: Option<String>,

        #[clap(group = "simple", long)]
        /// List all ssh profiles
        profiles: bool,

        #[clap(group = "simple", long)]
        /// List all clients
        clients: bool,
    },
    /// Adds a new Client
    #[command(arg_required_else_help = true)]
    New {
        #[clap(long)]
        /// Name of the client to add
        client: bool,

        #[clap(long)]
        /// Name of the profile to add
        profile: bool,
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
        Commands::Add { client, profile } => {
            if profile.is_none() {
                match add::add_profile_repository(client, profile) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
            if profile.is_some() {
                match add::add_profile_repository(client, profile) {
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
            profile,
            profiles,
            clients,
        } => {
            if client.is_some() && profile.is_some() {
                match list::client_profile_repositories(
                    &client.to_owned().unwrap(),
                    &profile.to_owned().unwrap(),
                ) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                };
            };
            if client.is_some() && profile.is_none() {
                match list::client_profile(&client.to_owned().unwrap()) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                }
            }
            if client.is_none() && profile.is_some() {
                match list::profile_key_alias_config(&profile.to_owned().unwrap()) {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("ERROR: {}", err);
                        exit(1);
                    }
                };
            };
            if *profiles {
                match list::profiles() {
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
        Commands::New { client, profile } => {
            if *client && !*profile {
                add::new_client();
            } else if !*client && *profile {
                match add::new_profile() {
                    Ok(_) => exit(0),
                    Err(err) => {
                        eprintln!("{}", err);
                        exit(1)
                    }
                }
            } else {
                eprintln!("ERROR: flags --client and --profile can't be used together")
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
