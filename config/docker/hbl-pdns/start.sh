#!/bin/bash

mkdir -p /etc/powerdns/pdns.d

PDNSVARS=`echo ${!PDNSCONF_*}`
touch /etc/powerdns/pdns.conf

for var in $PDNSVARS; do
  varname=`echo ${var#"PDNSCONF_"} | awk '{print tolower($0)}' | sed 's/_/-/g'`
  value=`echo ${!var} | sed 's/^$\(.*\)/\1/'`
  echo "$varname=$value" >> /etc/powerdns/pdns.conf
done

if [ ! -z $PDNSCONF_API_KEY ]; then
  cat >/etc/powerdns/pdns.d/api.conf <<EOF
api=yes
webserver=yes
webserver-address=0.0.0.0
webserver-allow-from=0.0.0.0/0
EOF

fi

mysqlcheck() {
  # Wait for MySQL to be available...
  COUNTER=20
  until mysql -h mysql -u $PDNSCONF_GMYSQL_USER -p$PDNSCONF_GMYSQL_PASSWORD -e "SHOW DATABASES" 2>/dev/null; do
    echo "WARNING: MySQL still not up. Trying again... $PDNSCONF_GMYSQL_DBNAME"
    sleep 10
    let COUNTER-=1
    if [ $COUNTER -lt 1 ]; then
      echo "ERROR: MySQL connection timed out. Aborting."
      exit 1
    fi
  done

  mysql -h mysql -u $PDNSCONF_GMYSQL_USER -p$PDNSCONF_GMYSQL_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $PDNSCONF_GMYSQL_DBNAME"

  count=`mysql -h mysql -u $PDNSCONF_GMYSQL_USER -p$PDNSCONF_GMYSQL_PASSWORD -e "SELECT COUNT(*) FROM information_schema.tables WHERE table_type='BASE TABLE' AND table_schema='$PDNSCONF_GMYSQL_DBNAME';" | tail -1`
  if [ "$count" == "0" ]; then
    echo "Database is empty. Importing PowerDNS schema..."
    mysql -h mysql -u $PDNSCONF_GMYSQL_USER -p$PDNSCONF_GMYSQL_PASSWORD $PDNSCONF_GMYSQL_DBNAME < /usr/share/doc/pdns-backend-mysql/schema.mysql.sql && echo "Import done."
  fi
}

mysqlcheck


# Start PowerDNS
# same as /etc/init.d/pdns monitor
echo "Starting PowerDNS..."

if [ "$#" -gt 0 ]; then
  exec /usr/sbin/pdns_server "$@"
else
  exec /usr/sbin/pdns_server --daemon=no --guardian=no --control-console --loglevel=9
fi
