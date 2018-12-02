/* Execute this script as database owner */
/* psql -h localhost -U postgres -p 5432 -f assets/deployment/database/schema.sql */

drop database todolist;
create database todolist;
create user todolist_user;
create user todolist_admin;
alter user todolist_admin with encrypted password 'todo101';
alter user todolist_user with encrypted password 'todo42';
revoke connect on database todolist from public;
grant all privileges on database todolist to todolist_admin;
alter database todolist owner to todolist_admin;

grant connect on database todolist to todolist_user;

\c todolist todolist_admin

grant select, insert, update, delete on all tables in schema public to todolist_user;
alter default privileges in schema public grant all on tables to todolist_user;
alter default privileges in schema public grant all on sequences to todolist_user;

drop table tasks;

create table tasks(
   id serial primary key,
   name varchar(50) not null,
   description text,
   status varchar(10) not null,
   duedate timestamp,
   created timestamp not null default current_timestamp
);
create index tasks_status on tasks (status);









