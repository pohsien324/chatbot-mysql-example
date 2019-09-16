create database chatbot;
use chatbot ;
create table `telegram`( 
   `id` int not null primary key auto_increment, 
   `keyword` text,
   `response` text  
) CHARSET=utf8 COLLATE=utf8_general_ci;
insert into `telegram`(keyword,response) values ('Hi','Hi there');
