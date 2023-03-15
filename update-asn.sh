#!/usr/bin/env zsh

/usr/local/bin/asnlist -country cn > asn.cn
/usr/local/bin/asnlist -country hk > asn.hk
/usr/local/bin/asnlist -country us > asn.us

rm -f ./asn.csv
for item in Google Microsoft Disney Netflix Amazon Cloudflare Apple; do
    asns=$(grep $item asn.us | awk '{print $2}' | tr ',' '|' | tr -d '\n')
    asns=$(echo $asns | sed 's/.$//')  # removes the last character (|)
    echo "$item,$asns" >> asn.csv
done

asns=$(cat asn.cn | awk '{print $2}' | tr ',' '|' | tr -d '\n')
asns=$(echo $asns | sed 's/.$//')  # removes the last character (|)
echo "cn,$asns" >> asn.csv

asns=$(cat asn.hk | awk '{print $2}' | tr ',' '|' | tr -d '\n')
asns=$(echo $asns | sed 's/.$//')  # removes the last character (|)
echo "hk,$asns" >> asn.csv
