# docker for photon on spectrum testnet with new random account
# you can use this dockerfile like this:

#   cd Photon/docker
#   docker build -t {{YOUR_IMAGE_NAME}} .
#   docker run -it -p 0.0.0.0:5001:5001 -p 0.0.0.0:40001:40001 {{YOUR_IMAGE_NAME}}

# you can test with command :
#   curl http://127.0.0.1:5001/api/1/address
# if photon works well,this will return the address of account
FROM golang:1.12

RUN mkdir -p /go/src/github.com/SmartMeshFoundation
WORKDIR /go/src/github.com/SmartMeshFoundation
# get spectrum and photon
RUN git clone -b master https://github.com/SmartMeshFoundation/Spectrum.git
RUN git clone -b master https://github.com/SmartMeshFoundation/Photon.git

# build smc and photon
WORKDIR /go/src/github.com/SmartMeshFoundation/Spectrum/cmd/smc
RUN go install
WORKDIR /go/src/github.com/SmartMeshFoundation/Photon/cmd/photon
RUN /bin/bash ./build.sh
RUN cp photon /go/bin/

# create new account with random passwd
RUN mkdir -p /photon/keystore
WORKDIR /photon/keystore
RUN echo -n `date +%s%N | md5sum |cut -c 1-9` > pass
RUN address=$(smc account new --keystore=/photon/keystore --password=/photon/keystore/pass | cut -c 11-50)

# run photon on spectrum testnet
WORKDIR /photon
CMD photon --datadir=/photon/dbdata --api-address=0.0.0.0:5001 --listen-address=0.0.0.0:40001 --address=$address --keystore-path=/photon/keystore  --password-file=/photon/keystore/pass  --eth-rpc-endpoint=http://transport01.smartmesh.cn:17888 --verbosity=5 --logfile=/photon/photon.log

# CMD ["photon"," --datadir=/photon/dbdata","--api-address=0.0.0.0:5001","--listen-address=0.0.0.0:40001","--address=$address","--keystore-path=/photon/keystore","  --password-file=/photon/keystore/pass","--eth-rpc-endpoint=http://transport01.smartmesh.cn:17888","--verbosity=5","--logfile=/photon/photon.log"]
