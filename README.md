# ctrace

The command line tool to check Cloud Trace existence.

## Usage

```
% ./ctrace
Usage: ./ctrace [-quiet] [-url] projects/PROJECT_ID/traces/xxxxxxxx
       ./ctrace [-quiet] [-url] PROJECT_ID xxxxxxxx
```

上の例の引数は HTTP(S) Load Balancer のログなどに入っている形式

例えば次のようにして Load Balancer のログを取り出して、
それを [jq](https://stedolan.github.io/jq/) とこのコマンドで処理する

```
#!/bin/bash

date=${1:-$(date -d yesterday +%Y-%m-%d)}

begin=$(date -d "$date" --utc +%FT%TZ)
end=$(date -d "$date 1 day" --utc +%FT%TZ)

if [ ! -f "${date}.json" ] ; then
  gcloud logging read \
    "resource.type=\"http_load_balancer\" \
     httpRequest.status=\"502\" \
     timestamp>=\"$begin\" timestamp<\"$end\"" \
    --format=json > $date.json
fi

while read line; do
  timestamp=$(echo "$line" | jq -r .timestamp)
  detail=$(echo "$line" | jq -r .jsonPayload.statusDetails)
  trace=$(echo "$line" | jq -r .trace)
  server=$(echo "$line" | jq -r .httpRequest.serverIp)
  latency=$(echo "$line" | jq -r .httpRequest.latency)
  url=$(echo "$line" | jq -r .httpRequest.requestUrl)
  trace=$(ctrace -quiet -url $trace)
  echo "$timestamp,$detail,$latency,$server,$url,$trace"
done < <(cat ${date}.json | jq -c '.[]')
```
