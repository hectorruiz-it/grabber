

// use std::{fs, path};

// use git2::build::RepoBuilder;
// use git2::{IndexAddOption, Repository, Signature};


// pub fn push() -> String {
//     let root_dir = path::Path::new("Z:/Temp");
//     let base_path = root_dir.join("base");
//     let remote_path = root_dir.join("remote");
//     let clone_path = root_dir.join("clone");
//     let author = Signature::now("user", "user@example.com").unwrap();

//     // create base repo and remote bare repo
//     let base_repo = Repository::init(&base_path).unwrap();
//     let remote_repo = Repository::init_bare(&remote_path).unwrap();
//     let remote_url = format!("file:///{}", remote_repo.path().display());

//     // create a text file and add it to index
//     fs::write(base_path.join("hello.txt"), "hello world!\n").unwrap();
//     let mut base_index = base_repo.index().unwrap();
//     base_index
//         .add_all(["."], IndexAddOption::DEFAULT, None)
//         .unwrap();
//     match base_index.write() {
//         Ok(_) => println!("some"),
//         Err(err) => err,
//     }

//     // make the commit, since it's the initial commit, there's no parent
//     let tree = base_repo
//         .find_tree(base_index.write_tree().expect_err())
//         .unwrap();
//     let commit_oid = base_repo
//         .commit(None, &author, &author, "initial", &tree, &[])
//         .unwrap();

//     // update branch pointer
//     let branch = base_repo
//         .branch("main", &base_repo.find_commit(commit_oid).unwrap(), true)
//         .unwrap();
//     let branch_ref = branch.into_reference();
//     let branch_ref_name = branch_ref.name().unwrap();
//     base_repo.set_head(branch_ref_name).unwrap();

//     // add remote as "origin" and push the branch
//     let mut origin = base_repo.remote("origin", &remote_url).unwrap();
//     origin.push(&[branch_ref_name], None).unwrap();

//     // clone from remote
//     let clone_repo = RepoBuilder::new()
//         .branch("main")
//         .clone(&remote_url, &clone_path)
//         .unwrap();

//     // examine the commit message:
//     println!(
//         "short commit message: {}",
//         clone_repo
//             .head()
//             .unwrap()
//             .peel_to_commit()
//             .unwrap()
//             .summary()
//             .unwrap()
//     );

//     let value = "a";
//     value.to_string();
// }