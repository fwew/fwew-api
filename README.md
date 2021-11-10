# fwew-api

Fwew REST API

Currently deployed at [tirea.learnnavi.org/api](https://tirea.learnnavi.org/api)

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
  "Port": "10000",
  "WebRoot": "http://localhost"
}
```

## deploy

simply run the binary resulting from the install step above.

```bash
./fwew-api
```

## endpoints

here is a quick run-down on the endpoints.

### list the endpoints and syntax

`/`

The root endpoint returns an object containing the endpoints with expected parameters as values.

### search Na'vi to local

`/fwew/{nav}`

`{nav}` can be any Na'vi word, plain or affixed.

Returns an array of Word objects.

### search local to Na'vi

`/fwew/r/{lang}/{local}`

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

Returns an array of Word objects.

### list all words

`/list`

Returns an array containing every Word object in the dictionary.

### list words with given properties

`/list/{args}`

`{args}` is any "what cond spec" string. spaces and all, as usual.
see [fwew-lib](https://github.com/fwew/fwew-lib#list) for more info about list syntax.

Returns an array of Word objects.

### random words

`/random/{n}`

`{n}` is any whole number > 0 of random words to retrieve.

Returns an array of Word objects.

### random words with given properties

`/random/{n}/{args}`

`{n}` is any whole number > 0 of random words to retrieve.
`{args}` is any "what cond spec" string. spaces and all, same as with /list/.
see [fwew-lib](https://github.com/fwew/fwew-lib#list) for more info about list syntax.

Returns an array of Word objects.

### number to Na'vi

`/number/r/{num}`

`{num}` is any integer between 0 and 32767 inclusive.

Returns a number entry object.

### Na'vi to number

`/number/{word}`

`{word}` is any Na'vi number spelled out as a word. For example, `mevolaw`.

Returns a number entry object.

### lenition table

`/lenition`

Returns the lenition object.

### version

`/version`

Returns the version information object.
