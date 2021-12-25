use paris::formatter::colorize_string;
use key_list::KeyList;
use std::fmt::Display;






pub struct Table {
    column_count: usize,
    paddings: Vec<usize>,
    alignment: Vec<char>,
    columns: Vec<Vec<String>>
}

impl Table {
    pub fn new(
        column_count: usize,
        mut alignment: Vec<char>
    ) -> Self {
        if alignment.len() < column_count {
            alignment.resize(column_count, '<');
        }

        let mut cols = vec![];
        for _ in 0..column_count {
            cols.push(vec![]);
        }

        let mut paddings = vec![];
        for _ in 0..column_count {
            paddings.push(0);
        }

        Table {
            column_count,
            paddings,
            alignment,
            columns: cols
        }
    }


    pub fn header(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }
    
    
    pub fn row(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }
    
    
    pub fn footer(&mut self, _text: &str) {
        // Format the text to be as wide as all the columns combined
        //let padding = self.paddings.iter().sum::<usize>();
        //let padded = format!("{:width$}", text, width = padding);

        // TODO: Make it span the all columns
    }


    pub fn add_row<T: Display>(&mut self, elements: Vec<T>) {
        for (index, element) in elements.iter().enumerate() {
            let element = element.to_string();

            let keys = KeyList::new(&element, '<', '>');
            let paris_key_count = keys.map(|s| s.len()).sum::<usize>();
            let len = element.len() - paris_key_count;

            if len > self.paddings[index] {
                self.paddings[index] = len;
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
                let padded = self.pad_string(&formatted, self.alignment[j], self.paddings[j]);
                // First column's first character opens the table row
                let start_char = if j == 0 { "|" } else { "" };

                // Format one more time with separator
                let ready = format!("<black>{}</> {} <black>|</>", start_char, padded);

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

        let left = format!("{}{}", c_start, "-".repeat(self.paddings[0] + 2));
        let right = format!("{}{}{}", c_mid, "-".repeat(self.paddings[self.column_count - 1] + 2), c_end);

        pieces.push(left);
        for i in 1..self.column_count - 1 {
            pieces.push(format!("{}{}", c_mid, "-".repeat(self.paddings[i] + 2)));
        }
        pieces.push(right);

        log!("<black>{}</>", pieces.join(""));
    }


    fn pad_string(&self, string: &str, alignment: char, padding: usize) -> String {
        // Find all the ANSI escape codes in the string to adapt the padding
        let keys = KeyList::new(string, '[', 'm');
        let ansi_char_count = keys.map(|s| s.len()).sum::<usize>();

        // Increase the padding to make it even since it'll
        // add less because of the invisible characters
        let width = if ansi_char_count > 0 {
            padding + ansi_char_count + 2
        } else {
            padding
        };

        // Format the string with the custom width padding
        match alignment {
            '<' => format!("{:<width$}", string, width = width),
            '>' => format!("{:>width$}", string, width = width),
            'v' | '^' => format!("{:^width$}", string, width = width),
            _ => panic!("Invalid alignment")
        }
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