use paris::formatter::colorize_string;
use key_list::KeyList;






pub struct Table {
    column_count: usize,
    column_padding: Vec<usize>,
    alignment: Vec<char>,
    rows: Vec<String>
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
            alignment,
            rows: vec![]
        }
    }


    pub fn header(&mut self, elements: Vec<&str>) {
        let line = self.pad_strings(&elements);
        
        self.spacer();
        self.rows.push(line);
        self.spacer();
    }
    
    
    pub fn row(&mut self, elements: Vec<&str>) {
        let line = self.pad_strings(&elements);
        self.rows.push(line);
    }
    
    
    pub fn footer(&mut self, text: &str) {
        // Format the text to be as wide as all the columns combined
        let padding = self.column_padding.iter().sum::<usize>();
        let padded = format!("{:width$}", text, width = padding);

        // Call log in case there are any colors in there
        self.rows.push(padded);
    }


    pub fn display(&mut self) {
        // Close the table
        self.spacer();

        for row in &self.rows {
            log!("{}", row);
        }
    }


    fn spacer(&mut self) {
        let mut spacer = vec![];

        for i in 0..self.column_count {
            spacer.push("-".repeat(self.column_padding[i] + 2));
        }

        let joined = format!("<black>+{}+</>", spacer.join("+"));
        self.rows.push(joined);
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
                
                let keys = KeyList::new(&formatted, '[', 'm');
                let ansi_char_count = keys.map(|s| s.len()).sum::<usize>();
                let start_char = if i == 0 { "|" } else { "" };

                // Increase the padding to make it even since it'll
                // add less because of the invisible characters
                let width = if ansi_char_count > 0 {
                    padding + ansi_char_count + 2
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
                format!("<black>{}</> {} <black>|</>", start_char, aligned)
            })
            .collect::<Vec<String>>();
    
    
        // Put all the strings back together into a line
        format!("{}", padded.join(""))
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