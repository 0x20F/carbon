use std::fmt;



pub type Result<T> = std::result::Result<T, CarbonError>;


pub enum CarbonError {
    ServiceNotDefined(String),
    ServiceNotRunning(String),
    
    UndefinedEnvVar(String, String),
    
    FileReadError(String),
    FileWriteError(String),

    DockerServiceStartup(String),
    DockerServiceShutdown(String)
}


impl fmt::Display for CarbonError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        use CarbonError::*;

        let string = match self {
            ServiceNotDefined(s) => format!("The <cyan>service</> you tried to access <bright-red>doesn't exist</>: <magenta>{}</>", s),
            ServiceNotRunning(s) => format!("The <cyan>service <magenta>{}</> isn't running. <bright-green>Try starting it first.</>", s),

            UndefinedEnvVar(s, p) => format!("The <b><yellow>{}</> environment <cyan>variable</> isn't defined in the provided dotenv file: <magenta>{}</>", s, p),
            
            FileReadError(s) => format!("Couldn't read <cyan>service file</>: <magenta>{}</>", s),
            FileWriteError(s) => format!("Couldn't write new <cyan>composed service file</>: <magenta>{}</>", s),
            
            DockerServiceStartup(stderr) => format!("Couldn't start services. <cyan>Docker info below</>:\n{}", stderr),
            DockerServiceShutdown(stderr) => format!("Couldn't stop services. <cyan>Docker info below</>:\n{}", stderr)
        };

        write!(f, "{}", string)
    }
}