# geo

🗺An easy way to manage all your Geo resources.

Support GeoIP code lookup and converting among all popular GeoIP databases.

## Install

Requirements:

- [Go](https://go.dev) 1.20+

```shell
go install -v github.com/metacubex/geo/cmd/geo
```

## Usage

Type `geo help` for more details.

### Look up codes from existing Geo databases

The default directory to find Geo databases is
`~/.geo`. You can specify it through `-D` argument.

#### IP

```shell
geo look 223.5.5.5
```

```shell
geo look 2001:250::
```

Supported databases for `look` command:

- MaxMind MMDB
- V2Ray dat
- sing-geoip MMDB
- Meta-geoip MMDB

### Convert

#### IP

```shell
geo convert ip -i <input_type> -o <output_type> -f [output_filename] input_filename
```

```shell
geo convert ip -i v2ray -o meta ./geoip.dat
```

Available types:

- maxmind (MaxMind MMDB)
- v2ray (V2Ray dat)
- sing (sing-geoip MMDB)
- meta (Meta-geoip MMDB)

Supported conversion pairs: (Column=From, Row=To)

|            | MaxMind | V2Ray | sing-geoip | Meta-geoip |
|------------|---------|-------|------------|------------|
| MaxMind    | -       |       |            |            |
| V2Ray      |         | -     |            |            |
| sing-geoip | √       | √     | -          | √          |
| Meta-geoip |         | √     |            | -          |

Conversion to MaxMind is not available for legal reasons.  
Conversion to V2Ray is on the TODO list.

## F&Q

### Why conversion MaxMind/sing-geoip -> Meta-geoip is not available?

Meta-geoip is designed to support IP with multiple results,
which will help users who use GeoIP functionality as IPList or IPSet.

For sources such as MaxMind and sing-geoip, which only have a single possible result,
according to the principle of Occam's razor,
there is no need to convert to Meta-geoip database.

Clash.Meta supports all of these databases, so everything is well. :-)
