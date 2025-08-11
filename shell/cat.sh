date=$1
cuid=$2
cat output/log/agent.log.$date | grep -a $cuid | grep -a INFO | while read -r line; do
	logid=$(echo $line | grep -oP " logid:(\S*)")
	hit_logid=$(echo $line | grep -oP " hit_logid:(\S*)")
	rtc_sess_id=$(echo $line | grep -oP " rtc_sess_id:(\S*)")
	query=$(echo $line | grep -oP " query:(\S*) asr_last" | awk -F"asr_last" '{print $1}')

	response=$(echo $line | grep -oP "response:(\S+) push" | awk -F "push" '{print $1}')
	echo $response

	push_response=$(echo $line | grep -oP "push_response:(\S+)")
	echo $push_response

	tts_text=$(echo $line | grep -oP "tts_text:(\S+)")
	hit_logid=$(echo $line | grep -oP "hit_logid:(\S+)")
	nlu=`grep -a $hit_logid output/log/agent.log.$date | grep -a "get nlu success"`
	echo $logid
	echo $hit_logid
	echo $rtc_sess_id
	echo $query
	echo $tts_text
	echo $nlu
	echo
done