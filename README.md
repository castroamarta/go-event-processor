# Event Processor

Implements a Go Command Line Interface that parses a stream of events and produces an aggregated output.

### Event *translation_delivered*


```events.json
{
	"timestamp": "2018-12-26 18:12:19.903159",
	"translation_id": "5aa5b2f39f7254a75aa4",
	"source_language": "en",
	"target_language": "fr",
	"client_name": "airliberty",
	"event_name": "translation_delivered",
	"duration": 20,
	"nr_words": 100
}
```

For each minute, to get the moving average delivery time of all translations for the past 30 minutes run

`go-event-processor --file_name events.json --window_size 30`

Accepts an `events.json` input file with format, with the input lines ordered by the `timestamp` key, from lower (oldest) to higher values
:

```
{"timestamp": "2024-11-12 11:20:00.908962","translation_id": "5aa5b2f39f7254a75aa5","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 20}
{"timestamp": "2024-11-12 11:21:00.908962","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}
{"timestamp": "2024-11-12 11:25:00.908962","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}
{"timestamp": "2024-11-12 11:31:00.908962","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}
{"timestamp": "2024-11-12 11:32:00.908962","translation_id": "5aa5b2f39f7254a75aa4","source_language": "en","target_language": "fr","client_name": "airliberty","event_name": "translation_delivered","nr_words": 30, "duration": 31}
{"timestamp": "2024-11-12 11:37:00.908962","translation_id": "5aa5b2f39f7254a75bb3","source_language": "en","target_language": "fr","client_name": "taxi-eats","event_name": "translation_delivered","nr_words": 100, "duration": 54}
{"timestamp": "2024-11-12 11:50:00.908962","translation_id": "5aa5b2f39f7254a75bb3","source_language": "en","target_language": "fr","client_name": "taxi-eats","event_name": "translation_delivered","nr_words": 100, "duration": 54}
```

The outputs a `average_delivery_times.json` file in the following format.

```
[
 {
  "average_delivery_time": 20,
  "date": "2024-11-12 11:20:00"
 },
 {
  "average_delivery_time": 31,
  "date": "2024-11-12 11:21:00"
 },
 {
  "average_delivery_time": 31,
  "date": "2024-11-12 11:25:00"
 },
 {
  "average_delivery_time": 31,
  "date": "2024-11-12 11:31:00"
 },
 {
  "average_delivery_time": 31,
  "date": "2024-11-12 11:32:00"
 },
 {
  "average_delivery_time": 54,
  "date": "2024-11-12 11:37:00"
 },
 {
  "average_delivery_time": 54,
  "date": "2024-11-12 11:50:00"
 }
]
```