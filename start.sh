#!/bin/sh
echo shut down existed docker service
echo you env is $1
if [ $1 == "TEST" ]
then
    export RUNTIME="test"
    docker stop megaoasis_profie_test
    docker stop megaoasis_email_test
    docker stop megaoasis_img_test

    docker container rm megaoasis_profie_test
    docker container rm megaoasis_email_test
    docker container rm megaoasis_img_test

    docker rmi test_megaoasis_profie -f
    docker rmi test_megaoasis_email -f
    docker rmi test_megaoasis_img -f

    docker-compose -p "test" up -d
fi

if [ $1 == "STAGING" ]
then
    export RUNTIME="main"
    docker stop megaoasis_profie_main
    docker stop megaoasis_email_main


    docker container rm megaoasis_profie_main
    docker container rm megaoasis_email_main

    docker rmi main_megaoasis_profie -f
    docker rmi main_megaoasis_email -f

    docker-compose -p "main" up -d
fi


