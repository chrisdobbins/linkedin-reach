#!/bin/sh

user=$(id -u)

if [ $user -ne 0 ]; then 
        echo "must be root to run this" 
        exit 1
fi

echo "starting install..."

echo "updating OS"
apt-get update -y
if [ $? -ne 0 ]; then
       echo "failed to update OS"
       exit 1
fi

echo "installing game"
wget https://github.com/chrisdobbins/linkedin-reach/releases/download/v1.1.2/game-linux-arm
if [ $? -ne 0 ]; then
        echo "game installation failed"
        exit 1
fi
mv game-linux-arm /usr/local/bin
chmod +x /usr/local/bin/game-linux-arm
echo "game installed successfully! to start, type game-linux-arm"


echo "downloading configuration for 3.5\" display"
hash git 2> /dev/null
if [ $? -ne 0 ]; then
       echo "installing git"
       apt-get install -y git
       if [ $? -ne 0 ]; then
               echo "failed to install git"
               exit 1
       fi
fi

git clone https://github.com/goodtft/LCD-show
cd LCD-show
# modify install script to continue instead of waiting for user approval
cat LCD35-show | sed s/"apt-get"/"apt-get -y"/g > LCD35-show-mod
mv LCD35-show-mod LCD35-show
chmod +x LCD35-show
./LCD35-show
if [ $? -ne 0 ]; then
        echo "failed to install config for display"
        exit 1
fi

