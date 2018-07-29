package main

const (
	jwtKeyUID           = "sub"
	traefikFrontendName = "modoki"
	traefikBackendName  = "modoki_backend"

	frontendFormat = "modokif_%d"
	backendFormat  = "modokib_%d"
	serverName     = "main"

	dockerLabelModokiID   = "com.cs3238.modoki.id"
	dockerLabelModokiUID  = "com.cs3238.modoki.uid"
	dockerLabelModokiName = "com.cs3238.modoki.name"

	// user.go
	defaultShellKVFormat = "modoki/users/%d/defaultShell"
)

const containerSchema = `
CREATE TABLE IF NOT EXISTS containers (
	id INT NOT NULL AUTO_INCREMENT,
	cid VARCHAR(128) UNIQUE,
	name VARCHAR(64) NOT NULL UNIQUE,
	uid VARCHAR(128) NOT NULL,
	status VARCHAR(32),
	message TEXT,
	defaultShell TEXT,
	PRIMARY KEY (id),
	INDEX(cid, name, uid)
);`

const authorizedKeysSchema = `
CREATE TABLE IF NOT EXISTS authorizedKeys (
	id INT NOT NULL AUTO_INCREMENT,
	uid VARCHAR(128) NOT NULL,
	label VARCHAR(32),
	` + "`key`" + ` TEXT,
	PRIMARY KEY(id),
	INDEX(uid, label)
);`
