CREATE TABLE Projects (
    id serial PRIMARY KEY, 
    title TEXT UNIQUE
);

CREATE TABLE Author (
    id serial PRIMARY KEY, 
    name TEXT UNIQUE
);

CREATE TABLE Issue (
    id serial PRIMARY KEY,
    projectId INT NOT NULL,
    authorId INT NOT NULL,
    assigneeId INT NOT NULL,
    key TEXT UNIQUE,
    summary TEXT,
    description TEXT,
    type TEXT,
    priority TEXT,
    status TEXT,
    createdTime TIMESTAMP WITHOUT TIME ZONE,
    closedTime TIMESTAMP WITHOUT TIME ZONE,
    updatedTime TIMESTAMP WITHOUT TIME ZONE,
    timeSpent INT,
    FOREIGN KEY (projectId) REFERENCES Projects (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (authorId) REFERENCES Author (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE StatusChanges (
    issueId INT NOT NULL,
    authorId INT NOT NULL,
    changeTime TIMESTAMP WITHOUT TIME ZONE,
    fromStatus TEXT,
    toStatus TEXT,
    FOREIGN KEY (issueId) REFERENCES Issue (id) ON DELETE CASCADE ON UPDATE CASCADE,
    FOREIGN KEY (authorId) REFERENCES Author (id) ON DELETE CASCADE ON UPDATE CASCADE
);