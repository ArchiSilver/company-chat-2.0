-- Down migration: Drop messenger tables

DROP TABLE IF EXISTS messenger.messages;
DROP TABLE IF EXISTS messenger.chat_members;
DROP TABLE IF EXISTS messenger.chats;
DROP SCHEMA IF EXISTS messenger CASCADE;
