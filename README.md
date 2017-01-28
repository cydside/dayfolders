# dayfolders
The program creates daily folders to store files in a selectable period of time.

## Examples:

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
2017-01-01 (Sun)  
2017-01-02 (Mon)  
...  
2017-12-31 (Sun)  

* `dayfolders -from=2017-02-01 -to=2017-02-14 -one`  
It crates folders as follow:  
2017-02-01  
...  
2017-02-14  

you can run `dayfolders -from=2017-02-01 -days=14 -one` to get the same.  

* `dayfolders -to=2017-03-31 -days=10 -one`   
It crates 10 folders as follow:  
2017-03-22  
...  
2017-03-31
