use std::fs;
use std::io::prelude::*;



pub fn get_contents(path: &str) -> String {
    fs::read_to_string(path).unwrap()
}


pub fn write_tmp(name: &str, content: &str) -> String {
    let path = format!("/tmp/{}", name);

    let mut f = fs::write(&path, content).expect("Was unable to write a temporay compose file");

    path
}