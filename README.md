# TerranSeed

## TerranSeed is a fork of [TinySeed](https://github.com/notional-labs/tinyseed), which is a fork of Binary Holding's Tenderseed, which is a fork of Polychain's Tenderseed

This tool runs a seed node for any tendermint based blockchain (i.e. Terra, Cosmos, Akash, etc, thanks to TinySeed), crawls the network and generates a map with the geolocalisation of the peers.
It is used to pinpoint centralization of network in common infrastructure hosts like AWS, GCP etc.

###Configuration

```bash
git clone https://github.com/Terran-Stakers/terranseed
go mod tidy
go install .
tenderseed
```

Then you'll become a seed node on Osmosis-1. Let's do Cosmoshub-4, shall we? We've made Osmosis zeroconf, but hey this
thing here reads 2 env vars!

```bash
export ID=cosmoshub-4
export SEEDS=bf8328b66dceb4987e5cd94430af66045e59899f@public-seed.cosmos.vitwit.com:26656,cfd785a4224c7940e9a10f6c1ab24c343e923bec@164.68.107.188:26656,d72b3011ed46d783e369fdf8ae2055b99a1e5074@173.249.50.25:26656,ba3bacc714817218562f743178228f23678b2873@public-seed-node.cosmoshub.certus.one:26656,3c7cad4154967a294b3ba1cc752e40e8779640ad@84.201.128.115:26656,366ac852255c3ac8de17e11ae9ec814b8c68bddb@51.15.94.196:26656
terranseed
```

## License

[Blue Oak Model License 1.0.0](https://blueoakcouncil.org/license/1.0.0)

# Angular
This project was generated with [Angular CLI](https://github.com/angular/angular-cli) version 13.0.3.

## Development server

Run `ng serve` for a dev server. Navigate to `http://localhost:4200/`. The app will automatically reload if you change any of the source files.

## Code scaffolding

Run `ng generate component component-name` to generate a new component. You can also use `ng generate directive|pipe|service|class|guard|interface|enum|module`.

## Build

Run `ng build` to build the project. The build artifacts will be stored in the `dist/` directory.

## Running unit tests

Run `ng test` to execute the unit tests via [Karma](https://karma-runner.github.io).

## Running end-to-end tests

Run `ng e2e` to execute the end-to-end tests via a platform of your choice. To use this command, you need to first add a package that implements end-to-end testing capabilities.

## Further help

To get more help on the Angular CLI use `ng help` or go check out the [Angular CLI Overview and Command Reference](https://angular.io/cli) page.
