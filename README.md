# Genshin Impact helper written in GO

Based on [takagg/genshin-impact-helper](https://github.com/takagg/genshin-impact-helper)

# How to use
My primary goal was to create a library, but you can run the program to automatically sign-in on the website.

- paste cookie string (without "quotation marks") to file `cookie.txt` 
- run (`go run .`), it will sign-in you if you haven't already done this today

### To get a cookie string
- navigate to the website
- press F12
- select console window
- paste `document.cookie` and press enter
- copy output without quotation marks (" ")
