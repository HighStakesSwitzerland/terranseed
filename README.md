# TerranSeed

## TerranSeed is a fork of [TinySeed](https://github.com/notional-labs/tinyseed), which is a fork of Binary Holding's Tenderseed, which is a fork of Polychain's Tenderseed

This tool runs a seed node for multiple tendermint based blockchain (i.e. Terra, Cosmos, Akash, etc, thanks to TinySeed), crawls the network and expose the list of peers with the geolocation.

###Configuration

```bash
git clone https://github.com/Terran-Stakers/terranseed
go mod tidy
npm install
npm run build
go install .
./terranseed
```

A file `$HOME/.terranseed/config/config.toml` will be generated if it doesn't exist yet, with some default parameters, and the program will exit.

You need to fill the `seeds` and `chain_id` for every chain and start it again.
It may take few minutes/hours before discovering peers, depending on the network. 

## License

[Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0)
