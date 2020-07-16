#!/bin/bash

openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout https-server.key -out https-server.crt -subj '/O=GoOut/C=US'
