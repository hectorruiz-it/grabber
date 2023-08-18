use colored::*;
use dialoguer::{theme::ColorfulTheme, Confirm, Input, Password};
use git2::{Cred, RemoteCallbacks};
use keyring::Entry;
use regex::Regex;
use serde::Deserialize;
use std::fs::{self, File, OpenOptions};
use std::io::Read;
use std::io::Write;
use std::path::{Path, PathBuf};
use std::process::exit;
use std::result::Result;
use std::str::FromStr;
use toml::de::Error;
use toml::map::Map;
use toml::value::Table;
use toml::Value;

#[derive(Deserialize)]
struct Identifier {
    platform: String,
}
struct Config {
    platform_name: String,
    repositories: Vec<Value>,
}

impl Identifier {
    fn get_credential(&self) -> String {
        let entry = Entry::new(self.platform.as_str(), "grabber").expect("Failed to create entry");
        entry
            .get_password()
            .expect("No password has been configured for this platform")
    }
}

impl Config {
    fn add_platform(self) -> Map<String, Value> {
        let mut platform: Map<String, Value> = Map::new();
        let mut repositories: Map<String, Value> = Map::new();
        repositories.insert(
            String::from("repositories"),
            Value::Array(self.repositories),
        );
        platform.insert(self.platform_name, Value::Table(repositories));
        platform
    }
}

pub fn add_platform_repository(client: &String, platform: &Option<String>) -> Result<(), Error> {
    match platform {
        Some(platform_name) => {
            match add(client, platform_name) {
                Ok(_) => println!("{}", "New repositories have been configured".green().bold()),
                Err(_) => eprintln!("{}", "ERROR: Unable to add repositories".red().bold()),
            }
            Ok(())
        }
        None => {
            let platform_name: String = Input::with_theme(&ColorfulTheme::default())
                .with_prompt("Your name")
                .interact_text()
                .unwrap();

            match add(client, &platform_name) {
                Ok(_) => println!("{}", "New repositories have been configured".green().bold()),
                Err(_) => eprintln!("{}", "ERROR: Unable to add repositories".red().bold()),
            }
            Ok(())
        }
    }
}

fn add(client: &String, platform_name: &str) -> Result<(), Error> {
    let repositories_config_file_path = format!(
        "{}/.grabber/grabber-repositories.toml",
        dirs::home_dir().unwrap().display()
    );

    let mut file =
        File::open(&repositories_config_file_path).expect("ERROR: Please run grabber setup first");
    let mut contents = String::new();

    file.read_to_string(&mut contents).expect("msg");
    let mut dead_repositories: Vec<String> = Vec::new();

    match toml_edit::Document::from_str(&contents) {
        Ok(mut file) => match file[client][platform_name]["repositories"].as_array_mut() {
            None => {
                eprintln!("ERROR: client or platform doesn't exist. Run grabber list -c <CLIENT> to list platforms");
                exit(3)
            }
            Some(repositories) => {
                println!(
                    "{}",
                    "Use 'quit' to stop adding repositories"
                        .truecolor(255, 171, 0)
                        .bold()
                );
                let platform: String = platform_name.to_string();
                let ssh_config: Identifier = Identifier { platform };
                let password = ssh_config.get_credential();

                loop {
                    let passhprase: String = password.clone();
                    let ssh_clone_uri_pattern =
                        Regex::new(r"^(git@)?[a-zA-Z0-9.-]+(:|/).+\.git$").unwrap();

                    let repository_url: String = Input::with_theme(&ColorfulTheme::default())
                        .with_prompt("Enter repository SSH url")
                        .validate_with({
                            move |input: &String| -> Result<(), &str> {
                                if input == "quit" || ssh_clone_uri_pattern.is_match(input) {
                                    Ok(())
                                } else {
                                    Err("Invalid SSH clone URI format.")
                                }
                            }
                        })
                        .interact_text()
                        .unwrap();

                    if repository_url.eq("quit") {
                        break;
                    } else if repositories
                        .iter()
                        .any(|repo| repo.as_str() == Some(&repository_url))
                    {
                        eprintln!("{}", "✘ Repository already exists".red().bold());
                    } else {
                        match clone(platform_name, &repository_url, passhprase, client) {
                            Ok(_) => repositories.push(repository_url),
                            Err(_) => dead_repositories.push(repository_url),
                        }
                    }
                }
                fs::write(&repositories_config_file_path, file.to_string())
                    .expect("ERROR: Unable to write to ~/.grabber/grabber-repositories.toml");
            }
        },
        Err(_) => println!("ERROR: Unable to edit document"),
    }
    if dead_repositories.is_empty() {
        dead_repositories
            .iter()
            .map(|failed| println!("The following repositories failed to clone {}", failed))
            .collect()
    }
    Ok(())
}

