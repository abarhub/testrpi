#!/bin/bash

# !/usr/bin/env bash

# to run in background :
# nohup ./run.sh >>custom-output.log 2>&1 &
# to stop
# sudo netstat -anp | grep 5000
# kill -9 XXXX

aide(){
  echo "Parametres possibles : start / stop / status"
}

if [[ -z $1 ]];
then
  echo "No parameter passed."
  aide
elif [[ "$1" = "start" ]];
then
  echo "start ..."
  nohup ./go_build_test5_go_arm_linux >>custom-output.log 2>&1 &
elif [[ "$1" = "stop" ]];
then
  echo "stop ..."
  PID=`sudo netstat -anp | grep 3000 | tr -s ' ' | cut -d ' ' -f 7 | cut -d '/' -f 1`
  echo "PID=${PID}"
  if [[ -z ${PID} ]];
  then
    echo "Déjà arrété"
  else
    kill -SIGINT ${PID}
    echo "Stoppé"
  fi
elif [[ "$1" = "status" ]];
then
  PID=`sudo netstat -anp | grep 3000 | tr -s ' ' | cut -d ' ' -f 7 | cut -d '/' -f 1`
  echo "PID=${PID}"
  if [[ -z ${PID} ]];
  then
    echo "Arrété"
  else
    echo "Démarré (pid=${PID})"
  fi
else
    echo "Parameter passed = $1"
    aide
fi

#nohup ./run.sh >>custom-output.log 2>&1 &