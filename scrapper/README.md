# SCRAPPER

- The reason of this project is the necessity to create the rentability of apartments in home portals in a fast way and detect good deals.

## GOALS

- [x] Scrape the list with all the apartments and go through all the pages automatically.
- [x] Scrape the apartment information and store on a database.
- [ ] Scrape the average rent in the cities depending on the location.
- [ ] Use this information to generate an estimate of the rentability of the house.

## THOUGHT PROCESS

1. How to scrap the webs, from what I researched it looks like there are some mechanisms in form to prevent people from scrape these webs, so what I have done is look at the network calls in the browser and I found out that most of the information is sent in a request as html server-side rendered, so I just need the link to get the html and then parse it and get all the information that I need.

2. For the html parsing and extraction I chose cheerio library, for the call I chose axios and for the database sqlite (easy to store information).

3. There is some information like the images and some stats that doesn't come in the server-side html, I have tried using puppetteer for that but looks like the webs catch it very easily and put a captcha or just return an error in the first or second call, there are some ways around that but it involves using proxies or captcha resolvers and for the information that is left out at the moment it's not necessary to go this far.

4. For the scraping of the list and the apartments has been very straight forward, just needed to call the endpoints, looks like the language header is needed to be able to call it, but that was an easy fix and now it's working fine. Everything else is just look for the classes and ids in the html and get them.
