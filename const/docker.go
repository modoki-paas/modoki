package constants

const (
	JWTKeyUID           = "sub"
	TraefikFrontendName = "modoki"
	TraefikBackendName  = "modoki_backend"

	FrontendFormat = "modokif_%d"
	BackendFormat  = "modokib_%d"
	ServerName     = "main"

	DockerLabelModokiID   = "com.cs3238.modoki.id"
	DockerLabelModokiUID  = "com.cs3238.modoki.uid"
	DockerLabelModokiName = "com.cs3238.modoki.name"

	// user.go
	DefaultShellKVFormat = "modoki/users/%s/defaultShell" // TODO: encode for security
)
