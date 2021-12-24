use paris::formatter::colorize_string;






pub struct Table {
    column_count: usize,
    column_padding: Vec<usize>,
    alignment: Vec<char>
}

impl Table {
    pub fn new(
        column_count: usize, 
        mut column_padding: Vec<usize>,
        mut alignment: Vec<char>
    ) -> Table {
        if column_padding.len() < column_count {
            column_padding.resize(column_count, 20);
        }

        if alignment.len() < column_count {
            alignment.resize(column_count, '<');
        }


        Table {
            column_count,
            column_padding,
            alignment
        }
    }


    pub fn header(&self, elements: Vec<&str>) {
        let line = self.pad_strings(&elements);
        let mut spacer = vec![];

        for i in 0..self.column_count {
            spacer.push("-".repeat(self.column_padding[i] + (i + 1)));
        }

        let joined = spacer.join("+");

        log!("<black>{}</>\n{}\n<black>{}</>", joined, line, joined);
    }
    
    
    pub fn row(&self, elements: Vec<&str>) {
        let line = self.pad_strings(&elements);
        log!("{}", line);
    }
    
    
    pub fn footer(&self, text: &str) {
        // Format the text to be as wide as all the columns combined
        let padding = self.column_padding.iter().sum::<usize>();
        let padded = format!("{:width$}", text, width = padding);

        // Call log in case there are any colors in there
        log!("{}", padded);
    }


    fn pad_strings(&self, strings: &Vec<&str>) -> String {
        let mut padded = strings
            .iter()
            .enumerate()
            .map(|(i, e)| {
                // Parse any paris colors from the string
                let formatted = colorize_string(*e);
                let alignment = self.alignment[i];
                let padding = self.column_padding[i];
    
                // Figure out what length of characters is "invisible" characters
                let non_text = formatted
                    .chars()
                    .filter(|c| !c.is_alphanumeric() && !&['/', '.'].contains(c) )
                    .map(|c| c.len_utf8())
                    .sum::<usize>();

                // Increase the padding to make it even since it'll
                // add less because of the invisible characters
                let width = if non_text > 0 {
                    padding + non_text * 2 + 1
                } else {
                    padding
                };

    
                // Format the string with the custom width padding
                let aligned = match alignment {
                    '<' => format!("{:<width$}", formatted, width = width),
                    '>' => format!("{:>width$}", formatted, width = width),
                    'v' | '^' => format!("{:^width$}", formatted, width = width),
                    _ => panic!("Invalid alignment")
                };

                // Format one more time with separator
                format!("{} <black>|</>", aligned)
            })
            .collect::<Vec<String>>();
    
    
        // Put all the strings back together into a line
        format!("{}", padded.join(" "))
    }
}








#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_header() {
        header(vec!["<cyan>a</>", "<bright-green>a</>", "a"]);
    }
}