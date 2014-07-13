#!/bin/sh

TRAVISDIR=`dirname $0`
cd ${TRAVISDIR}
ABSP=${PWD}
cd -

sudo rabbitmq-plugins enable rabbitmq_management
sudo rabbitmqadmin -v import ${PWD}/rabbit.config
