# BscFees

## Intro

[![StackShare](http://img.shields.io/badge/tech-stack-0690fa.svg?style=flat)](https://stackshare.io/bscfees/bscfees) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkatsadim%2Fbscfees.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkatsadim%2Fbscfees?ref=badge_shield)

[![codecov](https://codecov.io/gh/katsadim/bscfees/branch/main/graph/badge.svg)](https://codecov.io/gh/katsadim/bscfees)
[![API](https://github.com/katsadim/bscfees/workflows/API/badge.svg)](https://github.com/katsadim/bscfees/workflows/API/badge.svg) 
[![WEB](https://github.com/katsadim/bscfees/workflows/WEB/badge.svg)](https://github.com/katsadim/bscfees/workflows/WEB/badge.svg) 
[![Terraform](https://github.com/katsadim/bscfees/workflows/TF/badge.svg)](https://github.com/katsadim/bscfees/workflows/TF/badge.svg) 

Are you tired of manually going over your trades in [BscScan](https://bscscan.com/) and [Etherscan](https://etherscan.com) 
to calculate your transaction fees?

[BscFees](https://bscfees.com) is your go-to place when it comes to Binance Safe Chain and Ethereum Blockchain fees calculation!

<p align="center">
  <img src="/res/site.webp">
</p>

## How

Just input the wallet address and BscFees will:

* fetch the last [bnb/busd](https://www.binance.com/en/trade/BNB_BUSD) currency rate
* iterate over all your transactions and sum your fees
* multiply the above two values

And there you have it! The sum of fees you have paid for Bsc and Eth transactions right in front of your pretty eyes! 

## Repo contents

This repository hosts the bits and pieces that make BscFees a reality:

* Backend
* Frontend
* [Infrastructure as code](tf/)

## Future work

* Show final fees price based on the bnb/busd currency rate at the exact time the transaction completed.
* Support other currencies (eg EUR)
* Support of more than 10000 transactions (bots)



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fkatsadim%2Fbscfees.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fkatsadim%2Fbscfees?ref=badge_large)