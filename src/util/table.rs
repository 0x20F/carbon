use paris::formatter::colorize_string;
use key_list::KeyList;
use std::fmt::Display;






pub struct Table {
    column_count: usize,
    column_padding: Vec<usize>,
    alignment: Vec<char>,
    rows: Vec<String>,
    columns: Vec<Vec<String>>,
    longest: Vec<usize>
}

impl Table {
    pub fn new(
        column_count: usize, 
        mut column_padding: Vec<usize>,
        mut alignment: Vec<char>
    ) -> Self {
        if column_padding.len() < column_count {
            column_padding.resize(column_count, 20);
        }

        if alignment.len() < column_count {
            alignment.resize(column_count, '<');
        }

        let mut cols = vec![];
        for _ in 0..column_count {
            cols.push(vec![]);
        }

        let mut longest = vec![];
        for _ in 0..column_count {
            longest.push(0);
        }


        Table {
            column_count,
            column_padding,
            alignment,
            rows: vec![],
            columns: cols,
            longest
        }
    }


    pub fn header(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }
    
    
    pub fn row(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }
    
    
    pub fn footer(&mut self, text: &str) {
        // Format the text to be as wide as all the columns combined
        let padding = self.column_padding.iter().sum::<usize>();
        let padded = format!("{:width$}", text, width = padding);

        // Call log in case there are any colors in there
        self.rows.push(padded);
    }


    pub fn add_row<T: Display>(&mut self, elements: Vec<T>) {
        for (index, element) in elements.iter().enumerate() {
            let element = element.to_string();

            let keys = KeyList::new(&element, '<', '>');
            let paris_key_count = keys.map(|s| s.len()).sum::<usize>();
            let len = element.len() - paris_key_count;

            if len > self.longest[index] {
                self.longest[index] = len;
            }
            
            self.columns[index].push(element);
        }
    }


    pub fn display(&mut self) {
        // Print all the collumns with padding
        let row_count = self.columns[0].len();

        // For each column array, print the row:th element
        for i in 0..row_count {
            let mut row = vec![];

            // First row, Third row, and last row are spacers
            if i == 0 {
                self.top_spacer();
            }

            for (j, col) in self.columns.iter().enumerate() {
                let string = &col[i];

                let formatted = colorize_string(string);
                let alignment = self.alignment[j];
                let padding = self.longest[j];

                // Find all the ANSI escape codes in the string to adapt the padding
                let keys = KeyList::new(&formatted, '[', 'm');
                let ansi_char_count = keys.map(|s| s.len()).sum::<usize>();

                // First column's first character opens the table row
                let start_char = if j == 0 { "|" } else { "" };

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
                let ready = format!("<black>{}</> {} <black>|</>", start_char, aligned);

                row.push(ready);
            }

            log!("{}", row.join(""));

            // Last row, second row spacers
            if i == 0 {
                self.middle_spacer();
            }

            if i == row_count - 1 {
                self.bottom_spacer();
            }
        }
    }



    fn top_spacer(&self) {
        self.spacer('┌', '┬', '┐');
    }


    fn middle_spacer(&self) {
        self.spacer('├', '┼', '┤');
    }

    fn bottom_spacer(&self) {
        self.spacer('└', '┴', '┘');
    }


    fn spacer(&self, c_start: char, c_mid: char, c_end: char) {
        let mut pieces = vec![];

        let left = format!("{}{}", c_start, "-".repeat(self.longest[0] + 2));
        let right = format!("{}{}{}", c_mid, "-".repeat(self.longest[self.column_count - 1] + 2), c_end);

        pieces.push(left);
        for i in 1..self.column_count - 1 {
            pieces.push(format!("{}{}", c_mid, "-".repeat(self.longest[i] + 2)));
        }
        pieces.push(right);

        log!("<black>{}</>", pieces.join(""));
    }


    fn pad_strings(&self, strings: &Vec<&str>) -> String {
        let padded = strings
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