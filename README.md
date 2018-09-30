# SmartRaiden
  [SmartRaiden documentation](https://smartraiden.readthedocs.io/en/latest/)

  Smartraiden is an off-chain scaling solution, enabling instant, low-fee and scalable payments. It’s complementary to the Ethereum blockchain and Spectrum blockchain and works with ERC20 compatible token and ERC223 compatible token. SmartRaiden currently can works on Windows, Linux ,Android, iOS etc. The new version of smartraiden adds some new functions, such as , cooperative settlement,widraw without closing the channel, and more perfect third-party services. In order to better fit for the mobile network, Smartraiden adopts the XMPP communication mechanism and supports the crash recovery and channel charging function.
## Project Status
  This project is still very much a work in progress. It can be used for testing, but it should not be used for real funds. We are doing our best to identify and fix problems, and implement missing features. Any help testing the implementation, reporting bugs, or helping with outstanding issues is very welcome.

## Build
```
  go get github.com/SmartMeshFoundation/SmartRaiden/
  cd cmd/smartraiden
  go install
```
## Requirements
geth >=1.7.3
