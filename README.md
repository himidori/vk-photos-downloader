# tool to download photos from VK dialogs

# usage, CLI arguments
```
-u : your VK login (username or phone number. used if no token was provided)
-p : your VK password (used if no token was provided)
-t : VK access token (to be used instead of user/pass)
-uid : user id (dialog) to download photos from (omit to download photos from every available dialog)
-r : amount of goroutines to use for concurrent photo download
-d : device to use for authorization (0 - iPhone, 1 - android, 2 - WPhone)
-h : print help
```

# example of usage
downloading photos from one dialog using 20 goroutines, authorizing by login and password
```
./vkphotos -u xxx -p xxx -uid xxx -r 20
```
downloading photos from every dialog using 20 goroutines, authorizing by token
```
./vkphotos -t xxx -r 20
```
