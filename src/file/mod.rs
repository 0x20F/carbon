use std::fs;
use std::io::prelude::*;
use crate::util::generators;



pub fn get_contents(path: &str) -> String {
    fs::read_to_string(path).unwrap()
}


pub fn write_tmp(extension: &str, content: &str) -> String {
    let name = generators::random_string(15);
    let path = format!("/tmp/{}.{}", name, extension);

    let mut f = fs::write(&path, content).expect("Was unable to write a temporay compose file");

    path
}