-- name: CreateUser :one
insert into users (id, username, password) values ($1, $2, $3) returning *;

-- name: GetUserByName :one
select * from users where username=$1;

-- name: UpdateRefreshToken :one
update users set refresh_token=$1 where username=$2 returning *;

-- name: UpdateUsername :one
update users set username=$1 where id=$2 returning *;

-- name: UpdateProfile :one
update users set profile_picture=$1 where id=$2 returning *;
