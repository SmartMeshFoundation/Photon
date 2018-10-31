# Photon
![](http://img.shields.io/travis/dognie/Photon.svg)
![](https://github.com/dognie/Photon/blob/master/docs/images/photon1.png?raw=true)

  [Photon](https://PhotonNetwork.readthedocs.io/en/latest/) is an off-chain scaling solution, enabling instant, low-fee and scalable payments. It’s complementary to the Ethereum blockchain and Spectrum blockchain and works with ERC20 compatible token and ERC223 compatible token. Photon currently can works on Windows, Linux ,Android, iOS etc. Photon adds some new functions in version 0.9, such as , cooperative settlement,widraw without closing the channel, and more perfect third-party services( [Photon Monitoring](https://github.com/SmartMeshFoundation/Photon-Monitoring) and  [ Photon-Path-Finder](https://github.com/SmartMeshFoundation/Photon-Path-Finder)). In order to improve the user experience and to better fit for the mobile network, Photon adopts the  Matrix communication mechanism and supports the crash recovery and channel charging function.
## Project Status
  This project is still very much a work in progress. It can be used for testing, but it should not be used for real funds. We are doing our best to identify and fix problems, and implement missing features. Any help testing the implementation, reporting bugs, or helping with outstanding issues is very welcome.

## Build
```
  go get github.com/SmartMeshFoundation/Photon/
  cd $GOPATH/github.com/SmartMeshFoundation/Photon
  make 
  ./build/bin/photon
```

## mobile support
Photon can works on Android and iOS using mobile's API.  it needs [go mobile](https://github.com/golang/mobile) to build mobile library.
### build Android mobile library
```bash
cd mobile
./build_Android.sh 
```
then you can integrate `mobile.aar` into your project.
### build iOS mobile framework
```bash
./build_iOS.sh
```
then you can integrate `Mobile.framework` into your project.
## Requirements
Latest version of SMC

We need go's plugin module.
