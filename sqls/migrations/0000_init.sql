--CREATE TYPE IF NOT EXISTS status AS ENUM ('anons', 'ongoing', 'released');
--CREATE TYPE IF NOT EXISTS kind AS ENUM ('tv', 'movie', 'ova', 'ona', 'other');
--CREATE TYPE IF NOT EXISTS rating AS ENUM ('r_plus', 'pg_13', 'r', 'g', 'rx');
--CREATE TYPE IF NOT EXISTS kind_video AS ENUM ('pv', 'op', 'cm', 'ed', 'mv');

DO $$ BEGIN
    CREATE TYPE status AS ENUM ('anons', 'ongoing', 'released');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE kind AS ENUM ('tv', 'movie', 'ova', 'ona', 'other');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE rating AS ENUM ('r_plus', 'pg_13', 'r', 'g', 'rx');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE kind_video AS ENUM ('pv', 'op', 'cm', 'ed', 'mv');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

CREATE TABLE IF NOT EXISTS images_table(
    hash varchar(64) PRIMARY KEY, --IPFS CID
    source_img varchar(255), --URL адрес изображения большого размера на сайте shikimori
    meta jsonb --различные данные
);

CREATE TABLE IF NOT EXISTS ipfs_object(
    hash varchar(64) PRIMARY KEY, --IPFS CID
    mime_type varchar(64)
);

CREATE TABLE IF NOT EXISTS genres_table(
    id varchar(32) PRIMARY KEY, -- id жанра uuid
    shikimori_id integer, --id с сайта shikimori
    genre_name varchar(100), --название жанра на английском
    russian varchar(100) --название жанра на русском
);

CREATE TABLE IF NOT EXISTS studio_table(
    id uuid PRIMARY KEY, --id студии
    shikimori_id integer, --id с сайта shikimori
    studio_name varchar(100), --название студии
    --filtered_name varchar(100), --надо?
    --real boolean, --надо? и что это?
    image varchar(64) REFERENCES images_table (hash) --url логотипа студии
);

CREATE TABLE IF NOT EXISTS video_table(
    hash varchar(64) PRIMARY KEY, --id видео IPFS CID
    shikimori_id integer, --id с сайта shikimori
    url varchar(255), --ссылка на youtube видео
    --image_url varchar(255),  --ссылка на youtube картинка
    player_url varchar(255),  --ссылка на youtube видео на весь экран
    video_name varchar(100), --название эпизода
    video_kind kind_video --тип видео
    --hosting varchar(100)
);

CREATE TABLE IF NOT EXISTS audio_table(
    hash varchar(64) PRIMARY KEY, --id аудио IPFS CID
    source_url varchar(255), --url аудиозаписи
    duration integer --длительность аудио в минутах
);

CREATE TABLE IF NOT EXISTS animes(
    -- Основные поля:
    id uuid PRIMARY KEY, --наш уникальный id
    anime_name varchar(255) UNIQUE NOT NULL, --название анимэ
    name_russian varchar(255), --название анимэ на русском
    name_english varchar[], --название анимэ на английском
    name_japanese varchar[], --название анимэ на японском
    name_synonyms varchar[], --синонимы названия анимэ
    anime_status status NOT NULL, --статус: anons, ongoing, released
    episodes integer, --количество серий
    episodes_aired integer DEFAULT 0, --количество вышедших эпизодов
    aired_on timestamp with time zone, --начало выпуска
    released_on timestamp with time zone, --конец выпуска
    duration integer, --длительность серии в минутах
    licensors_ru jsonb, --лицензировано
    franchise jsonb, --франшиза
    updated_at timestamp with time zone DEFAULT NOW(), --дата обновления
    next_episode_at varchar(255), --следующая серия ссылка
    image varchar(64) REFERENCES images_table (hash), --постер аниме (изображения на сайте shikimori)
    genres varchar(32)[], --жанры, может быть несколько
    studios uuid[], --REFERENCES Studio (id), --студии, может быть несколько
    videos varchar(64)[], --REFERENCES Video (id), --эпизоды
    screenshots varchar(64)[], -- REFERENCES Sreenshot (id), --кадры

    -- shikimori data:
    shikimori_id integer UNIQUE NOT NULL, --id с сайта shikimori
    shikimori_kind kind NOT NULL, --тип анимэ на сайте shikimori
    shikimori_rating rating, --возрастной ценз
    shikimori_description varchar, --описание на сайте shikimori
    shikimori_description_html varchar, --описание с тегами html на сайте shikimori
    shikimori_last_revision timestamp with time zone, --дата обновления на сайте shikimori

    -- myanimelist data:
    myanimelist_id integer UNIQUE NOT NULL, --id с сайта myanimelist
    myanimelist_score real --рейтинг берется из myanimelist

    --description_source null, --Пока опускаем не понятно, что это
);

