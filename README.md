# fwew-api

Fwew REST API

## install

pre-requisites:

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/)

```bash
git clone https://github.com/fwew/fwew-api.git
cd fwew-api
go build ./...
```

## configure

edit the included `config.json` file with your desired port and web root link.

config.json (defaults):

```json
{
  "Port": "80",
  "WebRoot": "http://localhost"
}
```

## deploy

simply run the binary resulting from the install step above.

```bash
./fwew-api
```

## endpoints

here is a quick run-down on the endpoints. in the following,`DOMAIN` is your domain name. (`example.com` or similar)

### search Na'vi to local

`http://DOMAIN/fwew/{nav}`

`{nav}` can be any Na'vi word, plain or affixed.

Returns a list of entries.

### search local to Na'vi

`http://DOMAIN/fwew/r/{lang}/{local}`

`{lang}` is any one of the following 2-character language codes:

- `de` (German)
- `en` (English)
- `et` (Estonian)
- `fr` (French)
- `hu` (Hungarian)
- `nl` (Dutch)
- `pl` (Polish)
- `ru` (Russian)
- `sv` (Swedish)

`{local}` is any word in the given language, hopefully one that will be found in the dictionary.

Returns a list of entries.

### list all words

`http://DOMAIN/list`

Returns the entire list of entries.

### list words with given properties

`http://DOMAIN/list/{args}`

`{args}` is any "what cond spec" string. spaces and all, as usual.
see [fwew-lib](https://github.com/fwew/fwew-lib#list) for more info about list syntax.

Returns a list of entries.

### random words

`http://DOMAIN/random/{n}`

`{n}` is any whole number > 0 of random words to retrieve.

Returns a list of entries.

### random words with given properties

`http://DOMAIN/random/{n}/{args}`

`{n}` is any whole number > 0 of random words to retrieve.
`{args}` is any "what cond spec" string. spaces and all, same as with /list/.
see [fwew-lib](https://github.com/fwew/fwew-lib#list) for more info about list syntax.

Returns a list of entries.

### number to Na'vi

`http://DOMAIN/number/r/{num}`

`{num}` is any integer between 0 and 32767 inclusive.

Returns a number entry.

### Na'vi to number

`http://DOMAIN/number/{word}`

`{word}` is any Na'vi number spelled out as a word. For example, `mevolaw`.

Returns a number entry.

### lenition table

`http://DOMAIN/lenition`

Returns the lenition entry.

### version

`http://DOMAIN/version`

Returns the version information.
