# ugrsical

Parse ZJU ugrs class timetable and generate iCalender file for you. If you are grs student, please use [this](https://github.com/determ1ne/grsical) repo.

 

### How to use 

 we provide two binary file, ugrsical and ugrsicalsrv

#### use ugrsical(locally)

- download from release page 
- don't need to edit configs
- use ` ./ugrsical -u [userId] -p [password]` to generate

#### use ugrsicalsrv(as server)

- download from release page

- edit configs/server.json

  here is the fields description

  ~~~js
  {
    "enckey": "D*F-JaNdRgUkXp2s5v8y/B?E(H+KbPeS",//aes key(must be 128bits or 192bits or 256bits)"
    "host": "localhost", //host path(localhost or your server domain)
    "port": 3000, //must be number
    "config": "configs/config.json", //path to config.json, can be empty(default)
    "ip_header": "", //if you use nginx as reverse proxy, you can set field ip-header to track real ip
    "redis_addr": "",//can be empty(default)
    "redis_pass": "" //can be empty(default)
  }
  
  ~~~

- use `./ugrsicalsrv` and open `host::port` to use!