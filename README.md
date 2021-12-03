# TerranSeed

## TerranSeed is a fork of [TinySeed](https://github.com/notional-labs/tinyseed), which is a fork of Binary Holding's Tenderseed, which is a fork of Polychain's Tenderseed

This tool runs a seed node for any tendermint based blockchain (i.e. Terra, Cosmos, Akash, etc, thanks to TinySeed), crawls the network and generates a map with the geolocation of the peers.
It can be used to pinpoint centralization of networks in common infrastructure hosts like AWS, GCP etc.

###Configuration

```bash
git clone https://github.com/Terran-Stakers/terranseed
go mod tidy
npm install
npm run build
go install .
./tenderseed
```

Then you'll become a seed node on Columbus-5. 
A `$HONE/.terranseed/config/config.toml` file will be generated if it doesn't exist yet, with some default parameters.
You need to fill the `seeds` and `chain_id` and start the process again.
It may take few minutes before discovering peers, depending on the network.

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
