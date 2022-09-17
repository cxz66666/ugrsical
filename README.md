# ugrsical

Parse ZJU ugrs class timetable and generate iCalender file for you. If you are grs student, please use [this](https://github.com/determ1ne/grsical) repo.

 

### How to use 

 we provide two binary file, ugrsical and ugrsicalsrv, both of them are same function, but one is running on your machine but the other is running on server and expose web page.

#### use ugrsical(local)

- download from release page 
- don't need to edit configs
- use ` ./ugrsical -u [userId] -p [password]` to generate
- also can use `./ugrsical --help` to find advanced usage

#### use ugrsicalsrv(server)

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
    "cache_ttl": 0  //can be empty(default 48h), unit is hour
  }
  
  ~~~

- use `./ugrsicalsrv` and open `host::port` to use!

### More configs
 
- we provide three config file under /configs: `upfile.json`, `config.json`, `server.json`
- `upfile.json` is used for ugrsical to get username and password from file
- `config.json` is used for ugrsical and ugrsicalsrv to get generated config, which is highly recommended not to modify, but if you are not satisfied with the generated config, you can modify it with follow rules:
  ~~~js
  {
  "lastUpdated": 20220913, //just a flag
  "tweaks": [   //some holidays for zju course schedule changed
    {
      "TweakType": 2, //0:clear,1:copy,2:exchange
      "Description": "[国庆节] 放假调休，课程对调", //description
      "From": 20221006, //format YYYYMMDD
      "To": 20221008 
    },
  ],
  "termConfigs": [ // detailed term informations
    {
      "Year": "2022-2023",
      "Term": 0, // see pkg/zjuservice/class.go for details
      "Begin": 20220912,
      "End": 20221106,
      "FirstWeekNo": 1
    },
    {
      "Year": "2022-2023",
      "Term": 1,
      "Begin": 20221107,
      "End": 20230111,
      "FirstWeekNo": 1
    }
  ],
  "classTerms": [ //want query class terms, must be in termConfigs
    "2022-2023:0",//see pkg/zjuservice/class.go for details
    "2022-2023:1"
  ],
  "examTerms": [ // want query exam terms, must be in termConfigs
    "2022-2023:0"  //0 is fall and winter term, 1 is spring and summer term
  ]
  }
  ~~~
  
- server.json is used for ugrsicalsrv to get server config, which is must modified for yourown server
- anyway, please don't move config file to another path(which is hard code!)

### Q&A

- **Q: I generate a ics file by ugrsical, but can't import it on my iphone?**
- A: ical file can be directly import on Mac, Windows, Linux and Android, but can't be opened on iphone directly, you need send it to your email(or others email) and open it on `mail` APP`, or you can use our online playground!
- **Q: I can't generate ics file by ugrsical, even I can't run it!**
- A: please check your download file from release page, you need download the binary file which is suitable for your machine, if you are fresh man and not sure, please download the `X86_64` version, and if you are using `arm`, please download the `Arm64` version, if you are using M1/M2 mac, please download the `Darwin_arm64` version.
- **Q: I am worried about the security of online services, what should I do if my account number and password are stolen**
- **A: The author pledges his reputation not to conduct private collections, you also can build your own server by manual**

### Disclaimer
该项目仅供学习交流使用，作者不对产生结果正确性与时效性做实时保证，使用者需自行承担因程序逻辑错误或课程时间变动导致的后果。
