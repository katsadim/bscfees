# BscFees

## Intro

[![StackShare](http://img.shields.io/badge/tech-stack-0690fa.svg?style=flat)](https://stackshare.io/bscfees/bscfees) 
[![API](https://github.com/katsadim/bscfees/workflows/API/badge.svg)](https://github.com/katsadim/bscfees/workflows/API/badge.svg) 
[![WEB](https://github.com/katsadim/bscfees/workflows/WEB/badge.svg)](https://github.com/katsadim/bscfees/workflows/WEB/badge.svg) 
[![Terraform](https://github.com/katsadim/bscfees/workflows/TF/badge.svg)](https://github.com/katsadim/bscfees/workflows/TF/badge.svg) 

Are you tired of manually going over your trades in [BscScan](https://bscscan.com/) to calculate your transaction fees?

[BscFees](https://bscfees.com) is your go-to place when it comes to Binance Safe Chain fees calculation!

<p align="center">
  <img src="/res/site.webp">
</p>

## How

Just input the wallet address and BscFees will:

* fetch the last [bnb/busd](https://www.binance.com/en/trade/BNB_BUSD) currency rate
* iterate over all your transactions and sum your fees
* multiply the above two values

And there you have it! The sum of fees you have paid for Bsc transactions right in front of your pretty eyes! 

## Repo contents

This repository hosts the bits and pieces that make BscFees a reality:

* Backend
* Frontend
* [Infrastructure as code](tf/)

## Future work

* Support of ETH
* Show final fees price based on the bnb/busd currency rate at the exact time the transaction completed.
* Support other currencies (eg EUR)

