#Introduction
I had just (22 december 2013 01:58 UTC+00) forked https://github.com/oikomi/mygocrawl 
The commands below assume you have a working Go 1.2 compiler, and you have set the GOPATH enviroment variable.

# Go getting it

go get github.com/larsth/mygocrawl


#Getting the repository using Git

Using a bash shell, the commands are (should be):

mkdir -p $GOPATH/src/github.com/larsth
cd $GOPATH/src/github.com/larsth
git clone https://github.com/larsth/mygocrawl.git

# Install it

go install github.com/larsth/mygocrawl

# Important notes

I has not tested the golang code with a Go 1.2 compiler.
Also I had not done any tests, while the program is running (behaveiour tests, does it crash?, etc.).