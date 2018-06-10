package main

const jwtKeyUID = "uid"

const (
	TraefikFrontendName = "modoki"
	TraefikBackendName  = "modoki_backend"

	FrontendFormat = "modokif_%d"
	BackendFormat  = "modokib_%d"
	ServerName     = "main"

	DockerLabelModokiID   = "com.cs3238.modoki.id"
	DockerLabelModokiUID  = "com.cs3238.modoki.uid"
	DockerLabelModokiName = "com.cs3238.modoki.name"
)

const containerSchema = `
CREATE TABLE IF NOT EXISTS containers (
	id INT NOT NULL AUTO_INCREMENT,
	cid VARCHAR(128) UNIQUE,
	name VARCHAR(64) NOT NULL UNIQUE,
	uid INT NOT NULL,
	status VARCHAR(32),
	message TEXT,
	PRIMARY KEY (id),
	INDEX(cid, name, uid)
);`
