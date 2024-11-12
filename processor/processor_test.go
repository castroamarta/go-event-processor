package processor

import (
	"go-event-processor/types"
	"testing"
	"time"

	"github.com/go-playground/assert"
)

func TestAggregate(t *testing.T) {

	t.Run(`Should aggregate events successfully`, func(t *testing.T) {
		event := types.Event{
			Timestamp:      "2018-12-26 18:11:08.509654",
			TranslationID:  "5aa5b2f39f7254a75aa5",
			SourceLanguage: "en",
			TargetLanguage: "fr",
			ClientName:     "airliberty",
			EventName:      "translation_delivered",
			NrWords:        30,
			Duration:       43,
		}
		durations := map[string]types.DeliveryTimes{
			"2018-12-26 18:11:00": {
				DurationSum: 40,
				Count:       4,
			},
			"2018-12-26 18:23:00": {
				DurationSum: 45,
				Count:       2,
			},
		}
		aggregate(&event, durations)
		expected := map[string]types.DeliveryTimes{
			"2018-12-26 18:11:00": {
				DurationSum: 83,
				Count:       5,
			},
			"2018-12-26 18:23:00": {
				DurationSum: 45,
				Count:       2,
			},
		}
		assert.Equal(t, expected, durations)
	})
}

func TestTimestampIsValid(t *testing.T) {

	t.Run(`Should return true if valid`, func(t *testing.T) {
		timestamp, _ := time.Parse(types.Timelayout, "2018-12-26 18:23:19.903158")
		line := "{\"timestamp\": \"2018-12-26 18:23:19.903159\",\"translation_id\": \"5aa5b2f39f7254a75bb3\",\"source_language\": \"en\",\"target_language\": \"fr\",\"client_name\": \"taxi-eats\",\"event_name\": \"translation_delivered\",\"nr_words\": 100, \"duration\": 54}"
		actual := TimestampIsValid(timestamp, line)
		assert.Equal(t, true, actual)

	})
	t.Run(`Should return false if invalid`, func(t *testing.T) {
		timestamp, _ := time.Parse(types.Timelayout, "2018-12-26 18:25:19.903159")
		line := "{\"timestamp\": \"2018-12-26 18:23:19.903159\",\"translation_id\": \"5aa5b2f39f7254a75bb3\",\"source_language\": \"en\",\"target_language\": \"fr\",\"client_name\": \"taxi-eats\",\"event_name\": \"translation_delivered\",\"nr_words\": 100, \"duration\": 54}"
		actual := TimestampIsValid(timestamp, line)
		assert.Equal(t, false, actual)
	})
	t.Run(`Should return false if line parse fails`, func(t *testing.T) {
		timestamp, _ := time.Parse(types.Timelayout, "2018-12-26 18:25:19.903159")
		line := "{\"timestamp\": \"2018-12 18:23:19.903159\",\"translation_id\": \"5aa5b2f39f7254a75bb3\",\"source_language\": \"en\",\"target_language\": \"fr\",\"client_name\": \"taxi-eats\",\"event_name\": \"translation_delivered\",\"nr_words\": 100, \"duration\": 54}"
		actual := TimestampIsValid(timestamp, line)
		assert.Equal(t, false, actual)
	})
}

func TestAv(t *testing.T) {

	t.Run(`Should calculate average successfully`, func(t *testing.T) {
		durations := map[string]types.DeliveryTimes{
			"2018-12-26 18:11:00": {
				DurationSum: 40,
				Count:       4,
			},
			"2018-12-26 18:12:00": {
				DurationSum: 45,
				Count:       2,
			},
		}
		averages := Av(durations)
		expected := []map[string]interface{}{
			{
				"date":                  "2018-12-26 18:11:00",
				"average_delivery_time": float64(40) / float64(4),
			},
			{
				"date":                  "2018-12-26 18:12:00",
				"average_delivery_time": float64(45) / float64(2),
			},
		}
		assert.Equal(t, expected, averages)
	})
}