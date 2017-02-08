pushd ..
go build -v
popd
cp ../ham-lifx-bridge .
docker build -t immesys/eopdemo .
docker push immesys/eopdemo
