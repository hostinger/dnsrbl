from elasticsearch import Elasticsearch, ConnectionError
from elasticsearch_dsl import Search, Q
import powerdns
import datetime
import configargparse
import mysql.connector
import ipaddress


def parse_args():
    parser = configargparse.ArgParser(add_help=False, description='This script will fetch data from ElasticSearch '
                                                                  'and BAN IP addresses in hostinger DNSRBL.')
    required = parser.add_argument_group('required arguments')
    optional = parser.add_argument_group('optional arguments')
    '''Suppressing default help'''
    optional.add_argument('-h', '--help', action='help', default=configargparse.SUPPRESS,
                          help='show this help message and exit')
    required.add_argument('--es_user', help='Elastic http auth username', required=True, env_var='ES_USER')
    required.add_argument("--es_pass", help='Elastic http auth password', required=True, env_var='ES_PASS')
    required.add_argument("--es_url", help='ElasticSearch url', env_var='ES_URL', required=True)
    required.add_argument("--es_index", help='Index name, default: openresty-*', env_var='ES_INDEX', required=True)
    optional.add_argument("--es_scheme", default='https', help='Transport, default: https',
                          choices=('https', 'http'), env_var='ES_SCHEME')
    optional.add_argument("--es_timeout", default=15, help='ES read timeout, default: 15', type=int, env_var='ES_SCHEME')
    optional.add_argument("--es_port", default=443, help='80 or 443, default: 443', type=int, env_var='ES_PORT')
    optional.add_argument("--ban_threshold", default=1000,
                          help='Count to get banned, integer, default: 1000', type=int, env_var='BAN_THRESHOLD')

    optional.add_argument("--time_window", default=10, help="Time window, in minutes",
                          choices=(10, 15, 30, 60), type=int, env_var='TIME_WINDOW')
    optional.add_argument("--retention", default=10, help='Flush bans after, in days, integer, default: 10',
                          type=int, env_var='RETENTION')
    optional.add_argument("--pdns_api_url", help='PowerDNS API URL: http://fqdn:8081/api/v1', env_var='PDNS_API_URL')
    optional.add_argument("--manual_remove", help='list of IPS to remove from rbl',
                          env_var='MANUAL_REMOVE', nargs='+')
    optional.add_argument("--list_banned", default=False, action='store_true', help='list of banned IPS from rbl',
                          env_var='LIST_BANNED')
    optional.add_argument("--dry_run", default=False, action='store_true', help='Just print, do not change',
                          env_var='DRY_RUN')
    required.add_argument("--pdns_api_key", help='PowerDNS API KEY', required=True, env_var='PDNS_API_KEY')
    required.add_argument("--pdns_rbl_zone", default='hostinger.rbl.', help='RBL zone, default: hostinger.rbl',
                          required=True, env_var='PDNS_RBL_ZONE')
    required.add_argument("--mysql_host", help='Mysql DB host/ip', required=True, env_var='MYSQL_HOST')
    required.add_argument("--mysql_pass", help='Mysql password', required=True, env_var='MYSQL_PASS')
    required.add_argument("--mysql_user", help='Mysql user', required=True, env_var='MYSQL_USER')
    required.add_argument("--mysql_db", help='Mysql database', required=True, env_var='MYSQL_DB')
    return parser.parse_args()


class Es:
    def __init__(self, args):
        self.hosts = args.es_url
        self.http_auth = (args.es_user, args.es_pass)
        self.port = args.es_port
        self.scheme = args.es_scheme
        self.rq_timeout = args.es_timeout
        self.index = args.es_index
        self.time_window = args.time_window
        self.__es_client = self.__get_client()

    def __get_client(self):
        es = None
        try:
            es = Elasticsearch(hosts=self.hosts,
                               http_auth=self.http_auth,
                               port=self.port,
                               scheme=self.scheme)
        except ConnectionError:
            print("Error getting abusers from elastic: {}".format(ConnectionError.info))
        return es

    def __get_timestamps(self):
        dt_now_mills = round(datetime.datetime.now().timestamp() * 1000)
        dt_minus = datetime.datetime.now() - datetime.timedelta(minutes=self.time_window)
        past_mills = round(dt_minus.timestamp() * 1000)
        return dt_now_mills, past_mills

    def construct_query(self, lookup="POST /xmlrpc.php.*"):
        now, before = self.__get_timestamps()
        search_obj = Search(index=self.index)

        query = Q('bool', must=[Q('query_string', query='request:\"{}\"'.format(lookup),
                                  analyze_wildcard=True,
                                  default_field="*")])
        search_obj.aggs.bucket('ips', 'terms', field='remote_addr.keyword', size=500, order={'_count': 'desc'}) \
            .bucket('hosts', 'cardinality', field='host.keyword')
        search_obj = search_obj.query('range', **{'@timestamp': {"gte": before, "lte": now}})
        search_obj = search_obj.query(query)

        return search_obj

    def execute_search(self, search_object):
        ex = search_object.using(self.__es_client)
        return ex.execute()


