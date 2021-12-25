/// Helper macro to unwrap a std::process::Command output
/// into a string or throwing an error if the command wrote
/// to stderr.
macro_rules! unwrap_stderr {
    ( $output: expr, $enum_type: ident ) => {
        if !$output.status.success() {
            let stderr = std::str::from_utf8(&$output.stderr).unwrap();
            return Err(CarbonError::$enum_type(stderr.to_string()));
        } else {
            Ok(())
        }
    }
}


pub(crate) use unwrap_stderr;