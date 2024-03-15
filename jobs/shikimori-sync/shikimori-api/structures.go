package shikimori_api

import "time"

type Anime struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Russian string `json:"russian"`
	Image   struct {
		Original string `json:"original"`
		Preview  string `json:"preview"`
		X96      string `json:"x96"`
		X48      string `json:"x48"`
	} `json:"image"`
	Url               string        `json:"url"`
	Kind              string        `json:"kind"`
	Score             string        `json:"score"`
	Status            string        `json:"status"`
	Episodes          int           `json:"episodes"`
	EpisodesAired     int           `json:"episodes_aired"`
	AiredOn           string        `json:"aired_on"`
	ReleasedOn        string        `json:"released_on"`
	Rating            string        `json:"rating"`
	English           []string      `json:"english"`
	Japanese          []string      `json:"japanese"`
	Synonyms          []interface{} `json:"synonyms"`
	LicenseNameRu     string        `json:"license_name_ru"`
	Duration          int           `json:"duration"`
	Description       string        `json:"description"`
	DescriptionHtml   string        `json:"description_html"`
	DescriptionSource interface{}   `json:"description_source"`
	Franchise         string        `json:"franchise"`
	Favoured          bool          `json:"favoured"`
	Anons             bool          `json:"anons"`
	Ongoing           bool          `json:"ongoing"`
	ThreadId          int           `json:"thread_id"`
	TopicId           int           `json:"topic_id"`
	MyanimelistId     int           `json:"myanimelist_id"`
	RatesScoresStats  []struct {
		Name  int `json:"name"`
		Value int `json:"value"`
	} `json:"rates_scores_stats"`
	RatesStatusesStats []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"rates_statuses_stats"`
	UpdatedAt     time.Time   `json:"updated_at"`
	NextEpisodeAt interface{} `json:"next_episode_at"`
	Fansubbers    []string    `json:"fansubbers"`
	Fandubbers    []string    `json:"fandubbers"`
	Licensors     []string    `json:"licensors"`
	Genres        []struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Russian   string `json:"russian"`
		Kind      string `json:"kind"`
		EntryType string `json:"entry_type"`
	} `json:"genres"`
	Studios []struct {
		Id           int    `json:"id"`
		Name         string `json:"name"`
		FilteredName string `json:"filtered_name"`
		Real         bool   `json:"real"`
		Image        string `json:"image"`
	} `json:"studios"`
	Videos []struct {
		Id        int    `json:"id"`
		Url       string `json:"url"`
		ImageUrl  string `json:"image_url"`
		PlayerUrl string `json:"player_url"`
		Name      string `json:"name"`
		Kind      string `json:"kind"`
		Hosting   string `json:"hosting"`
	} `json:"videos"`
	Screenshots []struct {
		Original string `json:"original"`
		Preview  string `json:"preview"`
	} `json:"screenshots"`
	UserRate struct {
		Id        int         `json:"id"`
		Score     int         `json:"score"`
		Status    string      `json:"status"`
		Text      string      `json:"text"`
		Episodes  int         `json:"episodes"`
		Chapters  interface{} `json:"chapters"`
		Volumes   interface{} `json:"volumes"`
		TextHtml  string      `json:"text_html"`
		Rewatches int         `json:"rewatches"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt time.Time   `json:"updated_at"`
	} `json:"user_rate"`
}
