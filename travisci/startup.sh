#!/bin/sh

TRAVISDIR=`dirname $0`
cd ${TRAVISDIR}
ABSP=${PWD}
cd -

sudo rabbitmq-plugins enable rabbitmq_management
wget http://guest:guest@localhost:15672/cli/rabbitmqadmin
chmod a+x rabbitmqadmin
sudo ./rabbitmqadmin -v import ${PWD}/rabbit.config
