CREATE DATABASE db_forum;

CREATE USER forum_admin WITH password 'forum_password';

ALTER DATABASE db_forum OWNER TO forum_admin;

\connect db_forum

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;

CREATE TABLE Forum_User(
  id SERIAL PRIMARY KEY NOT NULL,
  nickname CITEXT UNIQUE NOT NULL,
  email CITEXT UNIQUE NOT NULL,
  fullname TEXT NOT NULL,
  about TEXT
);

CREATE TABLE Forum(
  id SERIAL PRIMARY KEY NOT NULL,
  slug CITEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  user_id INTEGER REFERENCES Forum_User(id) NOT NULL
);

CREATE TABLE Thread(
  id SERIAL PRIMARY KEY NOT NULL,
  message TEXT NOT NULL,
  title TEXT NOT NULL,
  slug CITEXT UNIQUE,
  created TIMESTAMPTZ DEFAULT transaction_timestamp() NOT NULL,
  forum_id INTEGER REFERENCES Forum(id) NOT NULL,
  user_id INTEGER REFERENCES Forum_User(id) NOT NULL
);

CREATE TABLE Post(
  id SERIAL PRIMARY KEY NOT NULL,
  parent INTEGER REFERENCES Post(id),
  user_id INTEGER REFERENCES Forum_User(id) NOT NULL,
  thread_id INTEGER REFERENCES Thread(id) NOT NULL,
  created TIMESTAMPTZ DEFAULT transaction_timestamp() NOT NULL,
  isEdited BOOLEAN DEFAULT FALSE NOT NULL,
  message TEXT NOT NULL
);

CREATE TABLE Vote(
  id SERIAL PRIMARY KEY NOT NULL,
  user_id INTEGER REFERENCES Forum_User(id) NOT NULL,
  thread_id INTEGER REFERENCES Thread(id) NOT NULL,
  value INTEGER NOT NULL,
  UNIQUE (user_id, thread_id)
);
