package structures

import (
	"github.com/google/uuid"
)

type AnimeInfo struct {
	Id            uuid.UUID                    `json:"id"`
	AnimeName     string                       `json:"anime_name"`
	NameRussian   string                       `json:"name_russian"`
	NameEnglish   *[]string                    `json:"name_english"`
	NameJapanese  *[]string                    `json:"name_japanese"`
	NameSynonyms  *[]string                    `json:"name_synonyms"`
	AnimeStatus   ValidatableEnum[AnimeStatus] `json:"anime_status"`
	Episodes      *int                         `json:"episodes"`
	EpisodesAired *int                         `json:"episodes_aired"`
	AiredOn       Time                         `json:"aired_on"`
	ReleasedOn    Time                         `json:"released_on"`
	Duration      *int                         `json:"duration"`

	//LicensorRu *[]string `json:"licensor_ru"`
	//Franchise  *string   `json:"franchise"`

	UpdatedAt     Time         `json:"updated_at"`
	NextEpisodeAt *string      `json:"next_episode_at"`
	Image         *string      `json:"image"`
	Genres        *[]string    `json:"genres"`
	Studios       *[]uuid.UUID `json:"studios"`
	Videos        *[]string    `json:"videos"`
	Screenshots   *[]string    `json:"screenshots"`

	ShikimoriId              int                          `json:"shikimori_id"`
	ShikimoriKind            ValidatableEnum[AnimeKind]   `json:"shikimori_kind"`
	ShikimoriRating          ValidatableEnum[AnimeRating] `json:"shikimori_rating"`
	ShikimoriDescription     *string                      `json:"shikimori_description"`
	ShikimoriDescriptionHtml *string                      `json:"shikimori_description_html"`
	ShikimoriLastRevision    Time                         `json:"shikimori_last_revision"`

	MyanimelistId    int     `json:"myanimelist_id"`
	MyanimelistScore float32 `json:"myanimelist_score"`
}

/*
licensors_ru jsonb, --лицензировано
franchise jsonb, --франшиза

videos varchar(255)[], --REFERENCES Video (id), --эпизоды

--description_source null, --Пока опускаем не понятно, что это
);
*/

type PersonInfo struct {
	Id          uuid.UUID `json:"id"`
	PeopleName  string    `json:"people_name"`
	Russian     *string   `json:"russian"`
	Japanese    *string   `json:"japanese"`
	Image       *string   `json:"image"`
	ShikimoriId int       `json:"shikimori_id"`
	JobTitle    *string   `json:"job_title"`
	Birthday    *struct {
		Day   *int
		Year  *int
		Month *int
	} `json:"birthday"`
	Deceased *struct {
		Day   *int
		Year  *int
		Month *int
	} `json:"deceased"`
	Website       *string        `json:"website"`
	GrouppedRoles map[string]int `json:"groupped_roles"`
	Producer      *bool          `json:"producer"`
	Mangaka       *bool          `json:"mangaka"`
	Seyu          *bool          `json:"seyu"`
	UpdatedAt     Time           `json:"updated_at"`
}
