# dayfolders
It's a command line tool that creates daily folders to store files in a selectable period of time.

You can add your language support, quick and easy, have a look at the source code. Enjoy!

## Usage examples:

* `dayfolders -year=2017`  
It crates a default folders tree:  
```
2017  
  -> 01  
     -> 01  
     -> 02  
     ...  
     -> 31  
  ...  
  -> 12  
     ...  
     -> 31  
```

* `dayfolders -year=2017 -one -dow`  
It crates folders as follow:  
```
2017-01-01_Sun  
2017-01-02_Mon  
...  
2017-12-31_Sun  
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
2017-03-22_081
...
2017-03-31_090
```

* `dayfolders -from=2017-02-01 -to=2017-02-14 -suffix="TO" -dom=1`
* `dayfolders -from=2017-02-01 -to=2017-02-14 -suffix="FROM" -dom=1`

The above two commands create a folder structure as follow:

```
2017  
  -> 02_Feb  
     -> 01_FROM  
     -> 01_TO  
     -> 02_FROM  
     -> 02_TO  
     ...  
     -> 14_FROM  
     -> 14_TO  
```