CREATE TABLE IF NOT EXISTS person(
    id uuid PRIMARY KEY, --уникальный id человека
    people_name VARCHAR(128) NOT NULL, --имя
    russian VARCHAR(128), --имя на русском
    japanese VARCHAR(64), ----имя на японском
    image VARCHAR(64) REFERENCES images_table (hash), --url фото человека
    shikimori_id integer UNIQUE, --id человека на сайте shikimori
    job_title VARCHAR(255), --основная работа
    birthday date, --дата рождения
    deceased date, --дата смерти
    website VARCHAR(255), --адрес сайта человека
    groupped_roles jsonb, --роли в аниме: название + количество
    --roles VARCHAR(255)[], --роли в аниме (Лучшие роли?) {[список аниме]}
    --works VARCHAR(255)[], --сделать таблицу anime_to_person (id_person, id_anime, role)
    producer BOOLEAN DEFAULT FALSE,
    mangaka BOOLEAN DEFAULT FALSE,
    seyu BOOLEAN DEFAULT FALSE,
    updated_at timestamp with time zone DEFAULT NOW() --дата обновления

);

CREATE OR REPLACE FUNCTION anime_validate() RETURNS trigger AS $anime_validate$
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
    --Создаем дату обновления
    NEW.updated_at := current_timestamp;

    --Проверить что id_genre задан верно
    IF NEW.genres IS NOT NULL THEN
        RAISE NOTICE 'NEW.genres = %', NEW.genres;
        --RAISE NOTICE 'NEW.genres.first = %', NEW.genres.first;
        FOR genre IN NEW.genres LOOP
            RAISE NOTICE 'genre = %', genre;
            SELECT id INTO add_genre FROM genres_table WHERE id = genre;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'genre % not found', genre;
            END IF;
        end LOOP;
    END IF;

    --Проверить что id_studio задан верно
    IF NEW.studios IS NOT NULL THEN
        FOR studio IN NEW.studios LOOP
            SELECT id INTO add_studio FROM studio_table WHERE id = studio;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'studio % not found', studio;
            END IF;
        end LOOP;
    END IF;

    --Проверить что id_video задан верно
    IF NEW.videos IS NOT NULL THEN
        FOR video IN NEW.videos LOOP
            SELECT hash INTO add_video FROM video_table WHERE hash = video;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'video % not found', video;
            END IF;
        end LOOP;
    END IF;

    --Проверить что id_screenshot задан верно ВСЕ ДОЛЖНО БЫТЬ В ИМАЖЕ
    IF NEW.screenshots IS NOT NULL THEN
        FOR screenshot IN NEW.screenshots LOOP
            SELECT hash INTO add_screenshot FROM screenshots_table WHERE hash = video;
            IF NOT FOUND THEN
                RAISE EXCEPTION 'screenshot % not found', screenshot;
            END IF;
        END LOOP;
    END IF;

    --Проверить, что количество вышедщих серий не больше общего количества серий
    IF NEW.episodes_aired > NEW.episodes THEN
       RAISE EXCEPTION 'episodes must be lower then episodes_aired';
    END IF;

    --Проверить, что дата начала выпуска анимэ не больше даты конца выпуска
    IF NEW.aired_on > NEW.released_on THEN
        RAISE EXCEPTION 'released_on must be lower then aired_on';
    END IF;

    RETURN NEW;
END;
$anime_validate$ LANGUAGE plpgsql;


CREATE OR REPLACE TRIGGER check_genre BEFORE INSERT OR UPDATE ON animes
    FOR EACH ROW EXECUTE PROCEDURE anime_validate();



