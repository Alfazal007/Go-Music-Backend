-- name: CreateSong :one
insert into songs (id, name, song_link, user_id) values ($1, $2, $3, $4) returning *;

-- name: GetSongByName :one
select * from songs where name=$1;

-- name: DeleteSong :one
delete from songs where id=$1 and user_id=$2 returning *;
