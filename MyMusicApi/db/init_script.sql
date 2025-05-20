create table if not exists music (
    Id integer primary key,
    Name text not null,
    Url text not null,
    Path text not null
)