pub fn new_platform() -> Result<(), String> {
    loop {
        let home: PathBuf = dirs::home_dir().expect("Home directory not found");
        let path: PathBuf = [".grabber", "grabber-config.toml"].iter().collect();
        let config_file: PathBuf = home.join(path);

        let platform_name: String = Input::with_theme(&ColorfulTheme::default())
            .with_prompt("Enter a platform alias for the ssh key")
            .interact_text()
            .unwrap();

        let key_path: String = Input::with_theme(&ColorfulTheme::default())
            .with_prompt("Enter the absolute path where the key is stored")
            .validate_with(|input: &String| -> Result<(), &str> {
                let key_path = Path::new(&input);
                if key_path.exists() {
                    Ok(())
                } else {
                    Err("ERROR: Private key does not exist. Remember to use an absolute path.")
                }
            })
            .interact_text()
            .unwrap();

        let password = Password::with_theme(&ColorfulTheme::default())
            .with_prompt("Enter passphrase")
            .validate_with(|input: &String| -> Result<(), &str> {
                if input.len() > 10 && !input.chars().all(char::is_alphabetic) && !input.chars().all(char::is_numeric) && !input.chars().all(char::is_lowercase){
                    Ok(())
                } else {
                    Err("Password must be longer than 10 and have at least 1 uppercase and 1 number")
                }
            })
            .with_confirmation("Enter same passphrase again", "Passphrases do not match. Try again")
            .interact()
            .unwrap();

        let entry = Entry::new(&platform_name, "grabber").expect("ERROR: entry already exists");
        entry
            .set_password(&password)
            .expect("ERROR: unable to open keyring platform");

        let public_key = format!("{}.pub", key_path);
        let mut config: Map<String, Value> = Map::new();
        let mut values: Map<String, Value> = Map::new();

        values.insert(String::from("private_key"), Value::String(key_path));
        values.insert(String::from("public_key"), Value::String(public_key));
        config.insert(platform_name, Value::Table(values));

        let toml_config_file =
            toml::to_string(&config).expect("ERROR: Unable to parse data to TOML");

        let mut platform_config_file: File = OpenOptions::new()
            .write(true)
            .append(true)
            .open(config_file)
            .expect("ERROR: Unable to open file with write permissions");

        platform_config_file
            .write_all(toml_config_file.as_bytes())
            .expect("ERROR: Unable to write data to config file");

        if !Confirm::with_theme(&ColorfulTheme::default())
            .with_prompt("Do you want to add another platform?")
            .interact()
            .expect("An unexpected error happened")
        {
            break;
        }
    }
    Ok(())
}

