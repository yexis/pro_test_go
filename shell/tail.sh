tail -f log/node.log.2025040212 | grep -a INFO | while read -r line; do
	query=$(echo $line | grep -oP " query:(\S*) dci" | awk -F"dci" '{print $1}')
	response=$(echo $line | grep -oP "response:(\S+)")
	echo $query" "$response
done