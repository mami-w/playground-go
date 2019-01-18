#!/bin/sh
echo $timetracker_endpoint
echo $timetracker_username
#echo $timetracker_pwd

./timetracker $timetracker_endpoint $timetracker_username $timetracker_pwd
