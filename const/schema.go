package constants

const ContainerSchema = `
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

const AuthorizedKeysSchema = `
CREATE TABLE IF NOT EXISTS authorizedKeys (
	id INT NOT NULL AUTO_INCREMENT,
	uid VARCHAR(128) NOT NULL,
	label VARCHAR(32),
	` + "`key`" + ` TEXT,
	PRIMARY KEY(id),
	INDEX(uid, label)
);`

const APIKeysSchema = `
CREATE TABLE IF NOT EXISTS apiKeys (
	id INT NOT NULL AUTO_INCREMENT,
	uid VARCHAR(128) NOT NULL UNIQUE,
	apiKey VARCHAR(256) NOT NULL,
	PRIMARY KEY(id),
	INDEX(apiKey)
);`
