#!/usr/bin/env bash

sessionID=`date +"%s"`

generate_post_data()
{
cat <<EOF
{
  "breakTts":"1",
  "callNumber":"10086",
  "calltype":"1",
  "mobile":"18565422776",
  "userId": "$sessionID"
}
EOF
}

curl -i \
-H "Accept: application/json" \
-H "Content-Type:application/json" \
-X POST --data "$(generate_post_data)" "http://127.0.0.1:9081/_set_/demo"

curl -XPOST  http://127.0.0.1:9081/_mock_/demo
exit 0;
