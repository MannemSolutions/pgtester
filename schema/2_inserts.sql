insert into t values(1);
insert into t2 values('["question","answer"]'::json),
    ('{"names":["me","myself","I"],"person":["second","third","first"]}'::json),
    ('{"answer":"fourtytwo"}'::json);
insert into t3 values('["question","answer"]'::json),
    ('{"names":["me","myself","I"],"person":["second","third","first"]}'::json),
    ('{"answer":"fourtytwo"}'::json);
insert into t4 values(true);
insert into t5 values('2024-01-11'::timestamp),
    (null);
insert into t6 values(1.0/134.0),
    (1.21);
insert into t7 values('192.168.0.1'::inet, '255.255.0.0'::cidr);
insert into t8 values('1 days 22 seconds'::interval),
    ('1 year 221 microseconds'::interval),
    ('1 year 221 milliseconds'::interval);
insert into t9 values(0.001);
insert into t10 values('PgTester');
insert into t11 values('{ 1, 2, 3 }');
insert into t12 values('<list><item1/></list>'::xml);