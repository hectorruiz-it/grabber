use colored::*;
use dialoguer::{theme::ColorfulTheme, Confirm, Input, Password};
use keyring::Entry;
use regex::Regex;
use std::fs::{File, OpenOptions};
use std::io::Write;
use std::path::{Path, PathBuf};
use std::result::Result;
use toml::map::Map;
use toml::Value;

struct Config {
    platform_name: String,
    repositories: Vec<Value>,
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

pub fn platform() -> Result<(), String> {
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

pub fn client() {
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
        loop {
            let ssh_clone_uri_pattern = Regex::new(r"^(git@)?[a-zA-Z0-9.-]+(:|/).+\.git$").unwrap();
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
                repositories.push(Value::String(repository_url));
            }
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
