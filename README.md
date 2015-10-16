# ipgrep

A simple tool for grepping for IPs (or networks) in a list of IP addresses.

## Usage

Simply as follows

    [user@server] $ echo "1.2.3.4" >> sources
    [user@server] $ echo "2.3.4.5" >> sources
    [user@server] $ echo "2001:db0::face" >> sources
    [user@server] $ 
    [user@server] $ ipgrep 2.3.4.0/24 < sources
    2.3.4.5
    [user@server] $ echo "2000::/3" > spec
    [user@server] $ 
    [user@server] $ ipgrep -patterns ./spec < sources
    2001:db0::face
    
## Missing Features

* No support for `-v` to invert matches.
* Others

## Blame

By James Hannah <james@tlyk.eu>
