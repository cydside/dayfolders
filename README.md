# dayfolders
A tool that creates daily folders for a selectable period of time.

You can add your language support, quick and easy, have a look at the source code. Enjoy!

## Usage examples:

* `dayfolders -f 2019`
It crates a default folders tree:
```
├── 2019
│   ├── 01
│   │   ├── 01
│   │   ├── ...
│   │   └── 31
│   ├── 02
│   ├── ...
│   └── 12
```

* `dayfolders -f 2019 -s -w 0`
It crates folders as follow:
```
├── 2019
│   ├── 2019-01-01_Tue
│   ├── 2019-01-02_Wed
│   ├── ...
│   └── 2019-12-31_Tue
```

* `dayfolders -f 2019-02 -s`
It crates folders as follow:
```
├── 2019
│   ├── 2019-02-01
│   ├── 2019-02-02
│   ├── ...
│   └── 2019-02-28
```

* `dayfolders -f 2019-11-01 -t 2019-11-08 -suffix="RX" -m 0`

The above two commands create a folder structure as follow:

```
├── 2019
│   ├── 11_Nov
│   │   ├── 01_RX
│   │   ├── ...
│   │   └── 08_RX
```

* `dayfolders -f 2019-11-01 -t 2019-11-08 -c "TODO,DONE" -m 0 -w 0`

The above two commands create a folder structure as follow:

```
├── 2019
│   ├── 11_Nov
│   │   ├── 01_Fri
|   |   |   ├──TODO
|   |   |   └──DONE
...
│   │   └── 08_Fri
|   |   |   ├──TODO
|   |   |   └──DONE
```

## TO DO:

* Import a json file for a new language.
