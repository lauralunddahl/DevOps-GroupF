## Week 1

This documents contains notes about the changes we did this week (following this week's tasks https://github.com/itu-devops/lecture_notes/blob/master/sessions/session_01/README_TASKS.md)

### Version control

After getting the files from the server we added version control. 
First we installed git on our Linux machines and followed the guide:
https://git-scm.com/book/en/v2/Git-Basics-Getting-a-Git-Repository (Initializing a Repository in an Existing Directory)

We have created a master branch which we plan on using as our development branch, and every week we will create a 
pull request with the work we have done and merge it into the main branch. 

### Migrating ITU-MiniTwit to run on a modern computer (running Linux)

We installed python3 (`apt search python3`) and pip3 for python (first we had to run `sudo apt update` and then `sudo apt install python3-pip`).
We also installed/ran the commands `sudo apt install libsqlite3-dev`, `apt search sqlite3` and lastly `sudo apt install sqlitebrowser` which we used to open/view the database.

We recompiled the flag_tool.c file using the following command `gcc flag_tool.c -l sqlite3 -o flag_tool`.

We then updated the python files from python2 to python3 using this command `2to3 -n -W --add-suffix=3 minitwit.py`
We saw that the following things was updated in the minitwit.py (by comparing it with the new .py3 file `diff minitwit.py minitwit.py3`):
```
11c11
< from __future__ import with_statement
---
> 
97c97
<     print "We got a visitor from: " + str(request.remote_addr)
---
>     print("We got a visitor from: " + str(request.remote_addr))
```

We deleted the .py3 file and ran the command `2to3 -w minitwit.py` to update the original file (not create a new output file).
In order to compile it, we had to install Flask `pip3 install Flask` and do some small changes in the code for it to run:

* Changed the import in the code from `werkzeug` to `werkzeug.security` as the `check_password_hash` and `generate_password_hash` packages are 
now located within .security.
* Removed the `tmp/` from the line: `DATABASE = '/tmp/minitwit.db'`
* Ended up changing it to the full path, as it couldn't locate it otherwise (not using ~ either): `DATABASE = '/home/parallels/Desktop/itu-minitwit/minitwit.db'`

### Run time dependencies

These were checked with the command `ldd ~/Desktop/itu-minitwit/flag_tool`. The output:
```
ldd ~/Desktop/itu-minitwit/flag_tool
	linux-vdso.so.1 (0x00007ffdc5d96000)
	libsqlite3.so.0 => /usr/lib/x86_64-linux-gnu/libsqlite3.so.0 (0x00007fd6af217000)
	libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007fd6af025000)
	libm.so.6 => /lib/x86_64-linux-gnu/libm.so.6 (0x00007fd6aeed6000)
	libpthread.so.0 => /lib/x86_64-linux-gnu/libpthread.so.0 (0x00007fd6aeeb3000)
	libdl.so.2 => /lib/x86_64-linux-gnu/libdl.so.2 (0x00007fd6aeead000)
	/lib64/ld-linux-x86-64.so.2 (0x00007fd6af35d000)
```

### Adapt shell script (`control.sh`)

We first installed shellcheck (`sudo apt install shellcheck`) and then ran the command `shellcheck control.sh`.
We then followed the recommendations and updated the `control.sh` file. The changes included:
* Adding `#!/bin/sh`
* Adding quotation marks around $1
