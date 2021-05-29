# Genshin Impact helper written in GO

Based on [takagg/genshin-impact-helper](https://github.com/takagg/genshin-impact-helper)

# How to use
My primary goal was to create a library, but you can run the program to automatically sign-in on the website.

Fill `config.toml` with your cookie and then just run the program with
flag `go run . -config config.toml`.

You can look at `example_config.toml` for references.

### To get a cookie string
- navigate to the website
- press F12
- select console window
- paste `document.cookie` and press enter
- copy output without quotation marks (" ")

# Changelog
## v0.1.0
* Support for multiple accounts
* Using toml config file