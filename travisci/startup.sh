#!/bin/sh

TRAVISDIR=`dirname $0`
cd ${TRAVISDIR}
PWD=`pwd`

rabbitmq-plugins enable rabbitmq_management
rabbitmqadmin -q import ${PWD}/rabbit.config
cd -
