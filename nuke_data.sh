#!/bin/bash
# this script should nuke downloaded clips as well as the .db files that manage them
rm tmp/raw/*
rm tmp/processed/*
rm *.sqlite
rm *.db
#rm compile/tmp/*
#rm content.txt
rm *.mp4
