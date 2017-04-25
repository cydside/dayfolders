# dayfolders
It's a command line tool that creates daily folders for a selectable period of time to store files.

## Examples:

* `dayfolders -year=2017`  
It crates a default folders tree:  
```
.
├── 2017
│   ├── 01
...
│   ├── 02
│   │   ├── 01
│   │   ├── 02
...
│   │   └── 31
...
│   ├── 12
│   │   ├── 01
...
```

* `dayfolders -year=2017 -one -dow`  
It crates folders as follow:  
```
2017-01-01 (Sun)  
2017-01-02 (Mon)  
...  
2017-12-31 (Sun)  
```

* `dayfolders -from=2017-02-01 -to=2017-02-14 -one`  
It crates folders as follow:  
```
2017-02-01  
...  
2017-02-14  
```
you can run `dayfolders -from=2017-02-01 -days=14 -one` to get the same.  

* `dayfolders -to=2017-03-31 -days=10 -one -doy`   
It crates 10 folders adding day of the year as follow:  
```
2017-03-22 081
...
2017-03-31 090
```
