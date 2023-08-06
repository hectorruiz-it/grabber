use colored::*;
use dialoguer::{theme::ColorfulTheme, Input};
use regex::Regex;
use std::fs::{self, File};
use std::io::Read;
use std::process::exit;
use std::str::FromStr;
use toml::de::Error;

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

fn add(client: &String, platform: &str) -> Result<(), Error> {
    let repositories_config_file_path = format!(
        "{}/.grabber/grabber-repositories.toml",
        dirs::home_dir().unwrap().display()
    );

    let mut file =
        File::open(&repositories_config_file_path).expect("ERROR: Please run grabber setup first");
    let mut contents = String::new();

    file.read_to_string(&mut contents).expect("msg");

    match toml_edit::Document::from_str(&contents) {
        Ok(mut file) => match file[client][platform]["repositories"].as_array_mut() {
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
                loop {
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
                        repositories.push(repository_url);
                    }
                }
                fs::write(&repositories_config_file_path, file.to_string())
                    .expect("ERROR: Unable to write to ~/.grabber/grabber-repositories.toml");
            }
        },
        Err(_) => println!("ERROR: Unable to edit document"),
    }
    Ok(())
}
