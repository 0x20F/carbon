use paris::formatter::colorize_string;



pub fn header(elements: Vec<&str>) {
    let line = pad_strings(elements);
    let spacer = "-".repeat(line.len());

    println!("{}\n{}", line, spacer);
}


pub fn row(elements: Vec<&str>) {
    let line = pad_strings(elements);
    println!("{}", line);
}


fn pad_strings(strings: Vec<&str>) -> String {
    let mut padded = strings
        .iter()
        .map(|e| {
            // Parse any paris colors from the string
            let formatted = colorize_string(*e);

            // Figure out what length of characters is "invisible" characters
            let non_text = formatted
                .chars()
                .filter(|c| !c.is_alphanumeric())
                .collect::<String>();

            // Increase the padding to make it even since it'll
            // add less because of the invisible characters
            let width = 20 + non_text.len() * 2;

            // Format the string with the custom width padding
            format!("{:<width$}", formatted, width = width)
        })
        .collect::<Vec<String>>();


    // Put all the strings back together into a line
    format!("{}", padded.join(" "))
}






#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_header() {
        header(vec!["<cyan>a</>", "<bright-green>a</>", "a"]);
    }
}