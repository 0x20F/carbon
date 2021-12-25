use paris::formatter::colorize_string;
use key_list::KeyList;
use std::fmt::Display;






/// Table structure to allow for easy printing
/// and formatting of data in a table.
///
/// Built in accounting for any Paris color keys
/// that may be added to text, as well as scalable
/// column width based on automatic detection of the
/// longest value in a given column.
///
/// Default alignment is "<" for left alignment.
pub struct Table {
    /// How many columns the table will have
    column_count: usize,

    /// The column widths of each column
    paddings: Vec<usize>,

    /// Text alignment in each of the columns
    /// valid values are: "<", "^", ">"
    alignment: Vec<char>,

    /// The column data to be printed
    columns: Vec<Vec<String>>
}

impl Table {
    /// Create a new table with a given number of columns
    /// If custom alignment is provided, it will be used
    /// for each provided column. If not, the default alignment is set.
    pub fn new(
        column_count: usize,
        mut alignment: Vec<char>
    ) -> Self {
        if alignment.len() < column_count {
            alignment.resize(column_count, '<');
        }

        // All columns start with an empty vector.
        let mut cols = vec![];
        for _ in 0..column_count {
            cols.push(vec![]);
        }

        // All columns start with a padding of 0
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


    /// Add the first row to the table.
    pub fn header(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }
    
    
    /// Add a new row to the table.
    pub fn row(&mut self, elements: Vec<&str>) {
        self.add_row(elements);
    }


    /// Display the table.
    /// This is where all the provided data gets properly
    /// padded and aligned, as well as where all the coloring
    /// that might be provided by the Paris color keys is applied.
    pub fn display(&mut self) {
        let row_count = self.columns[0].len();

        // For each column array, print the row:th element
        for i in 0..row_count {
            let mut row = vec![];

            // First row is a spacer
            if i == 0 { self.top_spacer(); }

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
            if i == 0 { self.middle_spacer(); }
            if i == row_count - 1 { self.bottom_spacer(); }
        }
    }


    /// Create a new spacer with edge pieces
    /// matching the top of a rectangle.
    fn top_spacer(&self) {
        self.spacer('┌', '┬', '┐');
    }


    /// Create a new spacer with edge pieces
    /// matching the middle of a rectangle.
    fn middle_spacer(&self) {
        self.spacer('├', '┼', '┤');
    }


    /// Create a new spacer with edge pieces
    /// matching the bottom of a rectangle.
    fn bottom_spacer(&self) {
        self.spacer('└', '┴', '┘');
    }


    /// Creates a new spacer by looking at all
    /// the generated column paddings and repeating
    /// a set of dashes to match how wide each column will be.
    fn spacer(&self, c_start: char, c_mid: char, c_end: char) {
        let mut pieces = vec![];

        let left = format!("{}{}", c_start, "—".repeat(self.paddings[0] + 2));
        let right = format!("{}{}{}", c_mid, "—".repeat(self.paddings[self.column_count - 1] + 2), c_end);

        pieces.push(left);
        for i in 1..self.column_count - 1 {
            pieces.push(format!("{}{}", c_mid, "—".repeat(self.paddings[i] + 2)));
        }
        pieces.push(right);

        log!("<black>{}</>", pieces.join(""));
    }


    /// Add a new row to the table.
    /// Looks through all the provided columns that will
    /// be part of the row and determines whether or not
    /// any of the elements are greater than the already
    /// defined padding for that column.
    ///
    /// This is what "automatically" scales the width of a column
    /// based on input values.
    fn add_row<T: Display>(&mut self, elements: Vec<T>) {
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


    /// Formats a string based on the provided alignment
    /// and padding. Making sure to account for any ANSI characters
    /// that may be present in the string from any added colors since
    /// those can mess up the padding that `format!()` does.
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