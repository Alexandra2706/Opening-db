CREATE TYPE status AS ENUM ('anons', 'ongoing', 'released');

CREATE TYPE kind AS ENUM ('tv', 'movie', 'ova');

CREATE TYPE rating AS ENUM ('r_plus', 'pg_13', 'r', 'g', 'rx');

CREATE TYPE kind_video AS ENUM ('pv', 'op', 'cm', 'ed');

CREATE TABLE IF NOT EXISTS Images_table(
    hash varchar(64), --IPFS CID
    source_img varchar(255), --URL адрес изображения большого размера на сайте shikimori
    meta jsonb --различные данные
);

CREATE TABLE IF NOT EXISTS IPFS_Object(
    hash varchar(64), --IPFS CID
    mime_type varchar(64)
);

CREATE TABLE IF NOT EXISTS Genres_table(
    id varchar(32) PRIMARY KEY, -- id жанра uuid
    shikimori_id integer, --id с сайта shikimori
    genre_name varchar(100), --название жанра на английском
    russian varchar(100) --название жанра на русском
);

CREATE TABLE IF NOT EXISTS Studio_table(
    id uuid, --id студии
    shikimori_id integer, --id с сайта shikimori
    studio_name varchar(100), --название студии
    filtered_name varchar(100), --надо?
    --real boolean, --надо? и что это?
    image text --url логотипа студии
);

CREATE TABLE IF NOT EXISTS Video_table(
    hash varchar(64), --id видео IPFS CID
    shikimori_id integer, --id с сайта shikimori
    url varchar(255), --ссылка на youtube видео
    image_url varchar(255),  --ссылка на youtube картинка
    player_url varchar(255),  --ссылка на youtube видео на весь экран
    video_name varchar(100), --название эпизода
    video_kind kind_video, --тип видео
    hosting varchar(100)
);

CREATE TABLE IF NOT EXISTS Screenshots_table(
    hash varchar(64), --id скриншота IPFS CID
    original varchar(255), --url исходного изображения
    preview varchar(255) --url превью
);

CREATE TABLE IF NOT EXISTS Animes(
    -- Основные поля:
    id uuid PRIMARY KEY, --наш уникальный id
    anime_name varchar(255) UNIQUE NOT NULL, --название анимэ
    russian varchar(255), --название анимэ на русском
    english varchar[], --название анимэ на английском
    japanese varchar[], --название анимэ на японском
    synonyms varchar[], --синонимы названия анимэ
    anime_status status NOT NULL, --статус: anons, ongoing, released
    episodes integer DEFAULT 0, --количество серий
    episodes_aired integer DEFAULT 0, --количество вышедших эпизодов
    aired_on timestamp with time zone, --начало выпуска, формат ISO 8601 with TimeZone
    released_on timestamp with time zone, --конец выпуска, формат ISO 8601 with TimeZone или тут должен быть text?
    duration integer, --длительность серии в минутах
    licensors jsonb, --лицензировано
    franchise jsonb, --франшиза
    updated_at timestamp with time zone DEFAULT NOW(), --дата обновления, формат ISO 8601 with TimeZone
    next_episode_at varchar(255), --следующая серия ссылка
    genres varchar(32)[], --REFERENCES Genre (id), --жанры, может быть несколько
    studios uuid[], --REFERENCES Studio (id), --студии, может быть несколько
    videos varchar(64)[], --REFERENCES Video (id), --эпизоды
    screenshots varchar(64)[], -- REFERENCES Sreenshot (id), --кадры

    -- shikimori data:
    shikimori_id integer UNIQUE NOT NULL, --id с сайта shikimori
    --image integer REFERENCES Image (hash), --изображения на сайте shikimori
    shikimori_kind kind NOT NULL, --тип анимэ на сайте shikimori
    shikimori_rating rating, --возрастной ценз
    shikimori_description varchar, --описание на сайте shikimori
    shikimori_description_html varchar, --описание с тегами html на сайте shikimori
    shikimori_last_revision timestamp with time zone, --дата обновления на сайте shikimori, формат ISO 8601 with TimeZone

    -- myanimelist data:
    myanimelist_id integer UNIQUE NOT NULL, --id с сайта myanimelist
    myanimelist_score real --рейтинг берется из myanimelist

    --description_source null, --Пока опускаем не понятно, что это
);


CREATE TABLE emp (
                     empname text,
                     salary integer,
                     last_date timestamp,
                     last_user text
);

    --genres varchar(32)[], --REFERENCES Genre (id), --жанры, может быть несколько
    --studios uuid[], --REFERENCES Studio (id), --студии, может быть несколько
    --videos integer[], --REFERENCES Video (id), --эпизоды
    --screenshots integer[], -- REFERENCES Sreenshot (id), --кадры

CREATE OR REPLACE FUNCTION check_values() RETURNS trigger AS $check_values$
     DECLARE
       --Создаем переменные для записи результата SELECT
       --и циклов
        add_genre varchar(32) := NULL;
        genre varchar(32);
        add_studio uuid := NULL;
        studio uuid;
        add_video varchar(64) := NULL;
        video varchar(64);
        add_screenshot varchar(64) := NULL;
        screenshot varchar(64);
BEGIN
    --Создаем дату обновления в формате timestamp with time zone
    NEW.updated_at := current_timestamp;

    --Проверить что id_genre задан
    IF NEW.genres IS NOT NULL THEN
        FOR genre IN NEW.genres.first..NEW.genres.last
        LOOP
            SELECT id INTO add_genre FROM Genres_table WHERE id = genre;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'genre % not found', genre;
            END IF;
        end LOOP;
    END IF;

    --Проверить что id_studio задан
    IF NEW.studios IS NOT NULL THEN
        FOR studio IN NEW.studios.first..NEW.studios.last
        LOOP
            SELECT id INTO add_studio FROM Studio_table WHERE id = studio;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'studio % not found', studio;
            END IF;
        end LOOP;
    END IF;

    --Проверить что id_video задан
    IF NEW.videos IS NOT NULL THEN
        FOR video IN NEW.videos.first..NEW.videos.last
        LOOP
            SELECT hash INTO add_video FROM Video_table WHERE hash = video;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'video % not found', video;
            END IF;
        end LOOP;
        END IF;

    --Проверить что id_screenshot задан
    IF NEW.screenshots IS NOT NULL THEN
        FOR screenshot IN NEW.screenshots.first..NEW.screenshots.last
        LOOP
            SELECT hash INTO add_screenshot FROM Screenshots_table WHERE hash = video;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'screenshot % not found', screenshot;
            END IF;
        end LOOP;
    END IF;

RETURN NEW;
END;
$check_values$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER check_genre BEFORE INSERT OR UPDATE ON Animes
    FOR EACH ROW EXECUTE PROCEDURE check_values();



