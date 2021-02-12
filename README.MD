# DNS RBL tool

`requires python3+`

```
required arguments:
  --es_user ES_USER     Elastic http auth username [env var: ES_USER]
  --es_pass ES_PASS     Elastic http auth password [env var: ES_PASS]
  --es_url ES_URL       ElasticSearch url [env var: ES_URL]
  --es_index ES_INDEX   Index name, default: openresty-* [env var: ES_INDEX]
  --pdns_api_key PDNS_API_KEY
                        PowerDNS API KEY [env var: PDNS_API_KEY]
  --pdns_rbl_zone PDNS_RBL_ZONE
                        RBL zone, default: hostinger.rbl [env var: PDNS_RBL_ZONE]
  --mysql_host MYSQL_HOST
                        Mysql DB host/ip [env var: MYSQL_HOST]
  --mysql_pass MYSQL_PASS
                        Mysql password [env var: MYSQL_PASS]
  --mysql_user MYSQL_USER
                        Mysql user [env var: MYSQL_USER]
  --mysql_db MYSQL_DB   Mysql database [env var: MYSQL_DB]

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
  --retention RETENTION
                        Flush bans after, in days, integer, default: 10 [env var: RETENTION]
  --pdns_api_url PDNS_API_URL
                        PowerDNS API URL: http://fqdn:8081/api/v1 [env var: PDNS_API_URL]
  --manual_remove MANUAL_REMOVE [MANUAL_REMOVE ...]
                        list of IPS to remove from rbl [env var: MANUAL_REMOVE]
  --list_banned         list of banned IPS from rbl [env var: LIST_BANNED]
  --dry_run             Just print, do not change [env var: DRY_RUN]

```