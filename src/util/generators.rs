use rand::{ distributions::Alphanumeric, Rng };



/// Generate a random alphanumeric string
/// of a given length.
pub fn random_string(len: usize) -> String {
    rand::thread_rng()
        .sample_iter(&Alphanumeric)
        .take(len)
        .map(char::from)
        .collect()
}