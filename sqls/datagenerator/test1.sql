
CREATE OR REPLACE FUNCTION genre_generate() RETURNS VOID AS $genre_generate$
DECLARE
    --Задаем количество записей в таблице жанров
    number integer := 10;
BEGIN
    INSERT INTO genres_table (id, shikimori_id, genre_name, russian)
    SELECT
        md5(random()::text)::char(32) AS id,
        s AS shikimori_id,
        md5(random()::text)::char(10) AS genre_name,
        md5(random()::text)::char(10) AS russian
    FROM
        generate_series(1, number) s;

END;
$genre_generate$ LANGUAGE plpgsql;

SELECT genre_generate();

-- DO $$ BEGIN
--     INSERT INTO genres_table (id, shikimori_id, genre_name, russian)
--     SELECT
--         md5(random()::text)::char(32) AS id,
--         s AS shikimori_id,
--         md5(random()::text)::char(10) AS genre_name,
--         md5(random()::text)::char(10) AS russian
--     FROM
--         generate_series(1, 10) s;
-- END $$;

SELECT * FROM genres_table;



