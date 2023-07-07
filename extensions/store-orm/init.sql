use test;

drop table test_suites;
create table test_suites (
    id int primary key auto_increment,
    name char(20),
    api char(100)
);

drop table test_cases;
create table test_cases (
    id int primary key auto_increment,
    suiteId int,
    suiteName char(20),
    name char(20)
);

insert into test_suites values (1, 'gitlab', '');
insert into test_suites values (2, 'gitee', '');

insert into test_cases values (1, 1, 'gitlab', 'gitlab-1');