#!/bin/bash

# see https://www.digitalocean.com/community/tutorials/how-to-install-and-use-graphite-on-an-ubuntu-14-04-server

echo "Installing apt packages"
apt-get update -qq
DEBIAN_FRONTEND=noninteractive apt-get install -qq -y --force-yes -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold" graphite-web graphite-carbon

# edit /etc/default/graphite-carbon
cat <<'EOT' > /etc/default/graphite-carbon
CARBON_CACHE_ENABLED=true
EOT

# edit /etc/carbon/storage-schemas.conf
cat <<'EOT' > /etc/carbon/storage-schemas.conf
[carbon]
pattern = ^carbon\.
retentions = 60:90d

[default]
pattern = .*
retentions = 5s:1d,30s:7d
EOT

# edit /etc/carbon/storage-aggregation.conf
cat <<'EOT' > /etc/carbon/storage-aggregation.conf
[min]
pattern = \.min$
xFilesFactor = 0.1
aggregationMethod = min

[max]
pattern = \.max$
xFilesFactor = 0.1
aggregationMethod = max

[sum]
pattern = \.count$
xFilesFactor = 0
aggregationMethod = sum

[default_average]
pattern = .*
xFilesFactor = 0.2
aggregationMethod = average
EOT

echo "Setting up the database"
# should replace this with postgres
graphite-manage syncdb

chown _graphite:_graphite /var/lib/graphite/graphite.db
chmod a+w /var/lib/graphite/graphite.db

echo "Starting carbon service"
service carbon-cache start

# install apache
echo "Installing apache"

apt-get install -qq -y --force-yes apache2 libapache2-mod-wsgi

a2dissite 000-default

cat <<'EOT' > /etc/apache2/sites-available/apache2-graphite.conf
<VirtualHost *:80>
	WSGIDaemonProcess _graphite processes=2 threads=2 display-name='%{GROUP}' inactivity-timeout=120 user=_graphite group=_graphite
	WSGIProcessGroup _graphite
	WSGIImportScript /usr/share/graphite-web/graphite.wsgi process-group=_graphite application-group=%{GLOBAL}
	WSGIScriptAlias / /usr/share/graphite-web/graphite.wsgi

	Alias /content/ /usr/share/graphite-web/static/
	<Location "/content/">
		SetHandler None
	</Location>

	ErrorLog ${APACHE_LOG_DIR}/graphite-web_error.log

	# Possible values include: debug, info, notice, warn, error, crit,
	# alert, emerg.
	LogLevel warn

	CustomLog ${APACHE_LOG_DIR}/graphite-web_access.log combined
</VirtualHost>
EOT

a2ensite apache2-graphite

service apache2 restart

echo "Graphite setup"
