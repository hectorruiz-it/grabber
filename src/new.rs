use colored::*;
use dialoguer::Confirm;
use dialoguer::{theme::ColorfulTheme, Input};
use regex::Regex;
use std::fs::OpenOptions;
use std::io::Write;
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

pub fn new(client_name: &String) {
    let repositories_config_file = format!(
        "{}/.grabber/grabber-repositories.toml",
        dirs::home_dir().unwrap().display()
    );
    let mut file = OpenOptions::new()
        .append(true)
        .open(repositories_config_file)
        .expect("ERROR: Run 'grabber setup' to configure the script");

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
            // .default(false)
            .with_prompt("Do you want to continue?")
            .interact()
            .expect("msg")
        {
            break;
        }
    }

    println!(
        "{}: {}",
        "New respository platforms have been configured for".green().bold(),
        client_name
    );
}
