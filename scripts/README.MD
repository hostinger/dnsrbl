# es_importer

```
This script will fetch data from ElasticSearch and BAN IP addresses in hbl API. If an arg is specified in more than one place, then commandline values override environment variables which override defaults.

required arguments:
  --es_user ES_USER     Elastic http auth username [env var: ES_USER]
  --es_pass ES_PASS     Elastic http auth password [env var: ES_PASS]
  --es_url ES_URL       ElasticSearch url [env var: ES_URL]
  --es_index ES_INDEX   Index name, default: openresty-* [env var: ES_INDEX]
  --hbl_url HBL_URL     HBL API base url [env var: HBL_URL]
  --hbl_key HBL_KEY     HBL API key for auth [env var: HBL_KEY]

optional arguments:
  -h, --help            show this help message and exit
  --es_scheme {https,http}
                        Transport, default: https [env var: ES_SCHEME]
  --es_timeout ES_TIMEOUT
                        ES read timeout, default: 15 [env var: ES_SCHEME]
  --es_port ES_PORT     80 or 443, default: 443 [env var: ES_PORT]
  --ban_threshold BAN_THRESHOLD
                        Count to get banned, integer, default: 1000 [env var: BAN_THRESHOLD]
  --time_window {10,15,30,60}
                        Time window, in minutes [env var: TIME_WINDOW]
  --dry_run             Just print, do not change [env var: DRY_RUN]
```