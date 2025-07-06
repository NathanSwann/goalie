# GOALIE - Goal Extensions

Some extensions for [goal](https://codeberg.org/anaseto/goal)

## DB

Allows reading and writting binary splayed files

- `s db.save t` write table `t` into directory `s`
- `db.get[s]` read table `s`

```
/ Save the table
"data/test"db.save![!"a b c";(1.0 2 3;4.0 5 6;!"one two three")]
/ makes a new folder:
/ data
/ ├── test
/ │   ├── a
/ │   ├── b
/ │   ├── c
/ │   └── c#
/ └── test.schema

/ Load table again:
db.get["data/test"]
/ !["a" "b" "c"(1.0 2.0 3.0;4.0 5.0 6.0;"one" "two" "three")]
```

## Dates

For handling dates

- `date.fs` turn a date string into a date ordinal (number of days from 2000-01-01)
- `date.sf` get the string representation of a date ordinal


```
/ convert string to a date ordinal
date.fs"2024-03-02"
/ 8827

/ works on lists
date.fs"2024-03-02","2024-01-03"
/ 8827 8768

/ Can covert back aswell
date.sf 8827 8768
/ "2024-03-02" "2024-01-03"
```

## HTTP

For performing requests:

- `http.get s` perform a get request agains url `s`

For making a http server:

- `http.register[s;f]` register a function `f` that serves requests to path `s`
- `http.serve s` Begin the web server binding to address `s`
