#!/bin/bash

TESTS_RESULT=0

echo "======================= Prepare environment ======================="
docker-compose --file functional-tests/docker-compose.yml down
docker-compose --file functional-tests/docker-compose.yml up -d postgres
#wait for postgres start. Can be improved
sleep 5;
docker-compose --file functional-tests/docker-compose.yml up -d rest

echo -e "\n======================= Run tests ======================="


docker-compose --file functional-tests/docker-compose.yml up --abort-on-container-exit functional-tests

TESTS_RESULT=$?

echo -e "\n======================= Destroy environment ======================="
docker-compose --file functional-tests/docker-compose.yml down

if [ "$TESTS_RESULT" -eq "0" ];
then
    tput setaf 4; echo -e "\n!!!!!!!!!!!!!!!!!!!!!!! Tests successfully completed !!!!!!!!!!!!!!!!!!!!!!!"; tput sgr0;
else
    tput setaf 1; echo -e "\n!!!!!!!!!!!!!!!!!!!!!!! Tests failed !!!!!!!!!!!!!!!!!!!!!!!"; tput sgr0;
fi

exit ${TESTS_RESULT}