pub fn new_client() {
    let home: PathBuf = dirs::home_dir().expect("Home directory not found");
    let path: PathBuf = [".grabber", "grabber-repositories.toml"].iter().collect();
    let repositories_config_file: PathBuf = home.join(path);

    let mut file = OpenOptions::new()
        .append(true)
        .open(repositories_config_file)
        .expect("ERROR: Run 'grabber setup' to configure the script");

    let client_name: String = Input::with_theme(&ColorfulTheme::default())
        .with_prompt("Enter a client name")
        .interact_text()
        .unwrap();

    let mut dead_repositories: Vec<String> = Vec::new();

    loop {
        let mut repositories: Vec<Value> = Vec::new();
        let platform_name: String = Input::with_theme(&ColorfulTheme::default())
            .with_prompt("Enter platform ssh key alias")
            .interact_text()
            .unwrap();
        println!(
            "{}",
            "Use 'quit' to stop adding repositories"
                .truecolor(255, 171, 0)
                .bold()
        );
        let platform: String = platform_name.clone();
        let ssh_config: Identifier = Identifier { platform };
        let password = ssh_config.get_credential();

        loop {
            let passhprase: String = password.clone();
            let ssh_clone_uri_pattern: Regex =
                Regex::new(r"^(git@)?[a-zA-Z0-9.-]+(:|/).+\.git$").unwrap();
            let repository_url: String = Input::with_theme(&ColorfulTheme::default())
                .with_prompt("Enter repository SSH url")
                .validate_with({
                    move |input: &String| -> Result<(), &str> {
                        if input == "quit" || ssh_clone_uri_pattern.is_match(input) {
                            Ok(())
                        } else {
                            Err("Invalid SSH clone URI format.")
                        }
                    }
                })
                .interact_text()
                .unwrap();

            if repository_url.eq("quit") {
                break;
            } else if repositories
                .iter()
                .any(|repo| repo.as_str() == Some(&repository_url))
            {
                eprintln!("{}", "✘ Repository already exists".red().bold());
            } else {
                match clone(
                    platform_name.as_str(),
                    &repository_url,
                    passhprase,
                    &client_name,
                ) {
                    Ok(_) => repositories.push(Value::String(repository_url)),
                    Err(_) => dead_repositories.push(repository_url),
                }
            }
        }

        if dead_repositories.is_empty() {
            dead_repositories
                .iter()
                .map(|failed| eprintln!("The following repositories failed to clone {}", failed))
                .collect()
        }

        let config: Config = Config {
            platform_name,
            repositories,
        };
        let platform = config.add_platform();

        let mut client: Map<String, Value> = Map::new();
        client.insert(client_name.to_ascii_lowercase(), Value::Table(platform));

        let toml_content = toml::to_string(&client).expect("ERROR: Parse TOML error");

        file.write_all(toml_content.as_bytes())
            .expect("ERROR: Unable to write TOML file");

        if !Confirm::with_theme(&ColorfulTheme::default())
            .with_prompt("Do you want to add another platform?")
            .interact()
            .expect("msg")
        {
            break;
        }
    }

    println!(
        "{}: {}",
        "New respository platforms have been configured for"
            .green()
            .bold(),
        client_name
    );
}

fn clone(
    platform_name: &str,
    repository_url: &str,
    password: String,
    client: &str,
) -> Result<(), Error> {
    let home: PathBuf = dirs::home_dir().expect("Home directory not found");
    let config_path: PathBuf = [".grabber", "grabber-config.toml"].iter().collect();
    let config_file: PathBuf = home.join(config_path);

    let mut file = File::open(config_file).expect("asdf");
    let mut contents = String::new();

    file.read_to_string(&mut contents).expect("asdf");

    let toml: Table = toml::from_str(&contents).expect("ERROR: Unable to parse TOML file");

    if let Some(inner_table) = toml.get(platform_name).and_then(|v| v.as_table()) {
        if let (Some(private_key), Some(public_key)) = (
            inner_table.get("private_key").and_then(|v| v.as_str()),
            inner_table.get("public_key").and_then(|v| v.as_str()),
        ) {
            // Prepare callbacks.
            let mut callbacks = RemoteCallbacks::new();
            callbacks.credentials(|_url, _username_from_url, _allowed_types| {
                Cred::ssh_key(
                    "git",
                    Some(Path::new(public_key)),
                    Path::new(private_key),
                    Some(&password),
                )
            });

            // Prepare fetch options.
            let mut fo = git2::FetchOptions::new();
            fo.remote_callbacks(callbacks);

            // Prepare builder.
            let mut builder = git2::build::RepoBuilder::new();
            builder.fetch_options(fo);

            // Convert the clone URL to a Path
            let path = std::path::Path::new(&repository_url);

            // Get the file stem (repository name) without the extension
            if let Some(repo_name) = path.file_stem() {
                println!("Repository name: {:?}", repo_name);
            } else {
                println!("Invalid URL format");
            }
            let repo_name = Path::new(&repository_url)
                .file_stem()
                .and_then(|stem| stem.to_str())
                .unwrap_or("ERROR: Unknown repository");

            let home: PathBuf = dirs::home_dir().expect("Home directory not found");
            let repo_path: PathBuf = ["Workspace", client, repo_name].iter().collect();
            let clone_path: PathBuf = home.join(repo_path);
            // Clone the project.
            match builder.clone(repository_url, &clone_path) {
                Ok(_) => println!("Respository cloned at: {}", clone_path.display()),
                Err(err) => eprintln!("Error while cloning: {}", err),
            };
        }
    }
    Ok(())
}
