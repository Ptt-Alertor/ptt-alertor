#!/bin/bash
export $(cat .env | grep -v ^# | xargs)
ptt-alertor