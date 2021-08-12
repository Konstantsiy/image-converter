-- sing up a user
insert into converter.users (email, password) values ('email1', 'password1');


-- view user password for sign in
select password from converter.users where users.email = 'email1';


-- view user requests
select i.name, r.source_format, r.target_format, r.status, r.created, r.updated
from converter.images i
join converter.requests r
on i.id = r.source_id
where r.user_id = '7186afcc-cae7-11eb-80ff-0bc45a674b3c';


-- view all users uploaded images
select i.name, i.format, i.location
from converter.images i
join converter.requests r
on r.user_id = '7186afcc-cae7-11eb-80ff-0bc45a674b3c'
and i.id = r.source_id;


-- get the user's needed image
select i.name, i.format, i.location
from converter.images i
join converter.requests r
on i.id = r.source_id or i.id = r.target_id
where i.id = '871fc030-cae7-11eb-80ff-0bc45a674b3c'
and r.user_id = '7186afcc-cae7-11eb-80ff-0bc45a674b3c';


-- upload an image and add the converted one
insert into converter.images (name, format)
values ('image1', 'png'),
       ('image2', 'jpeg');


-- create a request
insert into converter.requests (user_id, source_id, target_id, source_format, target_format, ratio, status)
values (
           '7186afcc-cae7-11eb-80ff-0bc45a674b3c',
           '871fc030-cae7-11eb-80ff-0bc45a674b3c',
           '871fc031-cae7-11eb-80ff-0bc45a674b3c',
           'png', 'jpeg', 2, 'queued'
       );


-- get the user with the most number of conversion requests
select r.user_id, count(r.id) as c from converter.requests r
group by r.user_id order by c desc limit 1;


-- get the list of users sorted by the number of requests
select r.user_id, count(r.id) as c from converter.requests r
group by r.user_id order by c;


-- get the list of users which did not request to compress images
select id from converter.users
where id not in (
    select u.id from converter.users u
    join converter.requests r
    on u.id = r.user_id and r.ratio is not null
);