class PowerDnsClient:

    def __init__(self, args):
        self.pdns_api_key = args.pdns_api_key
        self.pdns_api_url = args.pdns_api_url
        self.pdns_zone_name = args.pdns_rbl_zone
        self.api = self.__get_powerdns_client()
        if args.dry_run:
            print("Dry run, will not create DNS RBL zone")
        if not self.api.servers[0].get_zone(self.pdns_zone_name) and not args.dry_run:
            self.zone = self.__create_pdns_zone()
        else:
            self.zone = self.api.servers[0].get_zone(self.pdns_zone_name)

    def __get_powerdns_client(self):
        api_client = powerdns.PDNSApiClient(api_endpoint=self.pdns_api_url, api_key=self.pdns_api_key)
        api = powerdns.PDNSEndpoint(api_client)
        return api

    def __create_pdns_zone(self):
        print('Creating PowerDNS zone {}'.format(self.pdns_zone_name))
        serial = datetime.date.today().strftime("%Y%m%d00")
        soa = "ns1.%s %s %s 28800 7200 604800 86400" % (self.pdns_zone_name, self.pdns_zone_name, serial)
        soa_r = powerdns.RRSet(name='hostinger.rbl.',
                               rtype="SOA",
                               records=[(soa, False)],
                               ttl=86400)
        zone = self.api.servers[0].create_zone(name="{}".format(self.pdns_zone_name),
                                               kind="Native",
                                               rrsets=[soa_r],
                                               nameservers=["ns1.hostinger.rbl.",
                                                            "ns2.hostinger.rbl."])
        return zone

    def create_pdns_rrsets(self, iplist):
        print('Create new ban records for: ', *iplist, sep='\n')
        for ip in iplist:
            self.zone.create_records([
                powerdns.RRSet(self.reverse(ip[0]), 'A', ['127.0.0.1'])
            ])

    def remove_pdns_rrsets(self, iplist):
        print('Remove ban records from PowerDNS if any {}'.format(iplist))
        for ip in iplist:
            self.zone.delete_record([
                powerdns.RRSet(self.reverse(ip[0]), 'A', ['127.0.0.1'])
            ])
            print('Sent removal query for {}'.format(self.reverse(ip[0]) + '.' + self.pdns_zone_name))

    def get_pdns_rrsets(self):
        return self.zone.records

    @staticmethod
    def reverse(ip):
        if len(ip) <= 1:
            return ip
        return '.'.join(ip.split('.')[::-1])


class MysqlDB:
    def __init__(self, args):
        self.host = args.mysql_host
        self.user = args.mysql_user
        self.password = args.mysql_pass
        self.database = args.mysql_db
        self.db = self.__get_db()

    def __get_db(self):
        datab = mysql.connector.connect(
            host=self.host,
            user=self.user,
            password=self.password,
            database=self.database
        )
        return datab

    def create_db_records(self, banlist):
        values = []
        cursor = self.db.cursor()
        for element in banlist:
            values.append((element[0], self.__timestamp_now(), element[1]))
        cursor.executemany("REPLACE INTO hostinger_dnsrbl (IP,CREATED_AT,COMMENT) VALUES (%s,%s,%s)", values)
        self.db.commit()

    def remove_manually(self, iplist):
        cursor = self.db.cursor()
        cursor.executemany("DELETE FROM hostinger_dnsrbl WHERE ip= %s", iplist)
        self.db.commit()

    def remove_expired_bans(self, days_to_keep):
        cursor = self.db.cursor()
        cursor.execute("SELECT IP FROM hostinger_dnsrbl WHERE CREATED_AT <= NOW() - INTERVAL %s DAY", (days_to_keep,))
        ips = list(cursor.fetchall())
        print('Cleanup from DB if any: {}'.format(ips))
        cursor.execute("DELETE FROM hostinger_dnsrbl WHERE CREATED_AT <= NOW() - INTERVAL %s DAY",
                       (days_to_keep,))
        self.db.commit()
        return ips

    def __timestamp_now(self):
        self.ts_now = datetime.datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        return self.ts_now


def main():
    args = parse_args()
    pdns = PowerDnsClient(args)
    db = MysqlDB(args)
    if args.manual_remove:
        if args.dry_run:
            print('Would remove {} from rbl'.format(args.manual_remove))
            exit(0)
        else:
            cleanup(pdns_obj=pdns, db_obj=db, manual_remove=args.manual_remove)
            exit(0)
    if args.list_banned:
        bans = pdns.get_pdns_rrsets()
        for rrset in bans:
            if rrset['type'] == 'A':
                print(rrset.get('name'))
    else:
        elastic = Es(args)
        '''
        Separated query object for future, to be able to generate different objects,
        if same structure is needed, it should be fine using same object just passing different "lookup=" variable.
        '''
        query_obj = elastic.construct_query(lookup='POST /xmlrpc.php.*')
        response = elastic.execute_search(query_obj)
        badlist = []

        for bucket in response.aggregations.ips.buckets:
            try:
                ipaddress.IPv4Address(bucket['key'])
            except ValueError:
                continue
            if (bucket['doc_count'] > args.ban_threshold
                    and ipaddress.IPv4Address(bucket['key']).is_global
                    and bucket.hosts['value'] > 10):
                reason = 'Above threshold xmlrpc hits: {0} ' \
                         'Unique webs: {1} ' \
                         'Time window: {2} minutes'.format(bucket['doc_count'], bucket.hosts['value'], args.time_window)
                badlist.append((bucket['key'], reason))

        if args.dry_run:
            print('would ban: ', *badlist, sep='\n')
        else:
            pdns.create_pdns_rrsets(badlist)
            db.create_db_records(badlist)
            cleanup(args.retention, pdns, db)


def cleanup(retention=None, pdns_obj=None, db_obj=None, manual_remove=None):
    args = parse_args()
    if not pdns_obj or not db_obj:
        pdns_obj = PowerDnsClient(args)
        db_obj = MysqlDB(args)
    if manual_remove:
        iplist = []
        [iplist.append((i,)) for i in manual_remove]
        db_obj.remove_manually(iplist)
        pdns_obj.remove_pdns_rrsets(iplist)
    else:
        ips = db_obj.remove_expired_bans(retention)
        pdns_obj.remove_pdns_rrsets(ips)


if __name__ == "__main__":
    main()
