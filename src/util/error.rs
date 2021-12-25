use std::fmt;


/// Wrapper type to simplify specifying what error we're throwing
pub type Result<T> = std::result::Result<T, CarbonError>;


/// Enum representing all the possible errors that can be thrown
/// by Carbon. Most of them accept a string as a parameter, which
/// is usually a piece that goes into the error message, not the
/// entire message itself.
pub enum CarbonError {
    /// Given service name doesn't exist in the active environment
    ServiceNotDefined(String),

    /// Given service name is not a running service
    ServiceNotRunning(String),

    /// Given service name is already running
    ServiceAlreadyRunning(String),
    
    /// Variable not defined in active environment
    UndefinedEnvVar(String, String),

    /// There is no active environment
    NoActiveEnv,
    
    /// There was an error reading a file
    FileReadError(String),

    /// There was an error writing a file
    FileWriteError(String),

    /// Failed starting a docker service
    DockerServiceStartup(String),

    /// Failed stopping a docker service
    DockerServiceShutdown(String),

    /// Failed creating a docker network
    DockerNetworkCreate(String),

    /// Failed removing a docker network
    DockerNetworkRemove(String),

    /// Failed inspecting a docker network
    DockerNetworkInspect(String),

    /// Failed connecting to a docker network
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
            NoActiveEnv => format!("No active environment file found, aborting..."),
            
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