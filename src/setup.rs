use std::fs;
use std::path::{Path, PathBuf};


pub fn setup(force: bool) -> Result<(), String>{
    let home: PathBuf = dirs::home_dir().expect("Home directory not found");
    let config_path: PathBuf = [".grabber", "grabber-repositories.toml"].iter().collect();
    let repositories_config_path: PathBuf = [".grabber", "grabber-repositories.toml"].iter().collect();
    let config_file: PathBuf = home.join(config_path);    
    let repositories_config_file: PathBuf = home.join(repositories_config_path);    

    let config_files: [PathBuf; 2] = [config_file, repositories_config_file];

    match check_current_config(force, &config_files) {
        Ok(_) => {
            create_grabber_configuration(&home, config_files);
            Ok(())
        }
        Err(err) => {
            Err(err)
        }
    }
}

fn check_current_config(force: bool, config_files: &[PathBuf; 2]) -> Result<(), String> {
    for file in config_files {
        let path = Path::new(file);
        if let Ok(metadata) = fs::metadata(path) {
            let exists = metadata.is_file();
            if exists && !force {
                return Err(format!(
                    "File '{}' already exists. Use --force flag to overwrite.",
                    file.display()
                ));
            }
        }
    }
    Ok(())
}

fn create_grabber_configuration(home: &PathBuf, config_files: [PathBuf; 2]) {
    match fs::create_dir(home) {
        Ok(_) => println!("Directory has been created at: ~/.grabber/"),
        Err(_) => eprintln!("Directory already exists at: {}", &home.display()),
    }

    config_files
        .iter()
        .map(|file| match fs::File::create(file) {
            Ok(_) => println!("A new grabber file has been created at: {}", file.display()),
            Err(err) => eprintln!("ERROR: {}", err),
        })
        .collect()
}
