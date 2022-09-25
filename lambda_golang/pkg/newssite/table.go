package newssite

type EventName string

const (
	// `landing`
	EVENT_LANDING_PAGE_REQUESTED EventName = "LANDING_PAGE_REQUESTED" ✅ 📩
	// +`landing_s3_trigger` (put in db) ✅
	EVENT_LANDING_PAGE_FETCHED EventName = "LANDING_PAGE_FETCHED"
	// @`landing_metadata` -> `landing_metadata_cronjob` (query db; store metadata) ✅ (cronjob trigger) ✅
	EVENT_LANDING_METADATA_REQUESTED EventName = "LANDING_METADATA_REQUESTED"
	// `stories` (metadata triggers; sfn) ✅
	EVENT_LANDING_STORIES_REQUESTED EventName = "LANDING_STORIES_REQUESTED"
	// `story` (sfn map; archive story) ✅ 📩
	EVENT_STORY_REQUESTED EventName = "STORY_REQUESTED"
	// random wait
	EVENT_STORY_FETCHED EventName = "STORY_FETCHED"
	// +`stories_finalizer` (sfn last step)  ✅
	EVENT_LANDING_STORIES_FETCHED EventName = "LANDING_STORIES_FETCHED"
)

type MediaTableItemEvent struct {
	EventName EventName `json:"eventName"`
	Detail    string    `json:"detail"`
	EventTime string    `json:"eventTime"`
}

type DocType string

const (
	DOCTYPE_LANDING DocType = "LANDING"
)
