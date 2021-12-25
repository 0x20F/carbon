use std::fmt;



pub type Result<T> = std::result::Result<T, CarbonError>;


pub enum CarbonError {
    ServiceNotDefined(String),
    ServiceNotRunning(String),
    ServiceAlreadyRunning(String),
    
    UndefinedEnvVar(String, String),
    
    FileReadError(String),
    FileWriteError(String),

    DockerServiceStartup(String),
    DockerServiceShutdown(String),
    DockerNetworkCreate(String),
    DockerNetworkRemove(String),
    DockerNetworkInspect(String),
    DockerNetworkConnect(String),
}


impl fmt::Display for CarbonError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        use CarbonError::*;

        let string = match self {
            ServiceNotDefined(s) => format!("The <cyan>service</> you tried to access <bright-red>doesn't exist</>: <magenta>{}</>", s),
            ServiceNotRunning(s) => format!("The <cyan>service <magenta>{}</> isn't running. <bright-green>Try starting it first.</>", s),
            ServiceAlreadyRunning(s) => format!("The <cyan>service <magenta>{}</> is already running, aborting...\n  <bright-green>Try stopping it first.</>", s),
            
            UndefinedEnvVar(s, p) => format!("The <b><yellow>{}</> environment <cyan>variable</> isn't defined in the provided dotenv file: <magenta>{}</>", s, p),
            
            FileReadError(s) => format!("Couldn't read <cyan>service file</>: <magenta>{}</>", s),
            FileWriteError(s) => format!("Couldn't write new <cyan>composed service file</>: <magenta>{}</>", s),
            
            DockerServiceStartup(stderr) => format!("Couldn't start services. <cyan>Docker info below</>:\n{}", stderr),
            DockerServiceShutdown(stderr) => format!("Couldn't stop services. <cyan>Docker info below</>:\n{}", stderr),
            DockerNetworkCreate(stderr) => format!("Couldn't create new network! <cyan>Docker info below</>:\n{}", stderr),
            DockerNetworkRemove(stderr) => format!("Couldn't remove network! <cyan>Docker info below:</>\n{}", stderr),
            DockerNetworkInspect(stderr) => format!("Couldn't inspect network! <cyan>Docker info below:</>\n{}", stderr),
            DockerNetworkConnect(stderr) => format!("Couldn't connect containers to network! <cyan>Docker info below:</>\n{}", stderr),
        };

        write!(f, "{}", string)
    }
}