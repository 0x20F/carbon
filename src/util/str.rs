/// If a string is longer than the specified length
/// Get the important end pieces and add an ellipsis in the middle.
pub fn cut(s: &str, len: usize) -> String {
    if s.len() > len {
        format!("{}...{}", &s[0..10], &s[s.len() - 20..])
    } else {
        s.to_string()
    }
}