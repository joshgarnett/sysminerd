#!/bin/bash

wget -O - http://packages.elasticsearch.org/GPG-KEY-elasticsearch | apt-key add -

cat <<'EOT' > /etc/apt/sources.list.d/elasticsearch.list
deb http://packages.elasticsearch.org/elasticsearch/1.0/debian stable main
EOT

echo "Installing elasticsearch"
apt-get update -qq
apt-get install -qq -y openjdk-7-jre-headless elasticsearch

update-rc.d elasticsearch defaults 95 10

cp /etc/elasticsearch/elasticsearch.yml /etc/elasticsearch/elasticsearch.yml.bak

cat <<'EOT' > /etc/elasticsearch/elasticsearch.yml
cluster.name: elasticsearch
http.port: 9200
http.enabled: true
EOT

cat <<'EOT' > /etc/default/elasticsearch
ES_HEAP_SIZE=64m
EOT

/usr/share/elasticsearch/bin/plugin -install mobz/elasticsearch-head

echo "Starting elasticsearch"
/etc/init.d/elasticsearch start

echo "Waiting to create grafana-dash index"

sleep 15

echo "Creating grafana-dash index"

curl -XPUT 'http://localhost:9200/grafana-dash/' -d '
index :
    number_of_shards : 1
    number_of_replicas : 0
'

echo "Setting up Grafana"

curl http://grafanarel.s3.amazonaws.com/grafana-1.9.0.tar.gz > /tmp/grafana-1.9.0.tar.gz
cd /tmp
tar xzvf grafana-1.9.0.tar.gz
mv grafana-1.9.0 /usr/share/grafana
cd /usr/share/grafana

cat <<'EOT' > /usr/share/grafana/config.js
define(['settings'], function(Settings) {
  return new Settings({
      datasources: {
        graphite: {
          type: 'graphite',
          url: "http://192.168.33.50",
        },
        elasticsearch: {
          type: 'elasticsearch',
          url: "http://192.168.33.50:9200",
          index: 'grafana-dash',
          grafanaDB: true,
        }
      },
      search: {
        max_results: 100
      },
      default_route: '/dashboard/file/default.json',
      unsaved_changes_warning: true,
      playlist_timespan: "1m",
      admin: {
        password: ''
      },
      window_title_prefix: 'Grafana - ',
      // Add your own custom panels
      plugins: {
        // list of plugin panels
        panels: [],
        // requirejs modules in plugins folder that should be loaded
        // for example custom datasources
        dependencies: [],
      }
    });
});
EOT

chown -R www-data:www-data /usr/share/grafana

echo "alias /grafana /usr/share/grafana" > /etc/apache2/sites-enabled/grafana.conf

service apache2 restart
