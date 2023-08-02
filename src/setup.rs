use colored::*;
use std::fs;
use std::fs::OpenOptions;
use std::io;
use std::io::Write;
use std::path::Path;
use toml::map::Map;
use toml::Value;

struct Input {
    message: String,
}

impl Input {
    fn regular_input(&self) -> String {
        let mut value: String = String::new();
        print!("{}", self.message.bold());
        let _ = io::stdout().flush();
        io::stdin()
            .read_line(&mut value)
            .expect("Error reading from STDIN");
        value.pop();
        value
    }
    fn flow_control_input(&self) -> String {
        let mut value: String = String::new();
        print!("{}", self.message.bold().bold().truecolor(255, 171, 0));
        let _ = io::stdout().flush();
        io::stdin()
            .read_line(&mut value)
            .expect("Error reading from STDIN");
        value.pop();
        value
    }
}

pub fn setup(force: bool) {
    let home_directory_path = format!("{}/.grabber", dirs::home_dir().unwrap().display());
    let platform_config_file_path: String = format!("{}/grabber-config.toml", &home_directory_path);
    let repository_file_path: String =
        format!("{}/grabber-repositories.toml", &home_directory_path);

    let config_files: [String; 2] = [platform_config_file_path, repository_file_path];

    match check_current_config(force, &config_files) {
        Ok(_) => {
            create_grabber_configuration(&home_directory_path, config_files);
        }
        Err(err) => {
            eprintln!("{}", err);
            std::process::exit(2);
        }
    }

    let mut n: i32 = 0;
    while n == 0 {
        match create_config_file() {
            Ok(_) => println!(),
            Err(_) => println!("ERROR: Unable to create config file at ~/.grabber"),
        }
        let message: String = String::from("Do you want to add another SSH configuration [y/N]? ");
        let continue_creating: String = Input { message }.flow_control_input();

        if !continue_creating.eq("y") {
            n += 1;
        }
    }
}

fn check_current_config(force: bool, config_files: &[String; 2]) -> Result<(), String> {
    for file in config_files {
        let path = Path::new(file);
        if let Ok(metadata) = fs::metadata(path) {
            let exists = metadata.is_file();
            if exists && !force {
                return Err(format!(
                    "File '{}' already exists. Use --force flag to overwrite.",
                    file
                ));
            }
        }
    }

    Ok(())
}

fn create_grabber_configuration(home_directory_path: &str, config_files: [String; 2]) {
    match fs::create_dir(home_directory_path) {
        Ok(_) => println!("Directory has been created at: ~/.grabber/"),
        Err(_) => eprintln!("Directory already exists at: {}", &home_directory_path),
    }

    config_files
        .iter()
        .map(|file| match fs::File::create(file) {
            Ok(_) => println!("A new grabber file has been created at: {}", file),
            Err(err) => eprintln!("ERROR: {}", err),
        })
        .collect()
}

fn create_config_file() -> std::io::Result<()> {
    let message: String = String::from("Enter platform ssh key alias: ");
    let platform_ssh_key_alias: String = Input { message }.regular_input();

    let message: String = String::from("Enter private key absolute path: ");
    let private_key: String = Input { message }.regular_input();

    let message: String = String::from("Enter public key absolute path: ");
    let public_key: String = Input { message }.regular_input();

    let mut config: Map<String, Value> = Map::new();
    let mut values: Map<String, Value> = Map::new();

    values.insert(String::from("private_key"), Value::String(private_key));
    values.insert(String::from("public_key"), Value::String(public_key));
    config.insert(platform_ssh_key_alias, Value::Table(values));

    let config_file = toml::to_string(&config).expect("ERROR: Unable to parse data to TOML");

    let platform_config_file_path: String = format!(
        "{}/.grabber/grabber-config.toml",
        dirs::home_dir().unwrap().display()
    );

    let mut ssh_config_file = OpenOptions::new()
        .append(true)
        .open(platform_config_file_path)
        .expect("ERROR: Unable to open ssh config file. ¿Does it exist?");

    ssh_config_file
        .write_all(config_file.as_bytes())
        .expect("ERROR: Unable to create config file");
    Ok(())
}
