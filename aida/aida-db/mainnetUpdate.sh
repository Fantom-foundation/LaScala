. /etc/profile.d/golang.sh

DIR="aida-update"
BRANCH="develop"

if [ -d "$DIR" ]; then
  rm -rf "$DIR"
fi

# Clone the repository and build Aida
git clone https://github.com/Fantom-foundation/Aida "$DIR"
cd "$DIR"
git checkout $BRANCH
git submodule update --init --recursive
make clean
make

# Update AidaDB
./build/util-db update --aida-db=/var/opera/Aida/mainnet-data/aida-db --chainid 250 --db-tmp=/var/opera/Aida/dbtmpjenkins --log debug

# Cleanup
make clean
cd ../..
# rm -rf "$DIR"